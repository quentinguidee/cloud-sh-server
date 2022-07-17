package auth

import (
	"context"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newDB() (database.Database, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	dbx := sqlx.NewDb(db, "sqlmock")
	return database.New(dbx), mock
}

func newTX(db *database.Database) *sqlx.Tx {
	tx, err := db.Instance.BeginTxx(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	return tx
}

func TestGetUser(t *testing.T) {
	db, mock := newDB()

	rows := sqlmock.NewRows([]string{"id", "username", "name", "profile_picture"}).
		AddRow(2, "username", "Name", "https://google.com/")

	mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = \\$1$").
		WithArgs("username").
		WillReturnRows(rows)
	mock.ExpectCommit()

	tx := newTX(&db)

	user, _ := GetUser(tx, "username")

	assert.Equal(t, models.User{
		Id:             2,
		Username:       "username",
		Name:           "Name",
		ProfilePicture: "https://google.com/",
	}, user)
}
