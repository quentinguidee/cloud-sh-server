package database

import (
	"database/sql"
	"errors"
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
	statement, err := db.instance.Prepare("SELECT id, username, name FROM users WHERE username = ?")
	if err != nil {
		return User{}, errors.New("failed to prepare statement")
	}

	var user User
	err = statement.QueryRow(username).Scan(&user.Id, &user.Username, &user.Name)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *Database) CreateUser(user User) error {
	statement, err := db.instance.Prepare("INSERT INTO users(username, name) VALUES (?, ?)")
	if err != nil {
		return errors.New("failed to prepare statement")
	}

	_, err = statement.Exec(user.Username, user.Name)
	if err != nil {
		return err
	}

	return nil
}
