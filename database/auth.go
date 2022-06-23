package database

import "database/sql"

func (db *Database) CreateGitHubAuthTable() (sql.Result, error) {
	return db.instance.Exec(`
		CREATE TABLE IF NOT EXISTS auth_github (
			username VARCHAR(255),
			user_id  INTEGER
		)
	`)
}
