package storage

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/services"
	"strconv"

	"github.com/jmoiron/sqlx"
)

func SetupDefaultBucket(tx *sqlx.Tx, userId int) IServiceError {
	bucket, err := CreateBucket(tx, "Main bucket", "user_bucket")
	if err != nil {
		return err
	}

	_, err = CreateBucketRootNode(tx, userId, bucket.Id)
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

func GetUserBucket(tx *sqlx.Tx, userId int) (Bucket, IServiceError) {
	query := `
		SELECT buckets.*
		FROM buckets, buckets_to_users access
		WHERE buckets.id = access.bucket_id
		  AND buckets.type = 'user_bucket'
		  AND access.user_id = $1
	`

	var bucket Bucket

	err := database.
		NewRequest(tx, query).
		Get(&bucket, userId).
		OnError("error while getting user bucket")

	return bucket, err
}

func GetBucketRootNode(tx *sqlx.Tx, bucketId int) (Node, IServiceError) {
	query := "SELECT nodes.* FROM nodes WHERE nodes.bucket_id = $1 AND nodes.parent_uuid IS NULL"

	var node Node

	err := database.
		NewRequest(tx, query).
		Get(&node, bucketId).
		OnError("failed to get bucket root node")

	return node, err
}

func GetBucketPath(bucketId int) string {
	return filepath.Join(os.Getenv("DATA_PATH"), "buckets", strconv.Itoa(bucketId))
}
