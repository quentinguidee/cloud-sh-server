package database

func (db *Database) CreateSessionsTable() {
	_, _ = db.Instance.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id      INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			token   VARCHAR(255) UNIQUE
		)
	`)
}
