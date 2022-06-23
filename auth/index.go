package auth

import (
	"github.com/gin-gonic/gin"
	"self-hosted-cloud/server/auth/github"
	"self-hosted-cloud/server/auth/user"
)

func LoadRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		github.LoadRoutes(auth)
		user.LoadRoutes(auth)
	}
}
