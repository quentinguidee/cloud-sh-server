package storage

import (
	"database/sql"
	"errors"
	"net/http"
	. "self-hosted-cloud/server/models/storage"
	. "self-hosted-cloud/server/services"
)

func CreateBucketAccess(tx *sql.Tx, bucketId int, userId int) (BucketAccess, IServiceError) {
	access := BucketAccess{
		BucketId:   bucketId,
		UserId:     userId,
		AccessType: "admin",
	}

	request := `
		INSERT INTO buckets_access(bucket_id, user_id, access_type)
		VALUES (?, ?, ?)
		RETURNING id
	`

	err := tx.QueryRow(request,
		access.BucketId,
		access.UserId,
		access.AccessType,
	).Scan(&access.Id)

	if err != nil {
		err = errors.New("error while creating bucket access")
		return BucketAccess{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return access, nil
}
