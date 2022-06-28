package auth

import (
	"errors"
	"net/http"
	"self-hosted-cloud/server/commands/auth"
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
	"self-hosted-cloud/server/routes/auth/github"

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
	var session Session

	err := c.BindJSON(&session)
	if err != nil {
		err = errors.New("body can't be decoded into a Session object")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	db := database.GetDatabaseFromContext(c)
	commandError := auth.DeleteSessionCommand{
		Database: db,
		Session:  session,
	}.Run()

	if err != nil {
		c.AbortWithError(commandError.Code(), commandError.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Disconnected successfully.",
	})
}
