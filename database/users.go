package database

import (
	"database/sql"
	. "self-hosted-cloud/server/models"
)

func (db *Database) CreateUsersTable() (sql.Result, error) {
	return db.instance.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id       INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(255) UNIQUE,
			name     VARCHAR(255)
		)
	`)
}

func (db *Database) GetUser(username string) (User, error) {
	request := "SELECT id, username, name FROM users WHERE username = ?"

	var user User
	err := db.instance.QueryRow(request, username).Scan(&user.Id, &user.Username, &user.Name)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *Database) CreateUser(user User) (int, error) {
	request := "INSERT INTO users(username, name) VALUES (?, ?) RETURNING id"

	var id int
	err := db.instance.QueryRow(request, user.Username, user.Name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
