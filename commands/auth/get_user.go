package auth

import (
	"net/http"
	"self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
)

type GetUserFromGithubCommand struct {
	Database       Database
	GithubUsername string
	ReturnedUser   *User
}

func (c GetUserFromGithubCommand) Run() commands.ICommandError {
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

	if err != nil {
		return commands.NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c GetUserFromGithubCommand) Revert() commands.ICommandError {
	return nil
}
