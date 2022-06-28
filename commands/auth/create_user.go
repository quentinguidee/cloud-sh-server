package auth

import (
	"net/http"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
)

type CreateUserCommand struct {
	User     *User
	Database Database
}

func (c CreateUserCommand) Run() ICommandError {
	request := "INSERT INTO users(username, name, profile_picture) VALUES (?, ?, ?) RETURNING id"

	err := c.Database.Instance.QueryRow(request,
		c.User.Username,
		c.User.Name,
		c.User.ProfilePicture,
	).Scan(&c.User.Id)

	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c CreateUserCommand) Revert() ICommandError {
	request := "DELETE FROM users WHERE id = ?"

	_, err := c.Database.Instance.Exec(request, c.User.Id)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}
