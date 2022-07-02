package commands

import (
	"errors"
	"net/http"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models/storage"
)

type UpdateBucketRootNodeCommand struct {
	Bucket   *Bucket
	Database Database
	Node     *Node

	oldRootNode string
}

func (c UpdateBucketRootNodeCommand) Run() ICommandError {
	request := "UPDATE buckets SET root_node = ? WHERE id = ?"

	_, err := c.Database.Instance.Exec(request, c.Node.Uuid, c.Bucket.Id)
	if err != nil {
		err = errors.New("error while updating bucket")
		return NewError(http.StatusInternalServerError, err)
	}
	c.oldRootNode = c.Bucket.RootNodeUuid
	c.Bucket.RootNodeUuid = c.Node.Uuid
	return nil
}

func (c UpdateBucketRootNodeCommand) Revert() ICommandError {
	request := "UPDATE buckets SET root_node = ? WHERE id = ?"

	_, err := c.Database.Instance.Exec(request, c.oldRootNode, c.Bucket.Id)
	if err != nil {
		err = errors.New("error while updating bucket")
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}
