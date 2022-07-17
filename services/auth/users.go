package auth

import (
	"database/sql"
	"errors"
	"net/http"
	. "self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/services"

	"github.com/jmoiron/sqlx"
)

func CreateUser(tx *sqlx.Tx, username string, name string, profilePicture string, role string) (User, IServiceError) {
	request := "INSERT INTO users(username, name, profile_picture, role) VALUES ($1, $2, $3, $4) RETURNING id"

	user := User{
		Username:       username,
		Name:           name,
		ProfilePicture: profilePicture,
		Role:           role,
	}

	err := tx.QueryRow(request,
		username,
		name,
		profilePicture,
		role,
	).Scan(&user.Id)

	if err != nil {
		return User{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return user, nil
}

func GetUser(tx *sqlx.Tx, username string) (User, IServiceError) {
	request := "SELECT * FROM users WHERE username = $1"

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
		  AND sessions.token = $1
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

func GetUsersByRole(tx *sqlx.Tx, role string) ([]User, IServiceError) {
	request := "SELECT * FROM users WHERE role = $1"

	var users []User
	err := tx.Select(&users, request, role)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, err)
	}
	return users, nil
}
