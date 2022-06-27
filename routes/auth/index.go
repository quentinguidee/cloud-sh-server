package auth

import (
	"errors"
	"net/http"
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

func logout(context *gin.Context) {
	var session Session

	err := context.BindJSON(&session)
	if err != nil {
		err = errors.New("body can't be decoded into a Session object")
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	db := context.MustGet(database.KeyDatabase).(database.Database)
	err = db.CloseSession(session)
	if err != nil {
		err = errors.New("this session doesn't exists")
		context.AbortWithError(http.StatusNotFound, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Disconnected successfully.",
	})
}
