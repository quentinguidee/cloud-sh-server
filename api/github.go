package api

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

func Login(c *gin.Context) (int, error) {
	var config = getConfig()

	url := config.AuthCodeURL(os.Getenv("AUTH_STATE"), oauth2.AccessTypeOffline)

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})

	return http.StatusOK, nil
}

type CallbackParams struct {
	Code  string
	State string
}

type GitHubUser struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Login     string `json:"Login"`
	AvatarUrl string `json:"avatar_url"`
}

func Callback(c *gin.Context) (int, error) {
	var config = getConfig()
	var params CallbackParams

	// Decode JSON
	if err := c.BindJSON(&params); err != nil {
		err = errors.New("body can't be decoded into an CallbackParams object")
		return http.StatusBadRequest, err
	}

	// Ensure the AUTH_STATE is correct.
	if params.State != os.Getenv("AUTH_STATE") {
		err := errors.New("the states don't match")
		return http.StatusBadRequest, err
	}

	// Generate the Token
	cdb := authContext.Background()
	token, err := config.Exchange(cdb, params.Code)
	if err != nil {
		err = errors.New("the server failed to generate your token")
		return http.StatusInternalServerError, err
	}

	// Get user info
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		err = errors.New("failed to fetch GitHub user account")
		return http.StatusInternalServerError, err
	}

	// Get client info
	client := config.Client(cdb, token)
	res, err := client.Do(req)
	if err != nil {
		err = errors.New("failed to fetch GitHub user account")
		return res.StatusCode, err
	}

	// Decode client info
	var githubUser GitHubUser
	if err = json.NewDecoder(res.Body).Decode(&githubUser); err != nil {
		err = errors.New("failed to decode GitHub user account")
		return http.StatusInternalServerError, err
	}

	tx := database.NewTX(c)

	// Create account if it doesn't exist
	user, err := auth.GetGithubUser(tx, githubUser.Login)
	if err != nil {
		admins, err := auth.GetUsersByRole(tx, "admin")
		if err != nil {
			return http.StatusInternalServerError, err
		}

		var role string
		if len(admins) > 0 {
			role = "user"
		} else {
			role = "admin"
		}

		user = models.User{
			Username:       githubUser.Login,
			Name:           githubUser.Name,
			Email:          githubUser.Email,
			ProfilePicture: &githubUser.AvatarUrl,
			Role:           &role,
		}

		if err = auth.CreateUser(tx, &user); err != nil {
			return http.StatusInternalServerError, err
		}

		if err = auth.CreateGithubUser(tx, user.ID, user.Username); err != nil {
			return http.StatusInternalServerError, err
		}

		if err := storage.SetupDefaultBucket(tx, user.ID); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	// Open session
	session := models.Session{UserID: user.ID}
	if err := auth.CreateSession(tx, &session); err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	// OK
	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"session": session,
	})

	return http.StatusOK, nil
}
