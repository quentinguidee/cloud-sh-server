package utils

import (
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
	"self-hosted-cloud/server/services/auth"

	"github.com/gin-gonic/gin"
)

func GetTokenFromContext(c *gin.Context) string {
	return c.GetHeader("Authorization")
}

func GetUserFromContext(c *gin.Context) (User, error) {
	token := GetTokenFromContext(c)

	tx := database.NewTransaction(c)
	defer tx.Rollback()

	user, err := auth.GetUserFromToken(tx, token)
	if err != nil {
		return User{}, err.Error()
	}

	database.ExecTransaction(c, tx)

	return user, nil
}
