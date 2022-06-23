package user

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"self-hosted-cloud/server/database"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(testing *testing.T) {
	db, mock, _ := sqlmock.New()

	rows := sqlmock.NewRows([]string{"id", "username", "name"}).
		AddRow(2, "username", "Name")

	mock.ExpectPrepare("^SELECT (.+) FROM users WHERE username = \\?$").
		ExpectQuery().
		WithArgs("username").
		WillReturnRows(rows)

	router := gin.New()
	router.Use(database.Middleware(database.New(db)))

	LoadRoutes(router.Group("/auth"))

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/user/username", nil)

	router.ServeHTTP(recorder, req)

	assert.Equal(testing, http.StatusOK, recorder.Code)
	assert.Equal(testing, `{"id":2,"name":"Name","username":"username"}`, recorder.Body.String())
}

func TestGetNonExistingUser(testing *testing.T) {
	db, mock, _ := sqlmock.New()

	mock.ExpectPrepare("^SELECT (.+) FROM users WHERE username = \\?$").
		ExpectQuery().
		WithArgs("username").
		WillReturnError(sql.ErrNoRows)

	router := gin.New()
	router.Use(database.Middleware(database.New(db)))

	LoadRoutes(router.Group("/auth"))

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/user/username", nil)

	router.ServeHTTP(recorder, req)

	assert.Equal(testing, http.StatusNotFound, recorder.Code)
}
