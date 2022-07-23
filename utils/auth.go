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
	db := database.GetDatabaseFromContext(c)
	return auth.GetUserFromToken(db, token)
}
