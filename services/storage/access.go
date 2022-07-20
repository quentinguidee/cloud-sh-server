package storage

import (
	"errors"
	"fmt"
	"net/http"
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

	request := `
		INSERT INTO buckets_to_users(bucket_id, user_id, access_type)
		VALUES ($1, $2, $3)
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

func GetBucketUserAccess(tx *sqlx.Tx, bucketId int, userId int) (BucketAccess, IServiceError) {
	request := "SELECT * FROM buckets_to_users WHERE bucket_id = $1 AND user_id = $2"

	access := BucketAccess{
		BucketId: bucketId,
		UserId:   userId,
	}

	err := tx.Get(&access, request, access.BucketId, access.UserId)
	if err != nil {
		return BucketAccess{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return access, nil
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
