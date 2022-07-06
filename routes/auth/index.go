package auth

import (
	"errors"
	"net/http"
	. "self-hosted-cloud/server/models"
	"self-hosted-cloud/server/routes/auth/github"
	services "self-hosted-cloud/server/services/auth"
	. "self-hosted-cloud/server/utils"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine) {
	group := router.Group("/auth")
	{
		github.LoadRoutes(group)
		group.POST("/logout", logout)
	}
}

func logout(c *gin.Context) {
	// Get session from body
	var session Session
	err := c.BindJSON(&session)
	if err != nil {
		err = errors.New("body can't be decoded into a Session object")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tx := NewTransaction(c)
	defer tx.Rollback()

	serviceError := services.DeleteSession(tx, &session)
	if err != nil {
		serviceError.Throws(c)
		return
	}

	ExecTransaction(c, tx)

	c.JSON(http.StatusOK, gin.H{
		"message": "Disconnected successfully.",
	})
}
