package auth

import (
	"database/sql"
	"errors"
	"net/http"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
)

type GetUserCommand struct {
	Database     Database
	Username     string
	ReturnedUser *User
}

func (c GetUserCommand) Run() ICommandError {
	request := "SELECT id, username, name, profile_picture FROM users WHERE username = ?"

	err := c.Database.Instance.QueryRow(request, c.Username).Scan(
		&c.ReturnedUser.Id,
		&c.ReturnedUser.Username,
		&c.ReturnedUser.Name,
		&c.ReturnedUser.ProfilePicture)

	if err == sql.ErrNoRows {
		err = errors.New("the user 'username' doesn't exists")
		return NewError(http.StatusNotFound, err)
	}
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c GetUserCommand) Revert() ICommandError {
	return nil
}

type GetUserFromGithubCommand struct {
	Database       Database
	GithubUsername string
	ReturnedUser   *User
}

func (c GetUserFromGithubCommand) Run() ICommandError {
	request := `
		SELECT users.id, users.username, users.name, users.profile_picture
		FROM users, auth_github
		WHERE users.id = auth_github.user_id
		  AND auth_github.username = ?;
	`

	err := c.Database.Instance.QueryRow(request, c.GithubUsername).Scan(
		&c.ReturnedUser.Id,
		&c.ReturnedUser.Username,
		&c.ReturnedUser.Name,
		&c.ReturnedUser.ProfilePicture)

	if err == sql.ErrNoRows {
		return NewError(http.StatusNotFound, err)
	}
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c GetUserFromGithubCommand) Revert() ICommandError {
	return nil
}

type GetUserFromTokenCommand struct {
	Database     Database
	Token        string
	ReturnedUser *User
}

func (c GetUserFromTokenCommand) Run() ICommandError {
	request := `
		SELECT users.id, users.username, users.name, users.profile_picture
		FROM users, sessions
		WHERE sessions.user_id = users.id
		  AND sessions.token = ?
	`

	err := c.Database.Instance.QueryRow(request, c.Token).Scan(
		&c.ReturnedUser.Id,
		&c.ReturnedUser.Username,
		&c.ReturnedUser.Name,
		&c.ReturnedUser.ProfilePicture)

	if err == sql.ErrNoRows {
		err := errors.New("the user is not connected")
		return NewError(http.StatusNotFound, err)
	}
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c GetUserFromTokenCommand) Revert() ICommandError {
	return nil
}
