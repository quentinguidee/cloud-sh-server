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
	root, err := CreateBucketNode(tx, "root", "directory", "", 0)
	if err != nil {
		return err
	}

	bucket, err := CreateBucket(tx, "Main bucket", root.Uuid, "user_bucket")
	if err != nil {
		return err
	}

	err = CreateBucketToNodeAssociation(tx, bucket.Id, root.Uuid)
	if err != nil {
		return nil
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

func CreateBucket(tx *sqlx.Tx, name string, rootNode string, kind string) (Bucket, IServiceError) {
	request := "INSERT INTO buckets(name, root_node, type) VALUES ($1, $2, $3) RETURNING id"

	bucket := Bucket{
		Name:     name,
		RootNode: rootNode,
		Type:     kind,
	}

	err := tx.QueryRow(request,
		name,
		rootNode,
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
		FROM buckets, buckets_access access
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

func GetBucketPath(bucketId int) string {
	return filepath.Join(os.Getenv("DATA_PATH"), "buckets", strconv.Itoa(bucketId))
}
