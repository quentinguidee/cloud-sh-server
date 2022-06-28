package database

func (db *Database) CreateUsersTable() {
	_, _ = db.Instance.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id              INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
			username        VARCHAR(255) UNIQUE,
			name            VARCHAR(255),
			profile_picture VARCHAR(255)
		)
	`)
}
