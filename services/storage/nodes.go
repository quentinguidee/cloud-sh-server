package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	. "self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/services"

	"github.com/google/uuid"
)

func GetBucketNode(tx *sql.Tx, uuid string) (Node, IServiceError) {
	request := "SELECT uuid, name, type, size, bucket_id FROM buckets_nodes WHERE uuid = ?"

	var node Node
	err := tx.QueryRow(request, uuid).Scan(&node.Uuid, &node.Name, &node.Type, &node.Size, &node.BucketId)
	if err != nil {
		return Node{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return node, nil
}

func GetBucketNodes(tx *sql.Tx, parentUuid string) ([]Node, IServiceError) {
	request := `
		SELECT uuid, name, type, size, bucket_id
		FROM buckets_nodes nodes, buckets_nodes_associations associations
		WHERE associations.from_node = ?
		  AND associations.to_node = nodes.uuid
	`

	rows, err := tx.Query(request, parentUuid)
	if err != nil {
		return nil, NewServiceError(http.StatusInternalServerError, err)
	}

	var nodes []Node
	for rows.Next() {
		var node Node
		err := rows.Scan(&node.Uuid, &node.Name, &node.Type, &node.Size, &node.BucketId)
		if err != nil {
			err := errors.New("failed to decode nodes")
			return nil, NewServiceError(http.StatusInternalServerError, err)
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func GetBucketNodeParent(tx *sql.Tx, nodeUuid string) (Node, IServiceError) {
	request := `
		SELECT nodes.uuid, name, type, size, bucket_id
		FROM buckets_nodes nodes, buckets_nodes_associations associations
		WHERE associations.from_node = nodes.uuid
		  AND associations.to_node = ?
	`

	var parent Node
	err := tx.QueryRow(request, nodeUuid).Scan(&parent.Uuid, &parent.Name, &parent.Type, &parent.Size, &parent.BucketId)
	if err != nil {
		return Node{}, NewServiceError(http.StatusInternalServerError, err)
	}

	return parent, nil
}

func GetBucketNodePath(tx *sql.Tx, node Node, bucketId int, bucketRootNodeUuid string) (string, IServiceError) {
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

func CreateBucketNode(tx *sql.Tx, name string, kind string, size int64, bucketId int) (Node, IServiceError) {
	node := Node{
		Uuid:     uuid.NewString(),
		Name:     name,
		Type:     kind,
		Size:     size,
		BucketId: bucketId,
	}

	request := `
		INSERT INTO buckets_nodes(uuid, name, type, size, bucket_id)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := tx.Exec(request,
		node.Uuid,
		node.Name,
		node.Type,
		node.Size,
		node.BucketId)

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

func DeleteBucketNode(tx *sql.Tx, uuid string) IServiceError {
	request := "DELETE FROM buckets_nodes WHERE uuid = ?"

	_, err := tx.Exec(request, uuid)
	if err != nil {
		err = errors.New("error while deleting node")
		return NewServiceError(http.StatusInternalServerError, err)
	}

	request = "DELETE FROM buckets_nodes_associations WHERE to_node = ?"

	_, err = tx.Exec(request, uuid)
	if err != nil {
		err = errors.New("error while deleting node association")
		return NewServiceError(http.StatusInternalServerError, err)
	}

	return nil
}

func DeleteBucketNodeRecursively(tx *sql.Tx, node *Node) IServiceError {
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

func UpdateBucketNode(tx *sql.Tx, name string, previousType string, uuid string) IServiceError {
	request := "UPDATE buckets_nodes SET name = ?, type = ? WHERE uuid = ?"

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
