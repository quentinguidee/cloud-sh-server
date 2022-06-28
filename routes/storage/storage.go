package storage

import (
	"errors"
	"net/http"
	. "self-hosted-cloud/server/commands"
	commands "self-hosted-cloud/server/commands/storage"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/models"
	"self-hosted-cloud/server/models/storage"
	"self-hosted-cloud/server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine) {
	group := router.Group("/storage")
	{
		group.GET("", getNodes)
		group.PUT("", createNode)
		group.DELETE("", deleteNodes)
	}
}

type GetFilesParams struct {
	Path    string         `json:"path,omitempty"`
	Session models.Session `json:"session"`
}

func getNodes(c *gin.Context) {
	db := database.GetDatabaseFromContext(c)

	path := c.Query("path")

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	bucket, err := db.GetUserBucket(user.Id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	files, err := db.GetFiles(bucket, path)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

type CreateFilesParams struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

func createNode(c *gin.Context) {
	db := database.GetDatabaseFromContext(c)

	var params CreateFilesParams
	err := c.BindJSON(&params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if strings.TrimSpace(params.Name) == "" {
		err = errors.New("filename cannot be empty")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	path := c.Query("path")
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	bucket, err := db.GetUserBucket(user.Id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	directory, err := db.GetNode(bucket, path)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	node := storage.Node{
		Filename: params.Name,
		Filetype: params.Type,
		BucketId: bucket.Id,
	}

	transactionError := NewTransaction([]Command{
		commands.CreateBucketNodeCommand{
			Node:     &node,
			Bucket:   &bucket,
			Database: db,
		},
		commands.CreateBucketNodeAssociationCommand{
			FromNode: &directory,
			ToNode:   &node,
			Database: db,
		},
		commands.CreateBucketNodeInFileSystemCommand{
			Node: node,
			Path: path,
		},
	}).Try()

	if transactionError != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func deleteNodes(c *gin.Context) {
	db := database.GetDatabaseFromContext(c)

	path := c.Query("path")

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	bucket, err := db.GetUserBucket(user.Id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	node, err := db.GetNode(bucket, path)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = db.DeleteRecursively(node, path)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
