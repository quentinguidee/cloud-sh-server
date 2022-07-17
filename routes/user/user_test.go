package user

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/middlewares"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(testing *testing.T) {
	db, mock, _ := sqlmock.New()
	dbx := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{"id", "username", "name", "profile_picture", "role"}).
		AddRow(2, "username", "Name", "https://google.com/", "user")

	mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = \\$1$").
		WithArgs("username").
		WillReturnRows(rows)
	mock.ExpectCommit()

	d := database.New(dbx)

	router := gin.New()
	router.Use(middlewares.DatabaseMiddleware(&d))

	LoadRoutes(router)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/username", nil)

	router.ServeHTTP(recorder, req)

	assert.Equal(testing, http.StatusOK, recorder.Code)
	assert.Equal(testing, `{"id":2,"username":"username","name":"Name","profile_picture":"https://google.com/","role":"user"}`, recorder.Body.String())
}

func TestGetNonExistingUser(testing *testing.T) {
	db, mock, _ := sqlmock.New()
	dbx := sqlx.NewDb(db, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = \\$1$").
		WithArgs("username").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	d := database.New(dbx)

	router := gin.New()
	router.Use(middlewares.DatabaseMiddleware(&d))
	router.Use(middlewares.ErrorMiddleware())

	LoadRoutes(router)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/username", nil)

	router.ServeHTTP(recorder, req)

	assert.Equal(testing, http.StatusNotFound, recorder.Code)
	assert.Equal(testing, `{"message":"the user 'username' doesn't exists"}`, recorder.Body.String())
}
