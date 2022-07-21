package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"self-hosted-cloud/server/database"
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

	query := "INSERT INTO sessions(user_id, token) VALUES ($1, $2) RETURNING id"

	serviceError := database.
		NewRequest(tx, query).
		QueryRow(session.UserId, session.Token).
		Scan(&session.Id).
		OnError("failed to create the user session")

	return session, serviceError
}

func DeleteSession(tx *sqlx.Tx, session *Session) IServiceError {
	query := "DELETE FROM sessions WHERE token = $1 AND user_id = $2"

	res, serviceError := database.
		NewRequest(tx, query).
		Exec(session.Token, session.UserId).
		OnError("failed to delete the user session")

	if serviceError != nil {
		return serviceError
	}

	count, err := res.RowsAffected()
	if err != nil {
		return NewServiceError(http.StatusInternalServerError, err)
	}
	if count == 0 {
		err := errors.New("the session doesn't exists")
		return NewServiceError(http.StatusNotFound, err)
	}
	return nil
}
