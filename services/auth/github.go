package auth

import (
	"database/sql"
	"net/http"
	. "self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/services"
)

func CreateGithubUser(tx *sql.Tx, userId int, username string) IServiceError {
	request := "INSERT INTO auth_github(username, user_id) VALUES (?, ?)"

	_, err := tx.Exec(request, username, userId)
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, err)
	}
	return nil
}

func GetGithubUser(tx *sql.Tx, username string) (User, IServiceError) {
	request := `
		SELECT users.id, users.username, users.name, users.profile_picture
		FROM users, auth_github
		WHERE users.id = auth_github.user_id
		  AND auth_github.username = ?;
	`

	var user User
	err := tx.QueryRow(request, username).Scan(
		&user.Id,
		&user.Username,
		&user.Name,
		&user.ProfilePicture)

	if err == sql.ErrNoRows {
		return User{}, NewServiceError(http.StatusNotFound, err)
	}
	if err != nil {
		return User{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return user, nil
}
