package storage

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/models/types"
	. "self-hosted-cloud/server/services"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func GetBucketNode(tx *sqlx.Tx, uuid string) (Node, IServiceError) {
	query := "SELECT * FROM nodes WHERE uuid = $1"

	var node Node

	err := database.
		NewRequest(tx, query).
		Get(&node, uuid).
		OnError("error while getting bucket node")

	return node, err
}

func GetBucketNodes(tx *sqlx.Tx, parentUuid string) ([]Node, IServiceError) {
	query := `
		SELECT children.*
		FROM nodes parent, nodes children
		WHERE parent.uuid = children.parent_uuid
		  AND parent.uuid = $1
	`

	var nodes []Node

	err := database.
		NewRequest(tx, query).
		Select(&nodes, parentUuid).
		OnError("error while getting nodes")

	return nodes, err
}

func GetRecentFiles(tx *sqlx.Tx, userId int) ([]Node, IServiceError) {
	query := `
		SELECT nodes.*
		FROM nodes, nodes_to_users
		WHERE nodes_to_users.node_uuid = nodes.uuid
		  AND nodes_to_users.user_id = $1
		  AND nodes.type <> 'directory'
		ORDER BY nodes_to_users.last_view_timestamp DESC
	`

	var nodes []Node

	err := database.
		NewRequest(tx, query).
		Select(&nodes, userId).
		OnError("error while getting recent files")

	return nodes, err
}

func GetBucketNodeParent(tx *sqlx.Tx, uuid string) (Node, IServiceError) {
	query := `
		SELECT parent.*
		FROM nodes parent, nodes child
		WHERE child.parent_uuid = parent.uuid
		  AND child.uuid = $1
	`

	var parent Node

	err := database.
		NewRequest(tx, query).
		Get(&parent, uuid).
		OnError("error while getting node parent")

	return parent, err
}

func GetBucketNodePath(tx *sqlx.Tx, node Node, bucketId int, bucketRootNodeUuid string) (string, IServiceError) {
	var (
		i      = 50
		parent = node
		path   = node.Name
		err    IServiceError
	)

	for {
		parent, err = GetBucketNodeParent(tx, parent.Uuid)
		if parent.Uuid == bucketRootNodeUuid {
			return filepath.Join(GetBucketPath(bucketId), path), nil
		}
		if err != nil {
			return "", err
		}
		if i == 0 {
			err := errors.New("max recursion level reached")
			return "", NewServiceError(http.StatusInternalServerError, err)
		}
		path = filepath.Join(parent.Name, path)
		i--
	}
}

func CreateBucketRootNode(tx *sqlx.Tx, userId int, bucketId int) (Node, IServiceError) {
	return CreateBucketNode(tx, userId, NewNullString(), bucketId, "root", "directory", NewNullString(), NewNullInt64())
}

