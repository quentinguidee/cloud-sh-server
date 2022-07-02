package commands

import (
	"net/http"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models/storage"
	"strings"
)

type GetBucketNodeCommand struct {
	Database     Database
	Path         string
	Bucket       *Bucket
	ReturnedNode *Node
}

func (c GetBucketNodeCommand) Run() ICommandError {
	c.ReturnedNode.Uuid = c.Bucket.RootNodeUuid

	if len(c.Path) > 0 && c.Path[0] == '/' {
		c.Path = c.Path[1:]
	}

	if len(c.Path) == 0 {
		return nil
	}

	for _, filename := range strings.Split(c.Path, "/") {
		err := GetNodeInDirectory{
			Database:     c.Database,
			FromNodeUuid: c.ReturnedNode.Uuid,
			Name:         filename,
			ReturnedNode: c.ReturnedNode,
		}.Run()

		if err != nil {
			return NewError(err.Code(), err.Error())
		}
		if c.ReturnedNode.Type != "directory" {
			return nil
		}
	}

	return nil
}

func (c GetBucketNodeCommand) Revert() ICommandError {
	return nil
}

type GetNodesCommand struct {
	Database      Database
	Path          string
	Bucket        *Bucket
	ReturnedNodes *[]Node
}

func (c GetNodesCommand) Run() ICommandError {
	var node Node
	err := GetBucketNodeCommand{
		Database:     c.Database,
		Path:         c.Path,
		Bucket:       c.Bucket,
		ReturnedNode: &node,
	}.Run()

	if err != nil {
		return NewError(err.Code(), err.Error())
	}

	commandError := GetNodesInDirectoryCommand{
		Database:      c.Database,
		FromNodeUuid:  node.Uuid,
		ReturnedNodes: c.ReturnedNodes,
	}.Run()

	if commandError != nil {
		return NewError(commandError.Code(), commandError.Error())
	}
	return nil
}

func (c GetNodesCommand) Revert() ICommandError {
	return nil
}

type GetNodeInDirectory struct {
	Database     Database
	FromNodeUuid string
	Name         string
	ReturnedNode *Node
}

func (c GetNodeInDirectory) Run() ICommandError {
	request := `
		SELECT nodes.uuid, nodes.name, nodes.type, nodes.bucket_id
		FROM buckets_nodes nodes, buckets_nodes_associations associations
		WHERE associations.from_node = ?
		  AND associations.to_node = nodes.uuid
		  AND nodes.name = ?
	`

	err := c.Database.Instance.QueryRow(request, c.FromNodeUuid, c.Name).Scan(
		&c.ReturnedNode.Uuid,
		&c.ReturnedNode.Name,
		&c.ReturnedNode.Type,
		&c.ReturnedNode.BucketId)

	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c GetNodeInDirectory) Revert() ICommandError {
	return nil
}

type GetNodesInDirectoryCommand struct {
	Database      Database
	FromNodeUuid  string
	ReturnedNodes *[]Node
}

func (c GetNodesInDirectoryCommand) Run() ICommandError {
	request := `
		SELECT nodes.uuid, nodes.name, nodes.type, nodes.bucket_id
		FROM buckets_nodes nodes, buckets_nodes_associations associations
		WHERE associations.from_node = ?
          AND associations.to_node = nodes.uuid
	`

	rows, err := c.Database.Instance.Query(request, c.FromNodeUuid)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}

	for rows.Next() {
		var node Node
		err := rows.Scan(
			&node.Uuid,
			&node.Name,
			&node.Type,
			&node.BucketId)

		if err != nil {
			return NewError(http.StatusInternalServerError, err)
		}
		*c.ReturnedNodes = append(*c.ReturnedNodes, node)
	}
	return nil
}

func (c GetNodesInDirectoryCommand) Revert() ICommandError {
	return nil
}
