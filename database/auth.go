package database

func (db *Database) CreateGithubAuthTable() {
	_, _ = db.Instance.Exec(`
		CREATE TABLE IF NOT EXISTS auth_github (
			username VARCHAR(255) UNIQUE PRIMARY KEY,
			user_id  INTEGER,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)
	`)
}
