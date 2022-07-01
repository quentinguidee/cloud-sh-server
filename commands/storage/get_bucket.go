package commands

import (
	"errors"
	"net/http"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/models/storage"
)

type GetUserBucketCommand struct {
	Database       Database
	User           *User
	ReturnedBucket *Bucket
}

func (c GetUserBucketCommand) Run() ICommandError {
	request := `
		SELECT buckets.id, buckets.name, buckets.root_node, buckets.type
		FROM buckets, buckets_access access
		WHERE buckets.id = access.bucket_id
		  AND buckets.type = 'user_bucket'
		  AND access.user_id = ?
	`

	err := c.Database.Instance.QueryRow(request, c.User.Id).Scan(
		&c.ReturnedBucket.Id,
		&c.ReturnedBucket.Name,
		&c.ReturnedBucket.RootNode,
		&c.ReturnedBucket.Type)

	if err != nil {
		err = errors.New("error while getting user bucket")
		return NewError(http.StatusNotFound, err)
	}
	return nil
}

func (c GetUserBucketCommand) Revert() ICommandError {
	return nil
}
