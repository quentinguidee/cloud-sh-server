package commands

import (
	"net/http"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models/storage"
	"strings"
)

type GetNodeCommand struct {
	Database     Database
	Path         string
	Bucket       Bucket
	ReturnedNode *Node
}

func (c GetNodeCommand) Run() ICommandError {
	c.ReturnedNode.Id = c.Bucket.RootNode

	if len(c.Path) > 0 && c.Path[0] == '/' {
		c.Path = c.Path[1:]
	}

	if len(c.Path) == 0 {
		return nil
	}

	for _, filename := range strings.Split(c.Path, "/") {
		err := GetNodeInDirectory{
			Database:     c.Database,
			FromNode:     c.ReturnedNode.Id,
			Filename:     filename,
			ReturnedNode: c.ReturnedNode,
		}.Run()

		if err != nil {
			return NewError(err.Code(), err.Error())
		}
		if c.ReturnedNode.Filetype != "directory" {
			return nil
		}
	}

	return nil
}

func (c GetNodeCommand) Revert() ICommandError {
	return nil
}

type GetNodesCommand struct {
	Database      Database
	Path          string
	Bucket        Bucket
	ReturnedNodes *[]Node
}

func (c GetNodesCommand) Run() ICommandError {
	var node Node
	err := GetNodeCommand{
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
		FromNode:      node.Id,
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
	FromNode     int
	Filename     string
	ReturnedNode *Node
}

func (c GetNodeInDirectory) Run() ICommandError {
	request := `
		SELECT nodes.id, nodes.filename, nodes.filetype, nodes.bucket_id
		FROM buckets_nodes nodes, buckets_nodes_associations associations
		WHERE associations.from_node = ?
		  AND associations.to_node = nodes.id
		  AND nodes.filename = ?
	`

	err := c.Database.Instance.QueryRow(request, c.FromNode, c.Filename).Scan(
		&c.ReturnedNode.Id,
		&c.ReturnedNode.Filename,
		&c.ReturnedNode.Filetype,
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
	FromNode      int
	ReturnedNodes *[]Node
}

func (c GetNodesInDirectoryCommand) Run() ICommandError {
	request := `
		SELECT nodes.id, nodes.filename, nodes.filetype, nodes.bucket_id
		FROM buckets_nodes nodes, buckets_nodes_associations associations
		WHERE associations.from_node = ?
          AND associations.to_node = nodes.id
	`

	rows, err := c.Database.Instance.Query(request, c.FromNode)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}

	for rows.Next() {
		var node Node
		err := rows.Scan(
			&node.Id,
			&node.Filename,
			&node.Filetype,
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
