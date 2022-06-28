package commands

import (
	"net/http"
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
		VALUES ('/', 'directory', ?)
		RETURNING id
	`

	c.Node.BucketId = c.Bucket.Id
	err := c.Database.Instance.QueryRow(request, c.Node.BucketId).Scan(&c.Node.Id)
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
