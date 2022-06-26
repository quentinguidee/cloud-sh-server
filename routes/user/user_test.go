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

	rows := sqlmock.NewRows([]string{"id", "username", "name", "profile_picture"}).
		AddRow(2, "username", "Name", "https://google.com/")

	mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = \\?$").
		WithArgs("username").
		WillReturnRows(rows)

	router := gin.New()
	router.Use(database.Middleware(database.New(db)))

	LoadRoutes(router)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/username", nil)

	router.ServeHTTP(recorder, req)

	assert.Equal(testing, http.StatusOK, recorder.Code)
	assert.Equal(testing, `{"id":2,"username":"username","name":"Name","profile_picture":"https://google.com/"}`, recorder.Body.String())
}

func TestGetNonExistingUser(testing *testing.T) {
	db, mock, _ := sqlmock.New()

	mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = \\?$").
		WithArgs("username").
		WillReturnError(sql.ErrNoRows)

	router := gin.New()
	router.Use(database.Middleware(database.New(db)))

	LoadRoutes(router)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/username", nil)

	router.ServeHTTP(recorder, req)

	assert.Equal(testing, http.StatusNotFound, recorder.Code)
	assert.Equal(testing, `{"message":"User 'username' doesn't exists."}`, recorder.Body.String())
}
