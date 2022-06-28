package user

import (
	"net/http"
	"self-hosted-cloud/server/commands/auth"
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine) {
	group := router.Group("/user")
	{
		group.GET("/", getUser)
		group.GET("/:username", getUser)
	}
}

func getUser(c *gin.Context) {
	username := c.Param("username")
	db := database.GetDatabaseFromContext(c)

	var user User
	err := auth.GetUserCommand{
		Database:     db,
		Username:     username,
		ReturnedUser: &user,
	}.Run()

	if err != nil {
		c.AbortWithError(err.Code(), err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}
