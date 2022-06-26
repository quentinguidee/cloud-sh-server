package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
	"self-hosted-cloud/server/routes/auth/github"
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
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Body can't be decoded into a Session object.",
		})
		return
	}

	db := context.MustGet(database.KeyDatabase).(database.Database)
	err = db.CloseSession(session)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "This session doesn't exists.",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Disconnected successfully.",
	})
}
