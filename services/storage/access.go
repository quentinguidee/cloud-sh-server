package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	. "self-hosted-cloud/server/models"
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

func GetBucketUserAccess(tx *sql.Tx, bucketId int, userId int) (BucketAccess, IServiceError) {
	access := BucketAccess{
		BucketId: bucketId,
		UserId:   userId,
	}

	request := "SELECT id, access_type FROM buckets_access WHERE bucket_id = ? AND user_id = ?"

	err := tx.QueryRow(request, access.BucketId, access.UserId).Scan(&access.Id, &access.AccessType)
	if err != nil {
		return BucketAccess{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return access, nil
}

func GetBucketUserAccessType(tx *sql.Tx, bucketId int, userId int) (AccessType, IServiceError) {
	access, serviceError := GetBucketUserAccess(tx, bucketId, userId)
	if serviceError != nil {
		return Denied, serviceError
	}

	switch access.AccessType {
	case "admin":
		return Full, nil
	default:
		err := errors.New(fmt.Sprintf("the access_type '%s' is not supported", access.AccessType))
		return Denied, NewServiceError(http.StatusInternalServerError, err)
	}
}
