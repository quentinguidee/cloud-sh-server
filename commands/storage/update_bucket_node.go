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

type UpdateBucketNodeCommand struct {
	Database    Database
	Node        *Node
	NewFilename string

	oldFilename string
}

func (c UpdateBucketNodeCommand) Run() ICommandError {
	request := "UPDATE buckets_nodes SET filename = ? WHERE id = ?"

	res, err := c.Database.Instance.Exec(request, c.NewFilename, c.Node.Id)
	if err != nil {
		err = errors.New("failed to update the node")
		return NewError(http.StatusInternalServerError, err)
	}

	count, err := res.RowsAffected()
	if err != nil && count == 0 {
		err = errors.New("couldn't find the node")
		return NewError(http.StatusNotFound, err)
	}

	c.oldFilename = c.Node.Filename
	c.Node.Filename = c.NewFilename

	return nil
}

func (c UpdateBucketNodeCommand) Revert() ICommandError {
	err := UpdateBucketNodeCommand{
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

type UpdateBucketNodeInFileSystemCommand struct {
	CompletePath *string
	NewFilename  string

	oldPath string
	newPath string
}

func (c UpdateBucketNodeInFileSystemCommand) Run() ICommandError {
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

func (c UpdateBucketNodeInFileSystemCommand) Revert() ICommandError {
	err := UpdateBucketNodeInFileSystemCommand{
		CompletePath: &c.newPath,
		NewFilename:  path.Base(c.oldPath),
	}.Run()

	if err != nil {
		return NewError(err.Code(), err.Error())
	}
	return nil
}
