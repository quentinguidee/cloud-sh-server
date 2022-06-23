package github

import (
	authContext "context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
	"os"
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
		RedirectURL: "http://localhost:3000/login",
	}
}

func LoadRoutes(router *gin.RouterGroup) {
	github := router.Group("/github")
	{
		github.GET("/login", login)
		github.POST("/callback", callback)
	}
}

func login(context *gin.Context) {
	var config = getConfig()

	url := config.AuthCodeURL(os.Getenv("AUTH_STATE"), oauth2.AccessTypeOffline)

	context.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

type AuthorizeParams struct {
	Code  string
	State string
}

type User struct {
	Email string
	Name  string
}

func callback(context *gin.Context) {
	var config = getConfig()
	var params AuthorizeParams

	// Decode JSON
	err := context.BindJSON(&params)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Body can't be decoded into an AuthorizeParams object.",
		})
		return
	}

	// Ensure the AUTH_STATE is correct.
	if params.State != os.Getenv("AUTH_STATE") {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "The states don't match.",
		})
		return
	}

	// Generate the Token
	ctx := authContext.Background()
	token, err := config.Exchange(ctx, params.Code)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "The server failed to generate your token.",
		})
		return
	}

	// Get user info
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch GitHub user account.",
		})
		return
	}

	// Get client info
	client := config.Client(ctx, token)
	res, err := client.Do(req)
	if err != nil {
		context.JSON(res.StatusCode, gin.H{
			"message": "Failed to fetch GitHub user account.",
		})
		return
	}

	// Decode client info
	var user User
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to decode GitHub user account.",
		})
		return
	}

	// OK
	context.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
