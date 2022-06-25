package storage

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/models/storage"
	"strconv"
)

func LoadRoutes(router *gin.Engine) {
	user := router.Group("/storage")
	{
		user.GET("/:bucket", getFiles)
		user.PUT("/:bucket", createFile)
	}
}

type GetFilesParams struct {
	Path string
}

func getFiles(context *gin.Context) {
	bucket, err := strconv.Atoi(context.Param("bucket"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "The bucket ID must be a number.",
		})
		return
	}

	db := context.MustGet(database.KeyDatabase).(database.Database)

	var params GetFilesParams
	err = context.BindJSON(&params)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	files, err := db.GetFiles(bucket, params.Path)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

type CreateFilesParams struct {
	Path string
	Type string
	Name string
}

func createFile(context *gin.Context) {
	bucket, err := strconv.Atoi(context.Param("bucket"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "The bucket ID must be a number.",
		})
		return
	}

	db := context.MustGet(database.KeyDatabase).(database.Database)

	var params CreateFilesParams
	err = context.BindJSON(&params)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	directory, err := db.GetNode(bucket, params.Path)
	if err != nil {
		return
	}

	node := storage.Node{
		Filename:         params.Name,
		Filetype:         params.Type,
		InternalFilename: "UUID-TODO",
		BucketId:         bucket,
	}

	err = db.CreateNode(directory.Id, node)
	if err != nil {
		return
	}
}
