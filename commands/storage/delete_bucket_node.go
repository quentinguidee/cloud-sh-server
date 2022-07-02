package commands

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models/storage"
	"strconv"
)

type DeleteBucketNodeRecursivelyCommand struct {
	Node     *Node
	Path     string
	Database Database
}

func (c DeleteBucketNodeRecursivelyCommand) Run() ICommandError {
	var nodes []Node
	err := GetNodesInDirectoryCommand{
		Database:      c.Database,
		FromNodeUuid:  c.Node.Uuid,
		ReturnedNodes: &nodes,
	}.Run()

	if err != nil && err.Error() != sql.ErrNoRows {
		err := errors.New("error while deleting nodes")
		return NewError(http.StatusInternalServerError, err)
	}

	for _, node := range nodes {
		var err ICommandError

		path := filepath.Join(c.Path, node.Name)

		switch node.Type {
		case "directory":
			err = DeleteBucketNodeRecursivelyCommand{
				Node:     &node,
				Path:     path,
				Database: c.Database,
			}.Run()
		default:
			err = NewTransaction([]Command{
				DeleteBucketNodeCommand{
					Node:     &node,
					Database: c.Database,
				},
				DeleteBucketNodeInFileSystemCommand{
					Node:     &node,
					Path:     path,
					Database: c.Database,
				},
			}).Try()
		}

		if err != nil {
			return NewError(http.StatusInternalServerError, err.Error())
		}
	}

	transactionError := NewTransaction([]Command{
		DeleteBucketNodeCommand{
			Node:     c.Node,
			Database: c.Database,
		},
		DeleteBucketNodeInFileSystemCommand{
			Node:     c.Node,
			Path:     c.Path,
			Database: c.Database,
		},
	}).Try()

	if transactionError != nil {
		return NewError(transactionError.Code(), transactionError.Error())
	}
	return nil
}

func (c DeleteBucketNodeRecursivelyCommand) Revert() ICommandError {
	// TODO: Revert file deletion
	return nil
}

type DeleteBucketNodeCommand struct {
	Node     *Node
	Database Database
}

func (c DeleteBucketNodeCommand) Run() ICommandError {
	request := `
		BEGIN TRANSACTION;
		DELETE FROM buckets_nodes WHERE uuid = ?;
		DELETE FROM buckets_nodes_associations WHERE to_node = ?;
		COMMIT TRANSACTION;
	`

	_, err := c.Database.Instance.Exec(request, c.Node.Uuid, c.Node.Uuid)
	if err != nil {
		err = errors.New("error while deleting node")
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c DeleteBucketNodeCommand) Revert() ICommandError {
	// TODO: Revert file deletion
	return nil
}

type DeleteBucketNodeInFileSystemCommand struct {
	Node     *Node
	Path     string
	Database Database
}

func (c DeleteBucketNodeInFileSystemCommand) Run() ICommandError {
	if len(c.Path) > 0 && c.Path[0] == '/' {
		c.Path = c.Path[1:]
	}

	c.Path = filepath.Join(os.Getenv("DATA_PATH"), "buckets", strconv.Itoa(c.Node.BucketId), c.Path)
	err := os.RemoveAll(c.Path)
	if err != nil {
		err = errors.New("error while deleting node in file system")
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c DeleteBucketNodeInFileSystemCommand) Revert() ICommandError {
	// TODO: Revert file deletion
	return nil
}
