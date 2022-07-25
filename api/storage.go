package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/models"
	"self-hosted-cloud/server/services/storage"
	"self-hosted-cloud/server/utils"
	"strings"
)

func getBucketUUID(c *gin.Context) (uuid.UUID, error) {
	bucketUUID, err := uuid.Parse(c.Param("bucket_uuid"))
	if err != nil {
		err := errors.New("bad request: failed to parse bucket_uuid")
		return bucketUUID, err
	}
	return bucketUUID, nil
}

func GetNodes(c *gin.Context) (int, error) {
	parentUUID := c.Query("parent_uuid")
	bucketUUID, err := getBucketUUID(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	tx := database.NewTX(c)

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := storage.AuthorizeAccess(tx, models.ReadOnly, bucketUUID, user.ID); err != nil {
		return http.StatusInternalServerError, err
	}

	nodes, err := storage.GetNodes(tx, parentUUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
	})

	return http.StatusOK, nil
}

func GetRecentFiles(c *gin.Context) (int, error) {
	tx := database.NewTX(c)

	bucketUUID, err := getBucketUUID(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = storage.AuthorizeAccess(tx, models.ReadOnly, bucketUUID, user.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	nodes, err := storage.GetRecentFiles(tx, user.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
	})

	return http.StatusOK, nil
}

func EmptyBin(c *gin.Context) (int, error) {
	tx := database.NewTX(c)

	bucketUUID, err := getBucketUUID(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err = storage.AuthorizeAccess(tx, models.ReadOnly, bucketUUID, user.ID); err != nil {
		return http.StatusInternalServerError, err
	}

	nodes, err := storage.GetDeletedNodes(tx, bucketUUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	bucket, err := storage.GetBucket(tx, bucketUUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	for _, node := range nodes {
		tx := database.NewTX(c)

		var path string

		if path, err = storage.GetNodePath(tx, node, bucketUUID, bucket.RootNode.UUID); err != nil {
			return http.StatusInternalServerError, err
		}

		if err = storage.DeleteNodeRecursively(tx, &node); err != nil {
			return http.StatusInternalServerError, err
		}

		if err = storage.DeleteNodeInFileSystem(path); err != nil {
			return http.StatusInternalServerError, err
		}

		tx.Commit()
	}

	return http.StatusOK, nil
}

func GetBin(c *gin.Context) (int, error) {
	tx := database.NewTX(c)

	bucketUUID, err := getBucketUUID(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err = storage.AuthorizeAccess(tx, models.ReadOnly, bucketUUID, user.ID); err != nil {
		return http.StatusInternalServerError, err
	}

	nodes, err := storage.GetDeletedNodes(tx, bucketUUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
	})

	return http.StatusOK, nil
}

type CreateFilesParams struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

func CreateNode(c *gin.Context) (int, error) {
	var params CreateFilesParams

	if err := c.BindJSON(&params); err != nil {
		return http.StatusBadRequest, err
	}

	if strings.TrimSpace(params.Name) == "" {
		err := errors.New("filename cannot be empty")
		return http.StatusBadRequest, err
	}

	parentUUID := c.Query("parent_uuid")

	bucketUUID, err := getBucketUUID(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx := database.NewTX(c)

	if err = storage.AuthorizeAccess(tx, models.Write, bucketUUID, user.ID); err != nil {
		return http.StatusInternalServerError, err
	}

	nodeType := params.Type
	if nodeType != "directory" {
		nodeType = storage.DetectFileType(params.Name)
	}

	bucket, err := storage.GetBucket(tx, bucketUUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	node := models.Node{
		ParentUUID: parentUUID,
		BucketUUID: bucket.UUID,
		Name:       params.Name,
		Type:       nodeType,
	}

	if err = storage.CreateNode(tx, user.ID, &node); err != nil {
		return http.StatusInternalServerError, err
	}

	path, err := storage.GetNodePath(tx, node, bucket.UUID, bucket.RootNode.UUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = storage.CreateNodeInFileSystem(node.Type, path, "")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	return http.StatusOK, nil
}

func DeleteNodes(c *gin.Context) (int, error) {
	nodeUUID := c.Query("node_uuid")

	bucketUUID, err := getBucketUUID(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	softDeleteValue, softDelete := c.GetQuery("soft_delete")
	if softDeleteValue == "false" {
		softDelete = false
	}

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx := database.NewTX(c)

	if err := storage.AuthorizeAccess(tx, models.Write, bucketUUID, user.ID); err != nil {
		return http.StatusInternalServerError, err
	}

	node, err := storage.GetNode(tx, nodeUUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	bucket, err := storage.GetBucket(tx, bucketUUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	path, err := storage.GetNodePath(tx, node, bucket.UUID, bucket.RootNode.UUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if softDelete {
		err = storage.DeleteNode(tx, node.UUID, softDelete)
	} else {
		err = storage.DeleteNodeRecursively(tx, &node)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		err = storage.DeleteNodeInFileSystem(path)
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	return http.StatusOK, nil
}

func RenameNode(c *gin.Context) (int, error) {
	nodeUUID := c.Query("node_uuid")
	newName := c.Query("new_name")

	bucketUUID, err := getBucketUUID(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return 0, nil
	}

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx := database.NewTX(c)

	bucket, err := storage.GetBucket(tx, bucketUUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := storage.AuthorizeAccess(tx, models.Write, bucket.UUID, user.ID); err != nil {
		return http.StatusInternalServerError, err
	}

	node, err := storage.GetNode(tx, nodeUUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	path, err := storage.GetNodePath(tx, node, bucket.UUID, bucket.RootNode.UUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	node.Name = newName

	if err := storage.UpdateNode(tx, &node, user.ID); err != nil {
		return http.StatusInternalServerError, err
	}

	if err := storage.RenameNodeInFileSystem(path, newName); err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	return http.StatusOK, nil
}

func GetBucket(c *gin.Context) (int, error) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx := database.NewTX(c)

	bucket, err := storage.GetUserBucket(tx, user.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	c.JSON(http.StatusOK, bucket)

	return http.StatusOK, nil
}

func DownloadNodes(c *gin.Context) (int, error) {
	nodeUUID := c.Query("node_uuid")

	bucketUUID, err := getBucketUUID(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx := database.NewTX(c)

	if err := storage.AuthorizeAccess(tx, models.ReadOnly, bucketUUID, user.ID); err != nil {
		return http.StatusInternalServerError, err
	}

	bucket, err := storage.GetBucket(tx, bucketUUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	path, err := storage.GetDownloadPath(tx, user.ID, nodeUUID, bucket.UUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	c.File(path)

	return http.StatusOK, nil
}

func UploadNode(c *gin.Context) (int, error) {
	parentUUID := c.Query("parent_uuid")

	bucketUUID, err := getBucketUUID(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx := database.NewTX(c)

	if err := storage.AuthorizeAccess(tx, models.Write, bucketUUID, user.ID); err != nil {
		return http.StatusInternalServerError, err
	}

	bucket, err := storage.GetBucket(tx, bucketUUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	nodeType := storage.DetectFileType(file.Filename)
	mime := storage.DetectFileMime(file)

	node := models.Node{
		ParentUUID: parentUUID,
		BucketUUID: bucket.UUID,
		Name:       file.Filename,
		Type:       nodeType,
		Mime:       &mime,
		Size:       &file.Size,
	}

	if err := storage.CreateNode(tx, user.ID, &node); err != nil {
		return http.StatusInternalServerError, err
	}

	path, err := storage.GetNodePath(tx, node, bucket.UUID, bucket.RootNode.UUID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := c.SaveUploadedFile(file, path); err != nil {
		return http.StatusInternalServerError, err
	}

	tx.Commit()
	return 0, nil
}
