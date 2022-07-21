package storage

import (
	"fmt"
	"net/http"
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/services"

	"github.com/jmoiron/sqlx"
)

func CreateBucketAccess(tx *sqlx.Tx, bucketId int, userId int) (BucketAccess, IServiceError) {
	access := BucketAccess{
		BucketId:   bucketId,
		UserId:     userId,
		AccessType: "admin",
	}

	query := `
		INSERT INTO buckets_to_users(bucket_id, user_id, access_type)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := database.
		NewRequest(tx, query).
		QueryRow(access.BucketId, access.UserId, access.AccessType).
		Scan(&access.Id).
		OnError("error while creating bucket access")

	return access, err
}

func GetBucketUserAccess(tx *sqlx.Tx, bucketId int, userId int) (BucketAccess, IServiceError) {
	query := "SELECT * FROM buckets_to_users WHERE bucket_id = $1 AND user_id = $2"

	access := BucketAccess{
		BucketId: bucketId,
		UserId:   userId,
	}

	err := database.
		NewRequest(tx, query).
		Get(&access, access.BucketId, access.UserId).
		OnError("error while getting bucket user access")

	return access, err
}

func GetBucketUserAccessType(tx *sqlx.Tx, bucketId int, userId int) (AccessType, IServiceError) {
	access, serviceError := GetBucketUserAccess(tx, bucketId, userId)
	if serviceError != nil {
		return Denied, serviceError
	}

	switch access.AccessType {
	case "admin":
		return Full, nil
	default:
		err := fmt.Errorf("the access_type '%s' is not supported", access.AccessType)
		return Denied, NewServiceError(http.StatusInternalServerError, err)
	}
}
