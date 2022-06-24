package database

import (
	"database/sql"
	. "self-hosted-cloud/server/models"
)

func (db *Database) CreateGithubAuthTable() (sql.Result, error) {
	return db.instance.Exec(`
		CREATE TABLE IF NOT EXISTS auth_github (
			username VARCHAR(255) UNIQUE PRIMARY KEY,
			user_id  INTEGER,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)
	`)
}

func (db *Database) GetUserFromGithub(username string) (User, error) {
	request := `
		SELECT users.id, users.username, users.name
		FROM users, auth_github
		WHERE users.id = auth_github.user_id
		  AND auth_github.username = ?;
	`

	var user User
	err := db.instance.QueryRow(request, username).Scan(&user.Id, &user.Username, &user.Name)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *Database) CreateUserFromGithub(githubUser GithubUser) (User, error) {
	user := User{
		Username: githubUser.Login,
		Name:     githubUser.Name,
	}
	userId, err := db.CreateUser(user)
	if err != nil {
		return User{}, err
	}

	_, err = db.instance.Exec(`INSERT INTO auth_github(username, user_id) VALUES (?, ?)`, githubUser.Login, userId)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
