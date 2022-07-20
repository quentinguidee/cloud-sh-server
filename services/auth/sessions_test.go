package auth

import (
	"net/http"
	"self-hosted-cloud/server/models"
	"self-hosted-cloud/server/utils/tests"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateSession(t *testing.T) {
	db, mock := tests.NewDB()

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO sessions(.*)$").
		WithArgs(1, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	tx := tests.NewTX(&db)
	session, err := CreateSession(tx, 1)
	if err != nil {
		assert.NoError(t, err.Error())
	}

	assert.Equal(t, 1, session.Id)
	assert.Equal(t, 64, len(session.Token))
	assert.Equal(t, 1, session.UserId)
}

func TestDeleteSession(t *testing.T) {
	session := models.Session{
		UserId: 1,
		Token:  "ABCDEF",
	}

	t.Run("delete session", func(t *testing.T) {
		db, mock := tests.NewDB()
		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM sessions WHERE token = \\$1 AND user_id = \\$2$").
			WithArgs("ABCDEF", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := tests.NewTX(&db)
		err := DeleteSession(tx, &session)
		if err != nil {
			assert.NoError(t, err.Error())
		}
	})

	t.Run("delete session not found", func(t *testing.T) {
		db, mock := tests.NewDB()

		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM sessions WHERE token = \\$1 AND user_id = \\$2$").
			WithArgs("ABCDEF", 1).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		tx := tests.NewTX(&db)
		err := DeleteSession(tx, &session)
		assert.Equal(t, err.Error().Error(), "the session doesn't exists")
		assert.Equal(t, err.Code(), http.StatusNotFound)
	})
}
