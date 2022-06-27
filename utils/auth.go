package utils

import (
	"database/sql"
	"errors"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/models"

	"github.com/gin-gonic/gin"
)

func GetTokenFromContext(c *gin.Context) string {
	return c.GetHeader("Authorization")
}

func GetUserFromContext(c *gin.Context) (models.User, error) {
	token := GetTokenFromContext(c)
	db := database.GetDatabaseFromContext(c)
	user, err := db.GetUserFromSession(token)
	if err == sql.ErrNoRows {
		return models.User{}, errors.New("user not connected")
	}
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
