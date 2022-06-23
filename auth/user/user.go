package user

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"self-hosted-cloud/server/database"
)

func LoadRoutes(router *gin.RouterGroup) {
	user := router.Group("/")
	{
		user.GET("/user/:username", getUser)
	}
}

func getUser(context *gin.Context) {
	username := context.Param("username")
	db := context.MustGet(database.KeyDatabase).(database.Database)
	user, err := db.GetUser(username)
	if err == sql.ErrNoRows {
		context.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("User %s doesn't exists", username),
		})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Couldn't retrieve the user %s", username),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"id":       user.Id,
		"username": user.Username,
		"name":     user.Name,
	})
}
