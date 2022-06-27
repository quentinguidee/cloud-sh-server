package user

import (
	"database/sql"
	"errors"
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
		err = errors.New(fmt.Sprintf("the user '%s' doesn't exists", username))
		context.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		err = errors.New(fmt.Sprintf("Couldn't retrieve the user '%s'.", username))
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	context.JSON(http.StatusOK, user)
}
