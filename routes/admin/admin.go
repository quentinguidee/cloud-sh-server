package admin

import (
	"net/http"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/middlewares"
	. "self-hosted-cloud/server/models"
	"self-hosted-cloud/server/services/admin"
	"self-hosted-cloud/server/utils"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine) {
	group := router.Group("/admin")
	{
		group.Use(middlewares.AdminMiddleware())
		group.GET("/demo", getDemoMode)
		group.POST("/demo", enableDemoMode)
		group.POST("/reset", hardReset)
	}
}

func getDemoMode(c *gin.Context) {
	appIsInDemoMode, err := admin.AppIsInDemoMode()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !appIsInDemoMode {
		c.JSON(http.StatusOK, gin.H{
			"demo_mode": DemoMode{Enabled: false},
		})
		return
	}

	demoMode, err := admin.GetDemoMode()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"demo_mode": demoMode,
	})
}

func enableDemoMode(c *gin.Context) {
	tx := database.NewTX(c)

	err := admin.ResetServer(tx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = admin.SetupDemoMode(DemoMode{
		Enabled:       true,
		ResetInterval: "0 0 0 * * *",
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = admin.StartDemoMode(tx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx.Commit()
}

func hardReset(c *gin.Context) {
	_, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx := database.NewTX(c)

	err = admin.ResetServer(tx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx.Commit()
}
