package database

import (
	. "self-hosted-cloud/server/models"
)

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

func (db *Database) GetUser(username string) (User, error) {
	request := "SELECT id, username, name, profile_picture FROM users WHERE username = ?"

	var user User
	err := db.Instance.QueryRow(request, username).Scan(&user.Id, &user.Username, &user.Name, &user.ProfilePicture)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *Database) CreateUser(user User) (int, error) {
	request := "INSERT INTO users(username, name, profile_picture) VALUES (?, ?, ?) RETURNING id"

	var id int
	err := db.Instance.QueryRow(request, user.Username, user.Name, user.ProfilePicture).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
