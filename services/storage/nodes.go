package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	. "self-hosted-cloud/server/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetNode(tx *gorm.DB, uuid string) (Node, error) {
	var node Node
	err := tx.Preload("Parent").Find(&node, "uuid = ?", uuid).Error
	return node, err
}

func GetNodes(tx *gorm.DB, parentUUID string) ([]Node, error) {
	var nodes []Node
	err := tx.Preload("Parent", "uuid = ?", parentUUID).Find(&nodes, "parent_uuid = ?", parentUUID).Error
	return nodes, err
}

func GetRecentFiles(tx *gorm.DB, userID int) ([]Node, error) {
	var nodes []Node
	err := tx.Preload("NodeUsers", func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID).Order("last_view_at DESC")
	}).Where("type <> ?", "directory").Find(&nodes).Error
	return nodes, err
}

func GetNodeParent(tx *gorm.DB, uuid string) (Node, error) {
	node, err := GetNode(tx, uuid)
	return *node.Parent, err
}

func GetNodePath(tx *gorm.DB, node Node, bucketId int, bucketRootNodeUuid string) (string, error) {
	var (
		i      = 50
		parent = node
		path   = node.Name
		err    error
	)

	for {
		parent, err = GetNodeParent(tx, parent.UUID)
		if parent.UUID == bucketRootNodeUuid {
			return filepath.Join(GetBucketPath(bucketId), path), nil
		}
		if err != nil {
			return "", err
		}
		if i == 0 {
			err := errors.New("max recursion level reached")
			return "", err
		}
		path = filepath.Join(parent.Name, path)
		i--
	}
}

func CreateRootNode(tx *gorm.DB, userID int, bucketID int) (Node, error) {
	node := Node{
		BucketID: bucketID,
		Name:     "root",
		Type:     "directory",
	}
	return CreateNode(tx, userID, node)
}

func CreateNode(tx *gorm.DB, userID int, node Node) (Node, error) {
	if node.Size != nil {
		accepted, err := BucketCanAcceptNodeOfSize(tx, node.BucketID, *node.Size)
		if err != nil {
			return Node{}, err
		}
		if !accepted {
			err := errors.New("the storage is full")
			// TODO: http.StatusForbidden
			return Node{}, err
		}
	}

	now := time.Now()

	node.UUID = uuid.NewString()
	node.NodeUsers = []NodeUser{{
		UserID:     userID,
		LastViewAt: &now,
		EditedAt:   &now,
	}}

	err := tx.Create(&node).Error
	if err != nil {
		return node, err
	}

	if node.Size != nil && *node.Size != 0 {
		var bucket Bucket
		err := tx.Take(&bucket, node.BucketID).Error
		if err != nil {
			return node, err
		}
		bucket.Size += *node.Size
		err = tx.Save(&bucket).Error
	}

	return node, err
}

func CreateNodeInFileSystem(kind string, path string, content string) error {
	_, err := os.Stat(path)
	if err == nil {
		err := errors.New("error while creating node in file system: this file already exists")
		return err
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
				return err
			}
		}
	}
	return err
}

func DeleteNode(tx *gorm.DB, uuid string) error {
	err := tx.Delete(&NodeUser{}, "node_uuid = ?", uuid).Error
	if err != nil {
		return err
	}

	var node Node
	err = tx.Clauses(clause.Returning{}).Delete(&node, "uuid = ?", uuid).Error
	if err != nil {
		return err
	}

	return tx.Model(&Bucket{ID: node.BucketID}).UpdateColumn("size", gorm.Expr("size - ?", *node.Size)).Error
}

func DeleteNodeRecursively(tx *gorm.DB, node *Node) error {
	if node.Type == "directory" {
		children, err := GetNodes(tx, node.UUID)
		if err != nil {
			return err
		}
		for _, node := range children {
			err := DeleteNodeRecursively(tx, &node)
			if err != nil {
				return err
			}
		}
	}

	err := DeleteNode(tx, node.UUID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteNodeInFileSystem(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		err = errors.New("error while deleting node in file system")
		return err
	}
	return nil
}

func UpdateNode(tx *gorm.DB, node *Node, userId int) error {
	err := tx.Save(&node).Error
	if err != nil {
		return err
	}
	return UpdateNodeLastEditionTimestamp(tx, userId, node.UUID)
}

func RenameNodeInFileSystem(path string, name string) error {
	directoryPath := filepath.Dir(path)
	newPath := filepath.Join(directoryPath, name)

	if path == newPath {
		return nil
	}

	err := os.Rename(path, newPath)
	if err != nil {
		err = fmt.Errorf("failed to rename this file from %s to %s", path, newPath)
		return err
	}
	return nil
}

func GetDownloadPath(tx *gorm.DB, userId int, uuid string, bucketId int) (string, error) {
	bucket, err := GetBucket(tx, bucketId)
	if err != nil {
		return "", err
	}

	node, err := GetNode(tx, uuid)
	if err != nil {
		return "", err
	}

	path, err := GetNodePath(tx, node, bucketId, bucket.RootNode.UUID)
	if err != nil {
		return "", err
	}

	err = UpdateNodeLastViewTimestamp(tx, userId, uuid)
	return path, err
}

func UpdateNodeLastViewTimestamp(tx *gorm.DB, userID int, uuid string) error {
	err := tx.Model(&NodeUser{
		UserID:   userID,
		NodeUUID: uuid,
	}).UpdateColumn("last_view_at", time.Now()).Error
	return err
}

func UpdateNodeLastEditionTimestamp(tx *gorm.DB, userID int, uuid string) error {
	err := tx.Model(&NodeUser{
		UserID:   userID,
		NodeUUID: uuid,
	}).UpdateColumn("edited_at", time.Now()).Error
	return err
}
