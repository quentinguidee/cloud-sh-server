package commands

import (
	"errors"
	"net/http"
	"os"
	"path"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models/storage"
)

type UpdateBucketNodeFilenameCommand struct {
	Database    Database
	Node        *Node
	NewFilename string

	oldFilename string
}

func (c UpdateBucketNodeFilenameCommand) Run() ICommandError {
	request := "UPDATE buckets_nodes SET name = ? WHERE uuid = ?"

	res, err := c.Database.Instance.Exec(request, c.NewFilename, c.Node.Uuid)
	if err != nil {
		err = errors.New("failed to update the node")
		return NewError(http.StatusInternalServerError, err)
	}

	count, err := res.RowsAffected()
	if err != nil && count == 0 {
		err = errors.New("couldn't find the node")
		return NewError(http.StatusNotFound, err)
	}

	c.oldFilename = c.Node.Name
	c.Node.Name = c.NewFilename

	return nil
}

func (c UpdateBucketNodeFilenameCommand) Revert() ICommandError {
	err := UpdateBucketNodeFilenameCommand{
		Database:    c.Database,
		Node:        c.Node,
		NewFilename: c.oldFilename,
		oldFilename: c.NewFilename,
	}.Run()

	if err != nil {
		return NewError(err.Code(), err.Error())
	}
	return nil
}

type UpdateBucketNodeFilenameInFileSystemCommand struct {
	CompletePath *string
	NewFilename  string

	oldPath string
	newPath string
}

func (c UpdateBucketNodeFilenameInFileSystemCommand) Run() ICommandError {
	c.oldPath = *c.CompletePath

	c.newPath = path.Dir(c.oldPath)
	c.newPath = path.Join(c.newPath, c.NewFilename)

	err := os.Rename(c.oldPath, c.newPath)
	if err != nil {
		err = errors.New("failed to rename this file")
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c UpdateBucketNodeFilenameInFileSystemCommand) Revert() ICommandError {
	err := UpdateBucketNodeFilenameInFileSystemCommand{
		CompletePath: &c.newPath,
		NewFilename:  path.Base(c.oldPath),
	}.Run()

	if err != nil {
		return NewError(err.Code(), err.Error())
	}
	return nil
}
