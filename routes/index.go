package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAbout(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "Server is Running.",
	})
}
