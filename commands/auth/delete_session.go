package auth

import (
	"errors"
	"net/http"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
)

type DeleteSessionCommand struct {
	Database Database
	Session  Session
}

func (c DeleteSessionCommand) Run() ICommandError {
	request := "DELETE FROM sessions WHERE token = ? AND user_id = ?"

	res, err := c.Database.Instance.Exec(request, c.Session.Token, c.Session.UserId)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		err := errors.New("the session doesn't exists")
		return NewError(http.StatusNotFound, err)
	}

	return nil
}

func (c DeleteSessionCommand) Revert() ICommandError {
	return CreateSessionCommand{
		Database:        c.Database,
		UserId:          c.Session.UserId,
		ReturnedSession: &Session{},
	}.Run()
}
