package utils

import (
	"self-hosted-cloud/server/commands/auth"
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"

	"github.com/gin-gonic/gin"
)

func GetTokenFromContext(c *gin.Context) string {
	return c.GetHeader("Authorization")
}

func GetUserFromContext(c *gin.Context) (User, error) {
	token := GetTokenFromContext(c)
	db := database.GetDatabaseFromContext(c)

	var user User
	err := auth.GetUserFromTokenCommand{
		Database:     db,
		Token:        token,
		ReturnedUser: &user,
	}.Run()

	if err != nil {
		return User{}, err.Error()
	}
	return user, nil
}
