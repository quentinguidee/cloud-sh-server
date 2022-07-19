package auth

import (
	"self-hosted-cloud/server/models"
	"self-hosted-cloud/server/models/types"
	"self-hosted-cloud/server/utils/tests"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	db, mock := tests.NewDB()

	rows := sqlmock.NewRows([]string{"id", "username", "name", "profile_picture", "role"}).
		AddRow(2, "username", "Name", "https://google.com/", "user")

	mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE username = \\$1$").
		WithArgs("username").
		WillReturnRows(rows)
	mock.ExpectCommit()

	tx := tests.NewTX(&db)

	user, _ := GetUser(tx, "username")

	assert.Equal(t, models.User{
		Id:             2,
		Username:       "username",
		Name:           "Name",
		ProfilePicture: types.NewNullableString("https://google.com/"),
		Role:           types.NewNullableString("user"),
	}, user)
}
