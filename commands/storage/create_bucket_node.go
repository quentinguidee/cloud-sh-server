package commands

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models/storage"
)

type CreateBucketNodeCommand struct {
	Node     *Node
	Bucket   *Bucket
	Database Database
}

func (c CreateBucketNodeCommand) Run() ICommandError {
	request := `
		INSERT INTO buckets_nodes(filename, filetype, bucket_id)
		VALUES (?, ?, ?)
		RETURNING id
	`

	c.Node.BucketId = c.Bucket.Id
	err := c.Database.Instance.QueryRow(request,
		c.Node.Filename,
		c.Node.Filetype,
		c.Node.BucketId,
	).Scan(&c.Node.Id)

	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c CreateBucketNodeCommand) Revert() ICommandError {
	request := "DELETE FROM buckets_nodes WHERE id = ?"

	_, err := c.Database.Instance.Exec(request, c.Node.Id)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

type CreateBucketNodeInFileSystemCommand struct {
	Node *Node
	Path string

	filePath string
}

func (c CreateBucketNodeInFileSystemCommand) Run() ICommandError {
	if len(c.Path) > 0 && c.Path[0] == '/' {
		c.Path = c.Path[1:]
	}

	if len(c.Path) > 0 {
		c.Path += "/"
	}

	c.filePath = fmt.Sprintf("%s/buckets/%d/%s%s", os.Getenv("DATA_PATH"), c.Node.BucketId, c.Path, c.Node.Filename)

	var err error
	switch c.Node.Filetype {
	case "directory":
		err = os.Mkdir(c.filePath, os.ModePerm)
	default:
		err = errors.New("this filetype is not supported")
	}

	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c CreateBucketNodeInFileSystemCommand) Revert() ICommandError {
	err := os.RemoveAll(c.filePath)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}
