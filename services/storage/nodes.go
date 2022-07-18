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

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func GetBucketNode(tx *sqlx.Tx, uuid string) (Node, IServiceError) {
	request := "SELECT * FROM buckets_nodes WHERE uuid = $1"

	var node Node
	err := tx.Get(&node, request, uuid)
	if err != nil {
		return Node{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return node, nil
}

func GetBucketNodes(tx *sqlx.Tx, parentUuid string) ([]Node, IServiceError) {
	request := `
		SELECT nodes.*
		FROM buckets_nodes nodes, buckets_nodes_associations associations
		WHERE associations.from_node = $1
		  AND associations.to_node = nodes.uuid
	`

	var nodes []Node
	err := tx.Select(&nodes, request, parentUuid)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, err)
	}
	return nodes, nil
}

func GetBucketNodeParent(tx *sqlx.Tx, nodeUuid string) (Node, IServiceError) {
	request := `
		SELECT nodes.*
		FROM buckets_nodes nodes, buckets_nodes_associations associations
		WHERE associations.from_node = nodes.uuid
		  AND associations.to_node = $1
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

	for true {
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

	return "", NewServiceError(http.StatusInternalServerError, errors.New("unreachable code reached"))
}

func CreateBucketNode(tx *sqlx.Tx, name string, kind string, mime string, size int64) (Node, IServiceError) {
	node := Node{
		Uuid: uuid.NewString(),
		Name: name,
		Type: kind,
		Mime: NewNullableString(mime),
		Size: NewNullableInt64(size),
	}

	request := `
		INSERT INTO buckets_nodes(uuid, name, type, mime, size)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := tx.Exec(request,
		node.Uuid,
		node.Name,
		node.Type,
		node.Mime,
		node.Size)

	if err != nil {
		err := errors.New("error while creating node")
		return Node{}, NewServiceError(http.StatusInternalServerError, err)
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
	request := "DELETE FROM buckets_nodes_associations WHERE to_node = $1"

	_, err := tx.Exec(request, uuid)
	if err != nil {
		err = errors.New(fmt.Sprintf("error while deleting node association: %s", err.Error()))
		return NewServiceError(http.StatusInternalServerError, err)
	}

	request = "DELETE FROM buckets_to_node WHERE node_id = $1"

	_, err = tx.Exec(request, uuid)
	if err != nil {
		err = errors.New(fmt.Sprintf("error while deleting buckets_to_node association: %s", err))
		return NewServiceError(http.StatusInternalServerError, err)
	}

	request = "DELETE FROM buckets_nodes WHERE uuid = $1"

	_, err = tx.Exec(request, uuid)
	if err != nil {
		err = errors.New(fmt.Sprintf("error while deleting node: %s", err.Error()))
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

func UpdateBucketNode(tx *sqlx.Tx, name string, previousType string, uuid string) IServiceError {
	request := "UPDATE buckets_nodes SET name = $1, type = $2 WHERE uuid = $3"

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
	return nil
}

func RenameBucketNodeInFileSystem(path string, name string) IServiceError {
	directoryPath := filepath.Dir(path)
	newPath := filepath.Join(directoryPath, name)

	if path == newPath {
		return nil
	}

	err := os.Rename(path, newPath)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to rename this file from %s to %s", path, newPath))
		return NewServiceError(http.StatusInternalServerError, err)
	}
	return nil
}
