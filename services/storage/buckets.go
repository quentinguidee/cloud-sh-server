package storage

import (
	"errors"
	"os"
	"path/filepath"
	. "self-hosted-cloud/server/models"
	"strconv"

	"gorm.io/gorm"
)

func SetupDefaultBucket(tx *gorm.DB, userID int) error {
	bucket, err := CreateBucket(tx, "Main bucket", "user_bucket")
	if err != nil {
		return err
	}

	_, err = CreateRootNode(tx, userID, bucket.ID)
	if err != nil {
		return err
	}

	_, err = CreateBucketUser(tx, bucket.ID, userID)
	if err != nil {
		return err
	}

	err = CreateBucketInFileSystem(bucket.ID)
	if err != nil {
		return err
	}

	return nil
}

func CreateBucket(tx *gorm.DB, name string, kind string) (Bucket, error) {
	bucket := Bucket{
		Name: name,
		Type: kind,
	}

	err := tx.Create(&bucket).Error

	return bucket, err
}

func CreateBucketInFileSystem(bucketId int) error {
	err := os.MkdirAll(filepath.Join(os.Getenv("DATA_PATH"), "buckets", strconv.Itoa(bucketId)), os.ModePerm)
	if err != nil {
		err = errors.New("error while creating bucket in file system")
		return err
	}
	return nil
}

func GetBucket(tx *gorm.DB, bucketID int) (Bucket, error) {
	var bucket Bucket
	err := tx.Take(&bucket, bucketID).Error
	if err != nil {
		return bucket, err
	}

	err = tx.Where("bucket_id = ?", bucketID).Where("parent_uuid IS NULL").Take(&bucket.RootNode).Error
	return bucket, err
}

func GetUserBucket(tx *gorm.DB, userID int) (Bucket, error) {
	var bucketUser BucketUser
	err := tx.Take(&bucketUser, "user_id = ?", userID).Error
	if err != nil {
		return Bucket{}, err
	}

	bucket, err := GetBucket(tx, bucketUser.BucketID)
	return bucket, err
}

func GetBucketPath(bucketId int) string {
	return filepath.Join(os.Getenv("DATA_PATH"), "buckets", strconv.Itoa(bucketId))
}

func BucketCanAcceptNodeOfSize(tx *gorm.DB, bucketId int, nodeSize int64) (bool, error) {
	if nodeSize == 0 {
		return true, nil
	}

	bucket, err := GetBucket(tx, bucketId)
	if err != nil {
		return false, err
	}

	if bucket.MaxSize == nil {
		return true, nil
	}

	return (bucket.Size + nodeSize) <= *bucket.MaxSize, err
}
