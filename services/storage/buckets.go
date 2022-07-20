package storage

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	request := "INSERT INTO buckets(name, type) VALUES ($1, $2) RETURNING id"

	bucket := Bucket{
		Name: name,
		Type: kind,
	}

	err := tx.QueryRow(request,
		name,
		kind,
	).Scan(&bucket.Id)

	if err != nil {
		err := fmt.Errorf("error while creating bucket: %s", err.Error())
		return Bucket{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return bucket, nil
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
	request := `
		SELECT buckets.*
		FROM buckets, buckets_to_users access
		WHERE buckets.id = access.bucket_id
		  AND buckets.type = 'user_bucket'
		  AND access.user_id = $1
	`

	var bucket Bucket
	err := tx.Get(&bucket, request, userId)
	if err != nil {
		err = errors.New("error while getting user bucket")
		return Bucket{}, NewServiceError(http.StatusNotFound, err)
	}
	return bucket, nil
}

func GetBucketRootNode(tx *sqlx.Tx, bucketId int) (Node, IServiceError) {
	request := "SELECT nodes.* FROM nodes WHERE nodes.bucket_id = $1 AND nodes.parent_uuid IS NULL"

	var node Node
	err := tx.Get(&node, request, bucketId)
	if err != nil {
		err = fmt.Errorf("failed to get bucket root node: %s", err)
		return Node{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return node, nil
}

func GetBucketPath(bucketId int) string {
	return filepath.Join(os.Getenv("DATA_PATH"), "buckets", strconv.Itoa(bucketId))
}
