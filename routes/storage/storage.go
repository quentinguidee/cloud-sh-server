package storage

import (
	"errors"
	"net/http"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/models"
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

	tx := database.NewTX(c)

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	directory, err := storage.GetNode(tx, parentUuid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	accessType, err := storage.GetBucketUserAccessType(tx, directory.BucketID, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if accessType < models.ReadOnly {
		err := errors.New("cannot access this bucket: insufficient permissions")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	nodes, err := storage.GetNodes(tx, parentUuid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
	})
}

func getRecentFiles(c *gin.Context) {
	tx := database.NewTX(c)

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	bucket, err := storage.GetUserBucket(tx, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	accessType, err := storage.GetBucketUserAccessType(tx, bucket.ID, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if accessType < models.ReadOnly {
		err := errors.New("cannot access this bucket: insufficient permissions")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	nodes, err := storage.GetRecentFiles(tx, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx.Commit()

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

	tx := database.NewTX(c)

	bucket, err := storage.GetUserBucket(tx, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	accessType, err := storage.GetBucketUserAccessType(tx, bucket.ID, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
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

	node, err := storage.CreateNode(tx, user.ID, models.Node{
		ParentUUID: parentUuid,
		BucketID:   bucket.ID,
		Name:       params.Name,
		Type:       nodeType,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	path, err := storage.GetNodePath(tx, node, bucket.ID, bucket.RootNode.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = storage.CreateNodeInFileSystem(node.Type, path, "")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx.Commit()
}

func deleteNodes(c *gin.Context) {
	uuid := c.Query("node_uuid")

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx := database.NewTX(c)

	bucket, err := storage.GetUserBucket(tx, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	accessType, err := storage.GetBucketUserAccessType(tx, bucket.ID, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if accessType < models.Write {
		err := errors.New("cannot delete elements in this bucket: insufficient permissions")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	node, err := storage.GetNode(tx, uuid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	path, err := storage.GetNodePath(tx, node, bucket.ID, bucket.RootNode.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = storage.DeleteNodeRecursively(tx, &node)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = storage.DeleteNodeInFileSystem(path)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx.Commit()
}

func renameNode(c *gin.Context) {
	uuid := c.Query("node_uuid")
	newName := c.Query("new_name")

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx := database.NewTX(c)

	bucket, err := storage.GetUserBucket(tx, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	accessType, err := storage.GetBucketUserAccessType(tx, bucket.ID, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if accessType < models.Write {
		err := errors.New("cannot rename elements in this bucket: insufficient permissions")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	node, err := storage.GetNode(tx, uuid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	path, err := storage.GetNodePath(tx, node, bucket.ID, bucket.RootNode.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	node.Name = newName

	err = storage.UpdateNode(tx, &node, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = storage.RenameNodeInFileSystem(path, newName)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx.Commit()
}

func getBucket(c *gin.Context) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx := database.NewTX(c)

	bucket, err := storage.GetUserBucket(tx, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, bucket)
}

func downloadNodes(c *gin.Context) {
	uuid := c.Query("node_uuid")

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx := database.NewTX(c)

	bucket, err := storage.GetUserBucket(tx, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	accessType, err := storage.GetBucketUserAccessType(tx, bucket.ID, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if accessType < models.ReadOnly {
		err := errors.New("cannot download elements from this bucket: insufficient permissions")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	path, err := storage.GetDownloadPath(tx, user.ID, uuid, bucket.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx.Commit()

	c.File(path)
}

func uploadNode(c *gin.Context) {
	parentUUID := c.Query("parent_uuid")
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

	tx := database.NewTX(c)

	bucket, err := storage.GetUserBucket(tx, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	accessType, err := storage.GetBucketUserAccessType(tx, bucket.ID, user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if accessType < models.Write {
		err := errors.New("cannot upload elements in this bucket: insufficient permissions")
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	nodeType := storage.DetectFileType(file.Filename)
	mime := storage.DetectFileMime(file)

	node, err := storage.CreateNode(tx, user.ID, models.Node{
		ParentUUID: parentUUID,
		BucketID:   bucket.ID,
		Name:       file.Filename,
		Type:       nodeType,
		Mime:       &mime,
		Size:       &file.Size,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	path, err := storage.GetNodePath(tx, node, bucket.ID, bucket.RootNode.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tx.Commit()
}
