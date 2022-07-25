package api

import (
	"net/http"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/services/auth"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) (int, error) {
	username := c.Param("username")

	tx := database.NewTX(c)

	user, err := auth.GetUser(tx, username)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	c.JSON(http.StatusOK, user)

	return http.StatusOK, nil
}
