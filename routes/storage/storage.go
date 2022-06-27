package storage

import (
	"errors"
	"net/http"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/models"
	"self-hosted-cloud/server/models/storage"
	"strings"

	"github.com/gin-gonic/gin"
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

	path := context.Query("path")
	token := context.GetHeader("Authorization")

	// TODO: db.VerifySession()

	user, err := db.GetUserFromSession(token)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	bucket, err := db.GetUserBucket(user.Id)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	files, err := db.GetFiles(bucket, path)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
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
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if strings.TrimSpace(params.Name) == "" {
		err = errors.New("filename cannot be empty")
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// TODO: db.VerifySession()

	path := context.Query("path")
	token := context.GetHeader("Authorization")

	user, err := db.GetUserFromSession(token)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	bucket, err := db.GetUserBucket(user.Id)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	directory, err := db.GetNode(bucket, path)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	node := storage.Node{
		Filename: params.Name,
		Filetype: params.Type,
		BucketId: bucket.Id,
	}

	err = db.CreateNode(directory.Id, node)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
