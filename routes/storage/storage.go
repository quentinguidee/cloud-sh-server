package storage

import (
	"errors"
	"net/http"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/models"
	"self-hosted-cloud/server/models/types"
	"self-hosted-cloud/server/services/storage"
	"self-hosted-cloud/server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine) {
	group := router.Group("/storage")
	{
		group.GET("", getNodes)
		group.GET("/recent", getRecentFiles)
		group.PUT("", createNode)
		group.DELETE("", deleteNodes)
		group.PATCH("", renameNode)
		group.GET("/bucket", getBucket)
		group.GET("/download", downloadNodes)
		group.POST("/upload", uploadNode)
	}
}

func getNodes(c *gin.Context) {
	parentUuid := c.Query("parent_uuid")

	tx := database.NewTransaction(c)
	defer tx.Rollback()

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	directory, serviceError := storage.GetNode(tx, parentUuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	accessType, serviceError := storage.GetBucketUserAccessType(tx, directory.BucketId, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	if accessType < models.ReadOnly {
		err := errors.New("cannot access this bucket: insufficient permissions")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	nodes, serviceError := storage.GetNodes(tx, parentUuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	database.ExecTransaction(c, tx)

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
	})
}

func getRecentFiles(c *gin.Context) {
	tx := database.NewTransaction(c)
	defer tx.Rollback()

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	bucket, serviceError := storage.GetUserBucket(tx, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	accessType, serviceError := storage.GetBucketUserAccessType(tx, bucket.Id, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	if accessType < models.ReadOnly {
		err := errors.New("cannot access this bucket: insufficient permissions")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	nodes, serviceError := storage.GetRecentFiles(tx, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	database.ExecTransaction(c, tx)

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
	})
}

type CreateFilesParams struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

func createNode(c *gin.Context) {
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

	parentUuid := c.Query("parent_uuid")
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx := database.NewTransaction(c)
	defer tx.Rollback()

	bucket, serviceError := storage.GetUserBucket(tx, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	accessType, serviceError := storage.GetBucketUserAccessType(tx, bucket.Id, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	if accessType < models.Write {
		err := errors.New("cannot write in this bucket: insufficient permissions")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	nodeType := params.Type
	if nodeType != "directory" {
		nodeType = storage.DetectFileType(params.Name)
	}

	node, serviceError := storage.CreateNode(tx,
		user.Id,
		types.NewNullableString(parentUuid),
		bucket.Id,
		params.Name,
		nodeType,
		types.NewNullString(),
		types.NewNullableInt64(0),
	)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	path, serviceError := storage.GetNodePath(tx, node, bucket.Id, bucket.RootNodeUuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	serviceError = storage.CreateNodeInFileSystem(node.Type, path, "")
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	database.ExecTransaction(c, tx)
}

func deleteNodes(c *gin.Context) {
	uuid := c.Query("node_uuid")

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx := database.NewTransaction(c)
	defer tx.Rollback()

	bucket, serviceError := storage.GetUserBucket(tx, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	accessType, serviceError := storage.GetBucketUserAccessType(tx, bucket.Id, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	if accessType < models.Write {
		err := errors.New("cannot delete elements in this bucket: insufficient permissions")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	node, serviceError := storage.GetNode(tx, uuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	path, serviceError := storage.GetNodePath(tx, node, bucket.Id, bucket.RootNodeUuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	serviceError = storage.DeleteNodeRecursively(tx, &node)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	serviceError = storage.DeleteNodeInFileSystem(path)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	database.ExecTransaction(c, tx)
}

func renameNode(c *gin.Context) {
	uuid := c.Query("node_uuid")
	newName := c.Query("new_name")

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx := database.NewTransaction(c)
	defer tx.Rollback()

	bucket, serviceError := storage.GetUserBucket(tx, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	accessType, serviceError := storage.GetBucketUserAccessType(tx, bucket.Id, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	if accessType < models.Write {
		err := errors.New("cannot rename elements in this bucket: insufficient permissions")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	node, serviceError := storage.GetNode(tx, uuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	path, serviceError := storage.GetNodePath(tx, node, bucket.Id, bucket.RootNodeUuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	serviceError = storage.UpdateNode(tx, newName, node.Type, uuid, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	serviceError = storage.RenameNodeInFileSystem(path, newName)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	database.ExecTransaction(c, tx)
}

func getBucket(c *gin.Context) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx := database.NewTransaction(c)
	defer tx.Rollback()

	bucket, serviceError := storage.GetUserBucket(tx, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	database.ExecTransaction(c, tx)

	c.JSON(http.StatusOK, bucket)
}

func downloadNodes(c *gin.Context) {
	uuid := c.Query("node_uuid")

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx := database.NewTransaction(c)
	defer tx.Rollback()

	bucket, serviceError := storage.GetUserBucket(tx, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	accessType, serviceError := storage.GetBucketUserAccessType(tx, bucket.Id, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	if accessType < models.ReadOnly {
		err := errors.New("cannot download elements from this bucket: insufficient permissions")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	path, serviceError := storage.GetDownloadPath(tx, user.Id, uuid, bucket.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	database.ExecTransaction(c, tx)

	c.File(path)
}

func uploadNode(c *gin.Context) {
	parentUuid := c.Query("parent_uuid")
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx := database.NewTransaction(c)
	defer tx.Rollback()

	bucket, serviceError := storage.GetUserBucket(tx, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	accessType, serviceError := storage.GetBucketUserAccessType(tx, bucket.Id, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	if accessType < models.Write {
		err := errors.New("cannot upload elements in this bucket: insufficient permissions")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	nodeType := storage.DetectFileType(file.Filename)
	mime := storage.DetectFileMime(file)

	node, serviceError := storage.CreateNode(tx,
		user.Id,
		types.NewNullableString(parentUuid),
		bucket.Id,
		file.Filename,
		nodeType,
		types.NewNullableString(mime),
		types.NewNullableInt64(file.Size),
	)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	path, serviceError := storage.GetNodePath(tx, node, bucket.Id, bucket.RootNodeUuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	database.ExecTransaction(c, tx)
}
