package storage

import (
	"errors"
	"net/http"
	. "self-hosted-cloud/server/commands"
	commands "self-hosted-cloud/server/commands/storage"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/models/storage"
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

	var (
		bucket Bucket
		nodes  []Node
	)

	transactionError := NewTransaction([]Command{
		commands.GetUserBucketCommand{
			Database:       db,
			User:           &user,
			ReturnedBucket: &bucket,
		},
		commands.GetNodesCommand{
			Database:      db,
			Path:          path,
			Bucket:        &bucket,
			ReturnedNodes: &nodes,
		},
	}).Try()

	if transactionError != nil {
		c.AbortWithError(transactionError.Code(), transactionError.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": nodes,
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

	var (
		bucket    Bucket
		directory Node
	)

	node := Node{
		Filename: params.Name,
		Filetype: params.Type,
	}

	transactionError := NewTransaction([]Command{
		commands.GetUserBucketCommand{
			Database:       db,
			User:           &user,
			ReturnedBucket: &bucket,
		},
		commands.GetBucketNodeCommand{
			Database:     db,
			Path:         path,
			Bucket:       &bucket,
			ReturnedNode: &directory,
		},
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
			Node: &node,
			Path: path,
		},
	}).Try()

	if transactionError != nil {
		c.AbortWithError(transactionError.Code(), transactionError.Error())
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

	var (
		node   Node
		bucket Bucket
	)

	transactionError := NewTransaction([]Command{
		commands.GetUserBucketCommand{
			Database:       db,
			User:           &user,
			ReturnedBucket: &bucket,
		},
		commands.GetBucketNodeCommand{
			Database:     db,
			Path:         path,
			Bucket:       &bucket,
			ReturnedNode: &node,
		},
		commands.DeleteBucketNodeRecursivelyCommand{
			Node:     &node,
			Path:     path,
			Database: db,
		},
	}).Try()

	if transactionError != nil {
		c.AbortWithError(transactionError.Code(), transactionError.Error())
		return
	}
}
