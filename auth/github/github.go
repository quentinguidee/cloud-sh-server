package github

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func getConfig() oauth2.Config {
	return oauth2.Config{
		ClientID:     os.Getenv("AUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH_CLIENT_SECRET"),
		Scopes:       []string{"all"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		RedirectURL: "http://localhost:3000",
	}
}

func LoadRoutes(router *gin.RouterGroup) {
	github := router.Group("/github")
	{
		github.GET("/login", login)
		github.POST("/authorize", authorize)
	}
}

func login(context *gin.Context) {
	var config = getConfig()

	url := config.AuthCodeURL(os.Getenv("AUTH_STATE"), oauth2.AccessTypeOffline)

	context.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

func authorize(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "Server is Running.",
	})
}
