package admin

import (
	"net/http"
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
	"self-hosted-cloud/server/services/admin"
	"self-hosted-cloud/server/utils"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine) {
	group := router.Group("/admin")
	{
		group.POST("/demo", enableDemoMode)
		group.POST("/reset", hardReset)
	}
}

func enableDemoMode(c *gin.Context) {
	db := database.GetDatabaseFromContext(c)

	admin.SetupDemoMode(DemoMode{
		Enabled:       true,
		ResetInterval: "0 0 0 * * *",
	})
	admin.ResetServer(db)
	admin.StartDemoMode(db)
}

func hardReset(c *gin.Context) {
	_, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	db := database.GetDatabaseFromContext(c)

	admin.ResetServer(db)
}
