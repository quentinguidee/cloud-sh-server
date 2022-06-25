package storage

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"self-hosted-cloud/server/database"
)

func LoadRoutes(router *gin.Engine) {
	user := router.Group("/storage")
	{
		user.GET("/:bucket", get)
	}
}

type GetFilesParams struct {
	Path string
}

func get(context *gin.Context) {
	bucket := context.Param("bucket")
	db := context.MustGet(database.KeyDatabase).(database.Database)

	var params GetFilesParams
	files, err := db.GetFiles(bucket, params.Path)
	if err != nil {
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}
