package utils

import (
	"log"
	"net/http"
	"net/http/httptest"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/models"
	"self-hosted-cloud/server/utils/tests"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const token = "UNIQUE_TOKEN"

func TestGetTokenFromContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.Header.Set("Authorization", token)

	assert.Equal(t, token, GetTokenFromContext(c))
}

func TestGetUserFromContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.Header.Set("Authorization", token)

	t.Run("get user from context", func(t *testing.T) {
		user := models.User{
			Username: "jean.dupont",
			Name:     "Jean Dupont",
		}
		db, mock := tests.NewDB()
		c.Set(database.KeyDatabase, &db)

		usersRows := sqlmock.NewRows([]string{"username", "name"}).
			AddRow(user.Username, user.Name)

		mock.ExpectBegin()
		mock.ExpectQuery("^SELECT (.+) FROM users INNER JOIN sessions ON users.id = sessions.user_id WHERE sessions.token = \\$1$").
			WithArgs(token).
			WillReturnRows(usersRows)
		mock.ExpectCommit()

		userFromContext, err := GetUserFromContext(c)
		if err != nil {
			log.Fatalln(err.Error())
		}
		assert.Equal(t, user, userFromContext)
	})

	t.Run("get user from context not found", func(t *testing.T) {
		db, mock := tests.NewDB()
		c.Set(database.KeyDatabase, &db)

		usersRows := sqlmock.NewRows([]string{"username", "name"})

		mock.ExpectBegin()
		mock.ExpectQuery("^SELECT (.+) FROM users, sessions WHERE sessions.user_id = users.id AND sessions.token = \\$1$").
			WithArgs(token).
			WillReturnRows(usersRows)
		mock.ExpectCommit()

		_, err := GetUserFromContext(c)
		assert.Error(t, err)
	})
}
