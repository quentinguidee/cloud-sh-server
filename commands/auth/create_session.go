package auth

import (
	"crypto/rand"
	"fmt"
	"net/http"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
)

type CreateSessionCommand struct {
	Database        Database
	UserId          int
	ReturnedSession *Session
}

func (c CreateSessionCommand) Run() ICommandError {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}

	c.ReturnedSession.UserId = c.UserId
	c.ReturnedSession.Token = fmt.Sprintf("%X", token)

	request := "INSERT INTO sessions(user_id, token) VALUES (?, ?)"

	_, err = c.Database.Instance.Exec(request, c.ReturnedSession.UserId, c.ReturnedSession.Token)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}

	return nil
}

func (c CreateSessionCommand) Revert() ICommandError {
	request := "DELETE FROM sessions WHERE token = ?"

	_, err := c.Database.Instance.Exec(request, c.ReturnedSession.Token)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}
