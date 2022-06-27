package github

import (
	authContext "context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"

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

func login(c *gin.Context) {
	var config = getConfig()

	url := config.AuthCodeURL(os.Getenv("AUTH_STATE"), oauth2.AccessTypeOffline)

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

type AuthorizeParams struct {
	Code  string
	State string
}

func callback(c *gin.Context) {
	var config = getConfig()
	var params AuthorizeParams

	// Decode JSON
	err := c.BindJSON(&params)
	if err != nil {
		err = errors.New("body can't be decoded into an AuthorizeParams object")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Ensure the AUTH_STATE is correct.
	if params.State != os.Getenv("AUTH_STATE") {
		err = errors.New("the states don't match")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Generate the Token
	ctx := authContext.Background()
	token, err := config.Exchange(ctx, params.Code)
	if err != nil {
		err = errors.New("the server failed to generate your token")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Get user info
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		err = errors.New("failed to fetch GitHub user account")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Get client info
	client := config.Client(ctx, token)
	res, err := client.Do(req)
	if err != nil {
		err = errors.New("failed to fetch GitHub user account")
		c.AbortWithError(res.StatusCode, err)
		return
	}

	// Decode client info
	var githubUser GithubUser
	err = json.NewDecoder(res.Body).Decode(&githubUser)
	if err != nil {
		err = errors.New("failed to decode GitHub user account")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Create account if it doesn't exist
	db := database.GetDatabaseFromContext(c)
	user, err := db.GetUserFromGithub(githubUser.Login)
	if err == sql.ErrNoRows {
		user, err = db.CreateUserFromGithub(githubUser)
		if err != nil {
			err = errors.New("failed to create the new user")
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		_, err := db.CreateBucket(user.Id)
		if err != nil {
			return
		}
	}
	if err != nil {
		err = errors.New("failed to retrieve the user")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Open session
	session, err := db.CreateSession(user.Id)
	if err != nil {
		err = errors.New("failed to create user session")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// OK
	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"session": session,
	})
}
