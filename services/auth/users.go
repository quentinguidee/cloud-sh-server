package auth

import (
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
	"self-hosted-cloud/server/models/types"
	. "self-hosted-cloud/server/services"
	"time"

	"github.com/jmoiron/sqlx"
)

func CreateUser(tx *sqlx.Tx, username string, name string, profilePicture string, role string) (User, IServiceError) {
	query := `
		INSERT INTO users(username, name, profile_picture, role, creation_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	user := User{
		Username:       username,
		Name:           name,
		ProfilePicture: types.NewNullableString(profilePicture),
		Role:           types.NewNullableString(role),
		CreationDate:   types.NewNullableTime(time.Now()),
	}

	err := database.
		NewRequest(tx, query).
		QueryRow(user.Username, user.Name, user.ProfilePicture, user.Role, user.CreationDate).
		Scan(&user.Id).
		OnError("failed to create the user")

	return user, err
}

func GetUser(tx *sqlx.Tx, username string) (User, IServiceError) {
	query := "SELECT * FROM users WHERE username = $1"

	var user User

	err := database.
		NewRequest(tx, query).
		Get(&user, username).
		OnError("failed to retrieve the user")

	return user, err
}

func GetUserFromToken(tx *sqlx.Tx, token string) (User, IServiceError) {
	query := `
		SELECT users.*
		FROM users, sessions
		WHERE sessions.user_id = users.id
		  AND sessions.token = $1
	`

	var user User

	err := database.
		NewRequest(tx, query).
		Get(&user, token).
		OnError("the user is not connected")

	return user, err
}

func GetUsersByRole(tx *sqlx.Tx, role string) ([]User, IServiceError) {
	query := "SELECT * FROM users WHERE role = $1"

	var users []User

	err := database.
		NewRequest(tx, query).
		Select(&users, role).
		OnError("failed to retrieve the user by its role")

	return users, err
}
