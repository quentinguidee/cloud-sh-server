package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	. "self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/services"

	"github.com/jmoiron/sqlx"
)

func CreateSession(tx *sqlx.Tx, userId int) (Session, IServiceError) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return Session{}, NewServiceError(http.StatusInternalServerError, err)
	}

	session := Session{
		UserId: userId,
		Token:  fmt.Sprintf("%X", token),
	}

	request := "INSERT INTO sessions(user_id, token) VALUES (?, ?) RETURNING id"

	err = tx.QueryRow(request, session.UserId, session.Token).Scan(&session.Id)
	if err != nil {
		return Session{}, NewServiceError(http.StatusInternalServerError, err)
	}

	return session, nil
}

func DeleteSession(tx *sqlx.Tx, session *Session) IServiceError {
	request := "DELETE FROM sessions WHERE token = ? AND user_id = ?"

	res, err := tx.Exec(request, session.Token, session.UserId)
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, err)
	}

	count, err := res.RowsAffected()
	if count == 0 && err != nil {
		err := errors.New("the session doesn't exists")
		return NewServiceError(http.StatusNotFound, err)
	}

	return nil
}
