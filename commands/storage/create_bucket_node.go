package commands

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models/storage"

	"github.com/google/uuid"
)

type CreateBucketNodeCommand struct {
	Node     *Node
	Bucket   *Bucket
	Database Database
}

func (c CreateBucketNodeCommand) Run() ICommandError {
	c.Node.Uuid = uuid.NewString()

	request := `
		INSERT INTO buckets_nodes(uuid, filename, filetype, bucket_id)
		VALUES (?, ?, ?, ?)
	`

	c.Node.BucketId = c.Bucket.Id
	_, err := c.Database.Instance.Exec(request,
		c.Node.Uuid,
		c.Node.Filename,
		c.Node.Filetype,
		c.Node.BucketId)

	if err != nil {
		err := errors.New("error while creating node")
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c CreateBucketNodeCommand) Revert() ICommandError {
	request := "DELETE FROM buckets_nodes WHERE uuid = ?"

	_, err := c.Database.Instance.Exec(request, c.Node.Uuid)
	if err != nil {
		err := errors.New("error while deleting node")
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

type CreateBucketNodeInFileSystemCommand struct {
	Node    *Node
	Path    string
	Content string

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

	_, err := os.Stat(c.filePath)
	if err == nil {
		err := errors.New("error while creating node in file system: this file already exists")
		return NewError(http.StatusInternalServerError, err)
	}

	switch c.Node.Filetype {
	case "directory":
		err = os.Mkdir(c.filePath, os.ModePerm)
	case "file":
		var file *os.File
		file, err = os.Create(c.filePath)
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		if len(c.Content) > 0 {
			_, err := file.WriteString(c.Content)
			if err != nil {
				return NewError(http.StatusInternalServerError, err)
			}
		}
	default:
		err = errors.New(fmt.Sprintf("the filetype '%s' is not supported", c.Node.Filetype))
	}

	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c CreateBucketNodeInFileSystemCommand) Revert() ICommandError {
	err := os.RemoveAll(c.filePath)
	if err != nil {
		err = errors.New("error while deleting node in file system")
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}
