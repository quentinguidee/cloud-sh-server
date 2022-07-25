package api

import (
	"errors"
	"net/http"
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
	services "self-hosted-cloud/server/services/auth"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) (int, error) {
	// Get session from body
	var session Session
	if err := c.BindJSON(&session); err != nil {
		err = errors.New("body can't be decoded into a Session object")
		return http.StatusBadRequest, err
	}

	tx := database.NewTX(c)

	if err := services.DeleteSession(tx, &session); err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	return http.StatusOK, nil
}
