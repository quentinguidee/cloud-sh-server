package auth

import (
	"net/http"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
)

type CreateGithubUserCommand struct {
	User           *User
	Database       Database
	GithubUsername string
}

func (c CreateGithubUserCommand) Run() ICommandError {
	request := "INSERT INTO auth_github(username, user_id) VALUES (?, ?)"

	_, err := c.Database.Instance.Exec(request, c.GithubUsername, c.User.Id)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}

	return nil
}

func (c CreateGithubUserCommand) Revert() ICommandError {
	request := "DELETE FROM auth_github WHERE username = ? AND user_id = ?"

	_, err := c.Database.Instance.Exec(request, c.GithubUsername, c.User.Id)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}
