package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.GET("/login", login)
	}
}

func login(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "Server is Running.",
	})
}
