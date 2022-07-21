package auth

import (
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/services"

	"github.com/jmoiron/sqlx"
)

func CreateGithubUser(tx *sqlx.Tx, userId int, username string) IServiceError {
	query := "INSERT INTO auth_github(username, user_id) VALUES ($1, $2)"

	_, err := database.
		NewRequest(tx, query).
		Exec(username, userId).
		OnError("failed to create the github user")

	return err
}

func GetGithubUser(tx *sqlx.Tx, username string) (User, IServiceError) {
	query := `
		SELECT users.*
		FROM users, auth_github
		WHERE users.id = auth_github.user_id
		  AND auth_github.username = $1
	`

	var user User

	err := database.
		NewRequest(tx, query).
		Get(&user, username).
		OnError("failed to get the github user")

	return user, err
}
