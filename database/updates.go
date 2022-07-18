package database

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

type request struct {
	sql    string
	params []any
}

func (r *request) execute(tx *sqlx.Tx) error {
	_, err := tx.Exec(r.sql, r.params...)
	return err
}

type update struct {
	version  int
	requests []request
	db       *Database
}

func (db *Database) updateForVersion(version int) *update {
	return &update{
		version: version,
		db:      db,
	}
}

func (u *update) with(sql string, params ...any) *update {
	u.requests = append(u.requests, request{sql, params})
	return u
}

func (u *update) execute(currentVersion *int) {
	if u.version <= *currentVersion {
		return
	}
	log.Println("[DB UPDATE] Updating to version", u.version)

	u.requests = append(u.requests, request{
		sql:    "UPDATE servers SET database_version = $1 WHERE id = 1",
		params: []any{u.version},
	})

	tx, err := u.db.Instance.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Fatalln("Failed to initialize update transaction.")
	}
	defer tx.Rollback()

	for _, request := range u.requests {
		err := request.execute(tx)
		if err != nil {
			log.Fatalln("Failed to execute request:", request.sql)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalln("Failed to commit update transaction.")
	}

	*currentVersion = u.version
}

func (db *Database) Update() error {
	var currentVersion int
	err := db.Instance.QueryRowx("SELECT database_version FROM servers WHERE id = 1").Scan(&currentVersion)
	if err != nil {
		return err
	}

	db.updateForVersion(2).
		with("ALTER TABLE users ADD COLUMN creation_date TIMESTAMP").
		execute(&currentVersion)

	return nil
}
