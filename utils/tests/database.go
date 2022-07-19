package tests

import (
	"context"
	"self-hosted-cloud/server/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func NewDB() (database.Database, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	dbx := sqlx.NewDb(db, "sqlmock")
	return database.New(dbx), mock
}

func NewTX(db *database.Database) *sqlx.Tx {
	tx, err := db.Instance.BeginTxx(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	return tx
}
