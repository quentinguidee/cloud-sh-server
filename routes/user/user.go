package user

import (
	"net/http"
	"self-hosted-cloud/server/services/auth"
	"self-hosted-cloud/server/utils"

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

	tx := utils.NewTransaction(c)
	defer tx.Rollback()

	user, err := auth.GetUser(tx, username)
	if err != nil {
		err.Throws(c)
		return
	}

	utils.ExecTransaction(c, tx)

	c.JSON(http.StatusOK, user)
}
