package auth

import (
	"database/sql"
	"errors"
	"net/http"
	. "self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/services"

	"github.com/jmoiron/sqlx"
)

func CreateUser(tx *sqlx.Tx, username string, name string, profilePicture string) (User, IServiceError) {
	request := "INSERT INTO users(username, name, profile_picture) VALUES (?, ?, ?) RETURNING id"

	user := User{
		Username:       username,
		Name:           name,
		ProfilePicture: profilePicture,
	}

	err := tx.QueryRow(request, username, name, profilePicture).Scan(&user.Id)
	if err != nil {
		return User{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return user, nil
}

func GetUser(tx *sqlx.Tx, username string) (User, IServiceError) {
	request := "SELECT * FROM users WHERE username = ?"

	var user User
	err := tx.Get(&user, request, username)
	if err == sql.ErrNoRows {
		err = errors.New("the user 'username' doesn't exists")
		return User{}, NewServiceError(http.StatusNotFound, err)
	}
	if err != nil {
		return User{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return user, nil
}

func GetUserFromToken(tx *sqlx.Tx, token string) (User, IServiceError) {
	request := `
		SELECT users.*
		FROM users, sessions
		WHERE sessions.user_id = users.id
		  AND sessions.token = ?
	`

	var user User
	err := tx.Get(&user, request, token)
	if err == sql.ErrNoRows {
		err := errors.New("the user is not connected")
		return User{}, NewServiceError(http.StatusNotFound, err)
	}
	if err != nil {
		return User{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return user, nil
}
