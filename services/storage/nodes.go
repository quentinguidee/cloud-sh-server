package storage

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	. "self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/models/types"
	. "self-hosted-cloud/server/services"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func GetBucketNode(tx *sqlx.Tx, uuid string) (Node, IServiceError) {
	request := "SELECT * FROM nodes WHERE uuid = $1"

	var node Node
	err := tx.Get(&node, request, uuid)
	if err != nil {
		return Node{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return node, nil
}

func GetBucketNodes(tx *sqlx.Tx, parentUuid string) ([]Node, IServiceError) {
	request := `
		SELECT children.*
		FROM nodes parent, nodes children
		WHERE parent.uuid = children.parent_uuid
		  AND parent.uuid = $1
	`

	var nodes []Node
	err := tx.Select(&nodes, request, parentUuid)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, err)
	}
	return nodes, nil
}

func GetRecentFiles(tx *sqlx.Tx, userId int) ([]Node, IServiceError) {
	request := `
		SELECT nodes.*
		FROM nodes, nodes_to_users
		WHERE nodes_to_users.node_uuid = nodes.uuid
		  AND nodes_to_users.user_id = $1
		  AND nodes.type <> 'directory'
		ORDER BY nodes_to_users.last_view_timestamp DESC
	`

	var nodes []Node
	err := tx.Select(&nodes, request, userId)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, err)
	}
	return nodes, nil
}

func GetBucketNodeParent(tx *sqlx.Tx, nodeUuid string) (Node, IServiceError) {
	request := `
		SELECT parent.*
		FROM nodes parent, nodes child
		WHERE child.parent_uuid = parent.uuid
		  AND child.uuid = $1
	`

	var parent Node
	err := tx.Get(&parent, request, nodeUuid)
	if err != nil {
		return Node{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return parent, nil
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

func CreateBucketRootNode(tx *sqlx.Tx, userId int) (Node, IServiceError) {
	return CreateBucketNode(tx, userId, NewNullString(), "root", "directory", NewNullString(), NewNullInt64())
}

func CreateBucketNode(tx *sqlx.Tx, userId int, parentUuid NullableString, name string, kind string, mime NullableString, size NullableInt64) (Node, IServiceError) {
	node := Node{
		Uuid:       uuid.NewString(),
		ParentUuid: parentUuid,
		Name:       name,
		Type:       kind,
		Mime:       mime,
		Size:       size,
	}

	request := `
		INSERT INTO nodes(uuid, parent_uuid, name, type, mime, size)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := tx.Exec(request,
		node.Uuid,
		node.ParentUuid,
		node.Name,
		node.Type,
		node.Mime,
		node.Size)

	if err != nil {
		err := errors.New("error while creating node")
		return Node{}, NewServiceError(http.StatusInternalServerError, err)
	}

	request = `
		INSERT INTO nodes_to_users(user_id, node_uuid, last_view_timestamp, last_edition_timestamp)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(request,
		userId,
		node.Uuid,
		time.Now(),
		time.Now())

	if err != nil {
		err := fmt.Errorf("error while creating node user specific data: %s", err)
		return node, NewServiceError(http.StatusInternalServerError, err)
	}
	return node, nil
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
	request := "DELETE FROM buckets_to_nodes WHERE node_uuid = $1"

	_, err := tx.Exec(request, uuid)
	if err != nil {
		err = fmt.Errorf("error while deleting buckets_to_node association: %s", err)
		return NewServiceError(http.StatusInternalServerError, err)
	}

	request = "DELETE FROM nodes_to_users WHERE node_uuid = $1"

	_, err = tx.Exec(request, uuid)
	if err != nil {
		err = fmt.Errorf("error while deleting node user specific data: %s", err)
		return NewServiceError(http.StatusInternalServerError, err)
	}

	request = "DELETE FROM nodes WHERE uuid = $1"

	_, err = tx.Exec(request, uuid)
	if err != nil {
		err = fmt.Errorf("error while deleting node: %s", err.Error())
		return NewServiceError(http.StatusInternalServerError, err)
	}

	return nil
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
	request := "UPDATE nodes SET name = $1, type = $2 WHERE uuid = $3"

	nodeType := previousType
	if previousType != "directory" {
		nodeType = DetectFileType(name)
	}

	res, err := tx.Exec(request, name, nodeType, uuid)
	if err != nil {
		err = errors.New("failed to update the node")
		return NewServiceError(http.StatusInternalServerError, err)
	}

	count, err := res.RowsAffected()
	if err != nil && count == 0 {
		err = errors.New("couldn't find the node")
		return NewServiceError(http.StatusNotFound, err)
	}

	serviceError := UpdateLastEditionTimestamp(tx, userId, uuid)
	return serviceError
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

func GetDownloadPath(tx *sqlx.Tx, userId int, uuid string, bucketId int, bucketRootNode string) (string, IServiceError) {
	node, serviceError := GetBucketNode(tx, uuid)
	if serviceError != nil {
		return "", serviceError
	}

	path, serviceError := GetBucketNodePath(tx, node, bucketId, bucketRootNode)
	if serviceError != nil {
		return "", serviceError
	}

	serviceError = UpdateLastViewTimestamp(tx, userId, uuid)
	return path, serviceError
}

func UpdateLastViewTimestamp(tx *sqlx.Tx, userId int, uuid string) IServiceError {
	request := `
		UPDATE nodes_to_users
		SET last_view_timestamp = $1
		WHERE node_uuid = $2
		  AND user_id = $3
	`

	res, err := tx.Exec(request, time.Now(), uuid, userId)
	if err != nil {
		err = fmt.Errorf("failed to update node user specific data: %s", err)
		return NewServiceError(http.StatusInternalServerError, err)
	}

	count, err := res.RowsAffected()
	if err != nil && count == 0 {
		err = errors.New("couldn't find the node user specific data")
		return NewServiceError(http.StatusNotFound, err)
	}

	return nil
}

func UpdateLastEditionTimestamp(tx *sqlx.Tx, userId int, uuid string) IServiceError {
	request := `
		UPDATE nodes_to_users
		SET last_edition_timestamp = $1
		WHERE node_uuid = $2
		  AND user_id = $3
	`

	res, err := tx.Exec(request, time.Now(), uuid, userId)
	if err != nil {
		err = fmt.Errorf("failed to update node user specific data: %s", err)
		return NewServiceError(http.StatusInternalServerError, err)
	}

	count, err := res.RowsAffected()
	if err != nil && count == 0 {
		err = errors.New("couldn't find the node user specific data")
		return NewServiceError(http.StatusNotFound, err)
	}

	return nil
}
