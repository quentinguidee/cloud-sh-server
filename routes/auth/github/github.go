package github

import (
	authContext "context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/models"
	"self-hosted-cloud/server/services/auth"
	"self-hosted-cloud/server/services/storage"
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

type GithubUser struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
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
	cdb := authContext.Background()
	token, err := config.Exchange(cdb, params.Code)
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
	client := config.Client(cdb, token)
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

	tx := database.NewTX(c)

	// Create account if it doesn't exist
	user, err := auth.GetGithubUser(tx, githubUser.Login)
	if err != nil {
		admins, err := auth.GetUsersByRole(tx, "admin")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var role string
		if len(admins) > 0 {
			role = "user"
		} else {
			role = "admin"
		}

		user = models.User{
			ID:             0,
			Username:       githubUser.Login,
			Name:           githubUser.Name,
			Email:          githubUser.Email,
			ProfilePicture: &githubUser.AvatarUrl,
			Role:           &role,
		}
		err = auth.CreateUser(tx, &user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err = auth.CreateGithubUser(tx, user.ID, user.Username)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err = storage.SetupDefaultBucket(tx, user.ID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	// Open session
	session, err := auth.CreateSession(tx, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx.Commit()

	// OK
	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"session": session,
	})
}
