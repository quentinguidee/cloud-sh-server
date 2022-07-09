package github

import (
	authContext "context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	. "self-hosted-cloud/server/models"
	"self-hosted-cloud/server/services/auth"
	"self-hosted-cloud/server/services/storage"
	. "self-hosted-cloud/server/utils"

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
		RedirectURL: os.Getenv("AUTH_REDIRECT_URI"),
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

type CallbackParams struct {
	Code  string
	State string
}

func callback(c *gin.Context) {
	var config = getConfig()
	var params CallbackParams

	// Decode JSON
	err := c.BindJSON(&params)
	if err != nil {
		err = errors.New("body can't be decoded into an CallbackParams object")
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

	tx := NewTransaction(c)
	defer tx.Rollback()

	// Create account if it doesn't exist
	user, serviceError := auth.GetGithubUser(tx, githubUser.Login)
	if serviceError != nil {
		if serviceError.Error() != sql.ErrNoRows {
			serviceError.Throws(c)
			return
		}

		user, serviceError = auth.CreateUser(tx, githubUser.Login, githubUser.Name, githubUser.AvatarUrl)
		if serviceError != nil {
			serviceError.Throws(c)
			return
		}

		serviceError = auth.CreateGithubUser(tx, user.Id, user.Username)
		if serviceError != nil {
			serviceError.Throws(c)
			return
		}

		serviceError = storage.SetupDefaultBucket(tx, user.Id)
		if serviceError != nil {
			serviceError.Throws(c)
			return
		}
	}

	// Open session
	session, serviceError := auth.CreateSession(tx, user.Id)
	if err != nil {
		serviceError.Throws(c)
		return
	}

	ExecTransaction(c, tx)

	// OK
	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"session": session,
	})
}
