package api

import (
	"net/http"
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
	"self-hosted-cloud/server/services/admin"
	"self-hosted-cloud/server/utils"

	"github.com/gin-gonic/gin"
)

func GetDemoMode(c *gin.Context) (int, error) {
	appIsInDemoMode, err := admin.AppIsInDemoMode()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !appIsInDemoMode {
		c.JSON(http.StatusOK, gin.H{
			"demo_mode": DemoMode{Enabled: false},
		})
		return http.StatusOK, nil
	}

	demoMode, err := admin.GetDemoMode()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	c.JSON(http.StatusOK, gin.H{
		"demo_mode": demoMode,
	})

	return http.StatusOK, nil
}

func EnableDemoMode(c *gin.Context) (int, error) {
	tx := database.NewTX(c)

	if err := admin.ResetServer(tx); err != nil {
		return http.StatusInternalServerError, err
	}

	demoMode := DemoMode{
		Enabled:       true,
		ResetInterval: "0 0 0 * * *",
	}

	if err := admin.SetupDemoMode(demoMode); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := admin.StartDemoMode(tx); err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	return http.StatusOK, nil
}

func HardReset(c *gin.Context) (int, error) {
	if _, err := utils.GetUserFromContext(c); err != nil {
		return http.StatusInternalServerError, err
	}

	tx := database.NewTX(c)

	if err := admin.ResetServer(tx); err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	return http.StatusOK, nil
}
