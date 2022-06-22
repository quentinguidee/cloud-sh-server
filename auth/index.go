package auth

import (
	"self-hosted-cloud/server/auth/github"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		github.LoadRoutes(auth)
	}
}
