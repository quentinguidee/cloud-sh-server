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

	directory, serviceError := storage.GetBucketNode(tx, parentUuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	bucket, serviceError := storage.GetBucketFromNode(tx, directory.Uuid)
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

	nodes, serviceError := storage.GetBucketNodes(tx, parentUuid)
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

	node, serviceError := storage.CreateBucketNode(tx, user.Id, params.Name, nodeType, "", 0)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	serviceError = storage.CreateBucketToNodeAssociation(tx, bucket.Id, node.Uuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	serviceError = storage.CreateBucketNodeAssociation(tx, parentUuid, node.Uuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	path, serviceError := storage.GetBucketNodePath(tx, node, bucket.Id, bucket.RootNode)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	serviceError = storage.CreateBucketNodeInFileSystem(node.Type, path, "")
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

	node, serviceError := storage.GetBucketNode(tx, uuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	path, serviceError := storage.GetBucketNodePath(tx, node, bucket.Id, bucket.RootNode)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	serviceError = storage.DeleteBucketNodeRecursively(tx, &node)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	serviceError = storage.DeleteBucketNodeInFileSystem(path)
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

	node, serviceError := storage.GetBucketNode(tx, uuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	path, serviceError := storage.GetBucketNodePath(tx, node, bucket.Id, bucket.RootNode)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	serviceError = storage.UpdateBucketNode(tx, newName, node.Type, uuid, user.Id)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	serviceError = storage.RenameBucketNodeInFileSystem(path, newName)
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

	path, serviceError := storage.GetDownloadPath(tx, user.Id, uuid, bucket.Id, bucket.RootNode)
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

	node, serviceError := storage.CreateBucketNode(tx, user.Id, file.Filename, nodeType, mime, file.Size)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	serviceError = storage.CreateBucketToNodeAssociation(tx, bucket.Id, node.Uuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	serviceError = storage.CreateBucketNodeAssociation(tx, parentUuid, node.Uuid)
	if serviceError != nil {
		serviceError.Throws(c)
		return
	}

	path, serviceError := storage.GetBucketNodePath(tx, node, bucket.Id, bucket.RootNode)
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
