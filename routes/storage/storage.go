package storage

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/models"
	"self-hosted-cloud/server/models/storage"
)

func LoadRoutes(router *gin.Engine) {
	group := router.Group("/storage")
	{
		group.GET("", getFiles)
		group.PUT("", createFile)
	}
}

type GetFilesParams struct {
	Path    string         `json:"path,omitempty"`
	Session models.Session `json:"session"`
}

func getFiles(context *gin.Context) {
	db := context.MustGet(database.KeyDatabase).(database.Database)

	path := context.Param("path")
	token := context.GetHeader("Authorization")

	// TODO: db.VerifySession()

	user, err := db.GetUserFromSession(token)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	bucket, err := db.GetUserBucket(user.Id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	files, err := db.GetFiles(bucket, path)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

type CreateFilesParams struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

func createFile(context *gin.Context) {
	db := context.MustGet(database.KeyDatabase).(database.Database)

	var params CreateFilesParams
	err := context.BindJSON(&params)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// TODO: db.VerifySession()

	path := context.Param("path")
	token := context.GetHeader("Authorization")

	user, err := db.GetUserFromSession(token)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	bucket, err := db.GetUserBucket(user.Id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	directory, err := db.GetNode(bucket, path)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	node := storage.Node{
		Filename:         params.Name,
		Filetype:         params.Type,
		InternalFilename: "UUID-TODO",
		BucketId:         bucket.Id,
	}

	err = db.CreateNode(directory.Id, node)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
}
