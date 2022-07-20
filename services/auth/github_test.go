package auth

import (
	"database/sql"
	"net/http"
	"self-hosted-cloud/server/models"
	"self-hosted-cloud/server/utils/tests"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateGithubUser(t *testing.T) {
	db, mock := tests.NewDB()

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO auth_github(.*)$").
		WithArgs("jean.dupont", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx := tests.NewTX(&db)

	err := CreateGithubUser(tx, 1, "jean.dupont")
	if err != nil {
		assert.NoError(t, err.Error())
	}
}

func TestGetGithubUser(t *testing.T) {
	user := models.User{
		Username: "jean.dupont",
		Name:     "Jean Dupont",
	}

	t.Run("get github user", func(t *testing.T) {
		db, mock := tests.NewDB()

		rows := sqlmock.NewRows([]string{"username", "name"}).
			AddRow(user.Username, user.Name)

		mock.ExpectBegin()
		mock.ExpectQuery("^SELECT users(.*) FROM users, auth_github WHERE users.id = auth_github.user_id AND auth_github.username = \\$1$").
			WithArgs("jean.dupont").
			WillReturnRows(rows)
		mock.ExpectCommit()

		tx := tests.NewTX(&db)

		returnedUser, err := GetGithubUser(tx, "jean.dupont")
		if err != nil {
			assert.NoError(t, err.Error())
		}
		assert.Equal(t, user, returnedUser)
	})

	t.Run("get github user not found", func(t *testing.T) {
		db, mock := tests.NewDB()

		mock.ExpectBegin()
		mock.ExpectQuery("^SELECT users(.*) FROM users, auth_github WHERE users.id = auth_github.user_id AND auth_github.username = \\$1$").
			WithArgs("jean.dupont").
			WillReturnError(sql.ErrNoRows)
		mock.ExpectCommit()

		tx := tests.NewTX(&db)

		_, err := GetGithubUser(tx, "jean.dupont")
		assert.Error(t, err.Error())
		assert.Equal(t, err.Code(), http.StatusNotFound)
	})
}
