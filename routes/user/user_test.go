package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"self-hosted-cloud/server/middlewares"
	"self-hosted-cloud/server/utils/tests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	t.Run("get user", func(t *testing.T) {
		db, mock := tests.NewDB()

		creationDate := time.Now()
		creationDateString, _ := creationDate.MarshalJSON()

		rows := sqlmock.NewRows([]string{"id", "username", "name", "profile_picture", "role", "creation_date"}).
			AddRow(2, "username", "Name", "https://google.com/", "user", creationDate)

		mock.ExpectBegin()
		mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = \\$1$").
			WithArgs("username").
			WillReturnRows(rows)
		mock.ExpectCommit()

		router := gin.New()
		router.Use(middlewares.DatabaseMiddleware(&db))

		LoadRoutes(router)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/user/username", nil)

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, fmt.Sprintf(`{"id":2,"username":"username","name":"Name","profile_picture":"https://google.com/","role":"user","creation_date":%s}`, creationDateString), recorder.Body.String())
	})

	t.Run("get user not found", func(t *testing.T) {
		db, mock := tests.NewDB()

		mock.ExpectBegin()
		mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = \\$1$").
			WithArgs("username").
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		router := gin.New()
		router.Use(middlewares.DatabaseMiddleware(&db))
		router.Use(middlewares.ErrorMiddleware())

		LoadRoutes(router)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/user/username", nil)

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Equal(t, `{"message":"the user 'username' doesn't exists"}`, recorder.Body.String())
	})
}
