package commands

import (
	"net/http"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models/storage"
)

type UpdateBucketRootNodeCommand struct {
	Bucket   *Bucket
	Database Database
	Node     *Node

	oldRootNode int
}

func (c UpdateBucketRootNodeCommand) Run() ICommandError {
	request := "UPDATE buckets SET root_node = ? WHERE id = ?"

	_, err := c.Database.Instance.Exec(request, c.Node.Id, c.Bucket.Id)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	c.oldRootNode = c.Bucket.RootNode
	c.Bucket.RootNode = c.Node.Id
	return nil
}

func (c UpdateBucketRootNodeCommand) Revert() ICommandError {
	request := "UPDATE buckets SET root_node = ? WHERE id = ?"

	_, err := c.Database.Instance.Exec(request, c.oldRootNode, c.Bucket.Id)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}
