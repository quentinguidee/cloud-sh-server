package commands

import (
	"errors"
	"net/http"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models/storage"
)

type CreateBucketNodeAssociationCommand struct {
	FromNode *Node
	ToNode   *Node
	Database Database
}

func (c CreateBucketNodeAssociationCommand) Run() ICommandError {
	request := `
		INSERT INTO buckets_nodes_associations(from_node, to_node)
		VALUES (?, ?)
	`

	_, err := c.Database.Instance.Exec(request, c.FromNode.Id, c.ToNode.Id)
	if err != nil {
		err = errors.New("error while creating bucket node association")
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c CreateBucketNodeAssociationCommand) Revert() ICommandError {
	request := `
		DELETE FROM buckets_nodes_associations
		WHERE from_node = ?
		  AND to_node = ?
	`

	_, err := c.Database.Instance.Exec(request, c.FromNode.Id, c.ToNode.Id)
	if err != nil {
		err = errors.New("error while deleting bucket node association")
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}
