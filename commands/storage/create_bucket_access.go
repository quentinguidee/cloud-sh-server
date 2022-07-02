package commands

import (
	"errors"
	"net/http"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models/storage"
)

type CreateBucketAccess struct {
	Bucket   *Bucket
	Database Database
	UserId   int

	BucketAccess BucketAccess
}

func (c CreateBucketAccess) Run() ICommandError {
	c.BucketAccess = BucketAccess{
		BucketId:   c.Bucket.Id,
		UserId:     c.UserId,
		AccessType: "admin",
	}

	request := `
		INSERT INTO buckets_access(bucket_id, user_id, access_type)
		VALUES (?, ?, ?)
		RETURNING id
	`

	err := c.Database.Instance.QueryRow(request,
		c.BucketAccess.BucketId,
		c.BucketAccess.UserId,
		c.BucketAccess.AccessType,
	).Scan(&c.BucketAccess.Id)

	if err != nil {
		err = errors.New("error while creating bucket access")
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c CreateBucketAccess) Revert() ICommandError {
	request := "DELETE FROM buckets_access WHERE id = ?"

	_, err := c.Database.Instance.Exec(request, c.BucketAccess.Id)
	if err != nil {
		err = errors.New("error while deleting bucket access")
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}
