package storage

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
	"self-hosted-cloud/server/models/types"
	. "self-hosted-cloud/server/services"
	"strconv"

	"github.com/jmoiron/sqlx"
)

func SetupDefaultBucket(tx *sqlx.Tx, userId int) IServiceError {
	bucket, err := CreateBucket(tx, "Main bucket", "user_bucket")
	if err != nil {
		return err
	}

	_, err = CreateRootNode(tx, userId, bucket.Id)
	if err != nil {
		return err
	}

	_, err = CreateBucketAccess(tx, bucket.Id, userId)
	if err != nil {
		return err
	}

	err = CreateBucketInFileSystem(bucket.Id)
	if err != nil {
		return err
	}

	return nil
}

func CreateBucket(tx *sqlx.Tx, name string, kind string) (Bucket, IServiceError) {
	query := "INSERT INTO buckets(name, type) VALUES ($1, $2) RETURNING id"

	bucket := Bucket{
		Name: name,
		Type: kind,
	}

	err := database.
		NewRequest(tx, query).
		QueryRow(name, kind).
		Scan(&bucket.Id).
		OnError("error while creating bucket")

	return bucket, err
}

func CreateBucketInFileSystem(bucketId int) IServiceError {
	err := os.MkdirAll(filepath.Join(os.Getenv("DATA_PATH"), "buckets", strconv.Itoa(bucketId)), os.ModePerm)
	if err != nil {
		err = errors.New("error while creating bucket in file system")
		return NewServiceError(http.StatusInternalServerError, err)
	}
	return nil
}

func GetBucket(tx *sqlx.Tx, bucketId int) (Bucket, IServiceError) {
	query := "SELECT * FROM detailed_buckets WHERE id = $1"

	var bucket Bucket

	err := database.
		NewRequest(tx, query).
		Get(&bucket, bucketId).
		OnError("failed to retrieve the bucket")

	return bucket, err
}

func GetUserBucket(tx *sqlx.Tx, userId int) (Bucket, IServiceError) {
	query := `
		SELECT buckets.*
		FROM detailed_buckets buckets INNER JOIN buckets_to_users access
		ON buckets.id = access.bucket_id
		WHERE buckets.type = 'user_bucket'
		  AND access.user_id = $1
	`

	var bucket Bucket

	err := database.
		NewRequest(tx, query).
		Get(&bucket, userId).
		OnError("error while getting user bucket")

	return bucket, err
}

func GetBucketPath(bucketId int) string {
	return filepath.Join(os.Getenv("DATA_PATH"), "buckets", strconv.Itoa(bucketId))
}

func BucketCanAcceptNodeOfSize(tx *sqlx.Tx, bucketId int, nodeSize int64) (bool, IServiceError) {
	if nodeSize == 0 {
		return true, nil
	}

	bucket, err := GetBucket(tx, bucketId)
	if err != nil {
		return false, err
	}

	if bucket.MaxSize == types.NewNullInt64() {
		return true, nil
	}

	return (bucket.Size + nodeSize) <= bucket.MaxSize.Int64, err
}
