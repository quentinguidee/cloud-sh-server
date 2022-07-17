package auth

import (
	"database/sql"
	"net/http"
	. "self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/services"

	"github.com/jmoiron/sqlx"
)

func CreateGithubUser(tx *sqlx.Tx, userId int, username string) IServiceError {
	request := "INSERT INTO auth_github(username, user_id) VALUES ($1, $2)"

	_, err := tx.Exec(request, username, userId)
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, err)
	}
	return nil
}

func GetGithubUser(tx *sqlx.Tx, username string) (User, IServiceError) {
	request := `
		SELECT users.*
		FROM users, auth_github
		WHERE users.id = auth_github.user_id
		  AND auth_github.username = $1;
	`

	var user User
	err := tx.Get(&user, request, username)
	if err == sql.ErrNoRows {
		return User{}, NewServiceError(http.StatusNotFound, err)
	}
	if err != nil {
		return User{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return user, nil
}
