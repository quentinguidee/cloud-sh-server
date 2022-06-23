package auth

import (
	"github.com/gin-gonic/gin"
	"self-hosted-cloud/server/routes/auth/github"
)

func LoadRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		github.LoadRoutes(auth)
	}
}