func CreateBucketNode(tx *sqlx.Tx, userId int, parentUuid NullableString, bucketId int, name string, kind string, mime NullableString, size NullableInt64) (Node, IServiceError) {
	if size.Valid == true {
		accepted, err := BucketCanAcceptNodeOfSize(tx, bucketId, size.Int64)
		if err != nil {
			return Node{}, err
		}
		if !accepted {
			err := errors.New("the storage is full")
			return Node{}, NewServiceError(http.StatusForbidden, err)
		}
	}

	node := Node{
		Uuid:       uuid.NewString(),
		ParentUuid: parentUuid,
		BucketId:   bucketId,
		Name:       name,
		Type:       kind,
		Mime:       mime,
		Size:       size,
	}

	query := `
		INSERT INTO nodes(uuid, parent_uuid, bucket_id, name, type, mime, size)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := database.
		NewRequest(tx, query).
		Exec(node.Uuid, node.ParentUuid, node.BucketId, node.Name, node.Type, node.Mime, node.Size).
		OnError("error while creating node")

	if err != nil {
		return node, err
	}

	query = `
		INSERT INTO nodes_to_users(user_id, node_uuid, last_view_timestamp, last_edition_timestamp)
		VALUES ($1, $2, $3, $4)
	`

	_, err = database.
		NewRequest(tx, query).
		Exec(userId, node.Uuid, time.Now(), time.Now()).
		OnError("error while creating node user specific data")

	if err != nil {
		return node, err
	}

	query = "UPDATE buckets SET size = size + $1 WHERE id = $2"

	_, err = database.
		NewRequest(tx, query).
		Exec(size, bucketId).
		OnError("failed to change the bucket size")

	return node, err
}

func CreateBucketNodeInFileSystem(kind string, path string, content string) IServiceError {
	_, err := os.Stat(path)
	if err == nil {
		err := errors.New("error while creating node in file system: this file already exists")
		return NewServiceError(http.StatusInternalServerError, err)
	}

	if kind == "directory" {
		err = os.Mkdir(path, os.ModePerm)
	} else {
		var file *os.File
		file, err = os.Create(path)
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		if len(content) > 0 {
			_, err := file.WriteString(content)
			if err != nil {
				return NewServiceError(http.StatusInternalServerError, err)
			}
		}
	}

	if err != nil {
		return NewServiceError(http.StatusInternalServerError, err)
	}
	return nil
}

func DeleteBucketNode(tx *sqlx.Tx, uuid string) IServiceError {
	query := "DELETE FROM nodes_to_users WHERE node_uuid = $1"

	_, err := database.
		NewRequest(tx, query).
		Exec(uuid).
		OnError("error while deleting node user specific data")

	if err != nil {
		return err
	}

	query = "DELETE FROM nodes WHERE uuid = $1 RETURNING size, bucket_id"

	var (
		size     int64
		bucketId int
	)

	err = database.
		NewRequest(tx, query).
		QueryRow(uuid).
		Scan(&size, &bucketId).
		OnError("error while deleting node")

	if err != nil {
		return err
	}

	query = "UPDATE buckets SET size = size - $1 WHERE id = $2"

	_, err = database.
		NewRequest(tx, query).
		Exec(size, bucketId).
		OnError("failed to update the bucket size")

	return err
}

func DeleteBucketNodeRecursively(tx *sqlx.Tx, node *Node) IServiceError {
	if node.Type == "directory" {
		children, err := GetBucketNodes(tx, node.Uuid)
		if err != nil {
			return err
		}
		for _, node := range children {
			err := DeleteBucketNodeRecursively(tx, &node)
			if err != nil {
				return err
			}
		}
	}

	err := DeleteBucketNode(tx, node.Uuid)
	if err != nil {
		return err
	}
	return nil
}

func DeleteBucketNodeInFileSystem(path string) IServiceError {
	err := os.RemoveAll(path)
	if err != nil {
		err = errors.New("error while deleting node in file system")
		return NewServiceError(http.StatusInternalServerError, err)
	}
	return nil
}

func UpdateBucketNode(tx *sqlx.Tx, name string, previousType string, uuid string, userId int) IServiceError {
	query := "UPDATE nodes SET name = $1, type = $2 WHERE uuid = $3"

	nodeType := previousType
	if previousType != "directory" {
		nodeType = DetectFileType(name)
	}

	res, serviceError := database.
		NewRequest(tx, query).
		Exec(name, nodeType, uuid).
		OnError("failed to update the node")

	if serviceError != nil {
		return serviceError
	}

	count, err := res.RowsAffected()
	if err != nil && count == 0 {
		err = errors.New("couldn't find the node")
		return NewServiceError(http.StatusNotFound, err)
	}

	return UpdateLastEditionTimestamp(tx, userId, uuid)
}

func RenameBucketNodeInFileSystem(path string, name string) IServiceError {
	directoryPath := filepath.Dir(path)
	newPath := filepath.Join(directoryPath, name)

	if path == newPath {
		return nil
	}

	err := os.Rename(path, newPath)
	if err != nil {
		err = fmt.Errorf("failed to rename this file from %s to %s", path, newPath)
		return NewServiceError(http.StatusInternalServerError, err)
	}
	return nil
}

func GetDownloadPath(tx *sqlx.Tx, userId int, uuid string, bucketId int) (string, IServiceError) {
	rootNode, serviceError := GetBucketRootNode(tx, bucketId)
	if serviceError != nil {
		return "", serviceError
	}

	node, serviceError := GetBucketNode(tx, uuid)
	if serviceError != nil {
		return "", serviceError
	}

	path, serviceError := GetBucketNodePath(tx, node, bucketId, rootNode.Uuid)
	if serviceError != nil {
		return "", serviceError
	}

	serviceError = UpdateLastViewTimestamp(tx, userId, uuid)
	return path, serviceError
}

func UpdateLastViewTimestamp(tx *sqlx.Tx, userId int, uuid string) IServiceError {
	query := `
		UPDATE nodes_to_users
		SET last_view_timestamp = $1
		WHERE node_uuid = $2
		  AND user_id = $3
	`

	res, serviceError := database.
		NewRequest(tx, query).
		Exec(time.Now(), uuid, userId).
		OnError("failed to update node user specific data")

	if serviceError != nil {
		return serviceError
	}

	count, err := res.RowsAffected()
	if err != nil && count == 0 {
		err = errors.New("couldn't find the node user specific data")
		return NewServiceError(http.StatusNotFound, err)
	}

	return nil
}

func UpdateLastEditionTimestamp(tx *sqlx.Tx, userId int, uuid string) IServiceError {
	query := `
		UPDATE nodes_to_users
		SET last_edition_timestamp = $1
		WHERE node_uuid = $2
		  AND user_id = $3
	`

	res, serviceError := database.
		NewRequest(tx, query).
		Exec(time.Now(), uuid, userId).
		OnError("failed to update node user specific data")

	if serviceError != nil {
		return serviceError
	}

	count, err := res.RowsAffected()
	if err != nil && count == 0 {
		err = errors.New("couldn't find the node user specific data")
		return NewServiceError(http.StatusNotFound, err)
	}

	return nil
}
