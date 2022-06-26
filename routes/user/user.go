package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"self-hosted-cloud/server/database"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine) {
	group := router.Group("/user")
	{
		group.GET("/", getUser)
		group.GET("/:username", getUser)
	}
}

func getUser(context *gin.Context) {
	username := context.Param("username")
	db := context.MustGet(database.KeyDatabase).(database.Database)
	user, err := db.GetUser(username)
	if err == sql.ErrNoRows {
		context.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("User '%s' doesn't exists.", username),
		})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Couldn't retrieve the user '%s'.", username),
		})
		return
	}

	context.JSON(http.StatusOK, user)
}
