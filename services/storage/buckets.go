package storage

import (
	"errors"
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

	root, err := CreateBucketNode(tx, "root", "directory", "", 0, bucket.Id)
	if err != nil {
		return err
	}

	err = UpdateBucketRootNode(tx, bucket.Id, root.Uuid)
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
	request := "INSERT INTO buckets(name, type) VALUES (?, ?) RETURNING id"

	bucket := Bucket{
		Name: name,
		Type: kind,
	}

	err := tx.QueryRow(request, name, kind).Scan(&bucket.Id)
	if err != nil {
		err := errors.New("error while creating bucket")
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

func UpdateBucketRootNode(tx *sqlx.Tx, bucketId int, rootNodeUuid string) IServiceError {
	request := "UPDATE buckets SET root_node = ? WHERE id = ?"

	_, err := tx.Exec(request, rootNodeUuid, bucketId)
	if err != nil {
		err = errors.New("error while updating bucket")
		return NewServiceError(http.StatusInternalServerError, err)
	}
	return nil
}

func GetUserBucket(tx *sqlx.Tx, userId int) (Bucket, IServiceError) {
	request := `
		SELECT buckets.id, buckets.name, buckets.root_node, buckets.type
		FROM buckets, buckets_access access
		WHERE buckets.id = access.bucket_id
		  AND buckets.type = 'user_bucket'
		  AND access.user_id = ?
	`

	var bucket Bucket
	err := tx.QueryRow(request, userId).Scan(
		&bucket.Id,
		&bucket.Name,
		&bucket.RootNodeUuid,
		&bucket.Type)

	if err != nil {
		err = errors.New("error while getting user bucket")
		return Bucket{}, NewServiceError(http.StatusNotFound, err)
	}
	return bucket, nil
}

func GetBucketPath(bucketId int) string {
	return filepath.Join(os.Getenv("DATA_PATH"), "buckets", strconv.Itoa(bucketId))
}
