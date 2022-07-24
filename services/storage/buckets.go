package storage

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	. "self-hosted-cloud/server/models"
)

func SetupDefaultBucket(tx *gorm.DB, userID int) error {
	bucket, err := CreateBucket(tx, "Main bucket", "user_bucket")
	if err != nil {
		return err
	}

	_, err = CreateRootNode(tx, userID, bucket.UUID)
	if err != nil {
		return err
	}

	_, err = CreateBucketUser(tx, bucket.UUID, userID)
	if err != nil {
		return err
	}

	err = CreateBucketInFileSystem(bucket.UUID)
	if err != nil {
		return err
	}

	return nil
}

func CreateBucket(tx *gorm.DB, name string, kind string) (Bucket, error) {
	bucket := Bucket{
		UUID: uuid.New(),
		Name: name,
		Type: kind,
	}

	err := tx.Create(&bucket).Error

	return bucket, err
}

func CreateBucketInFileSystem(bucketUUID uuid.UUID) error {
	err := os.MkdirAll(filepath.Join(os.Getenv("DATA_PATH"), "buckets", bucketUUID.String()), os.ModePerm)
	if err != nil {
		err = fmt.Errorf("error while creating bucket in file system: %s", err)
		return err
	}
	return nil
}

func GetBucket(tx *gorm.DB, bucketUUID uuid.UUID) (Bucket, error) {
	var bucket Bucket
	err := tx.Take(&bucket, bucketUUID).Error
	if err != nil {
		return bucket, err
	}

	err = tx.Where("bucket_uuid = ?", bucketUUID).Where("parent_uuid IS NULL").Take(&bucket.RootNode).Error
	return bucket, err
}

func GetUserBucket(tx *gorm.DB, userID int) (Bucket, error) {
	var bucketUser BucketUser
	err := tx.Take(&bucketUser, "user_id = ?", userID).Error
	if err != nil {
		return Bucket{}, err
	}

	bucket, err := GetBucket(tx, bucketUser.BucketUUID)
	return bucket, err
}

func GetBucketPath(bucketUUID uuid.UUID) string {
	return filepath.Join(os.Getenv("DATA_PATH"), "buckets", bucketUUID.String())
}

func BucketCanAcceptNodeOfSize(tx *gorm.DB, bucketUUID uuid.UUID, nodeSize int64) (bool, error) {
	if nodeSize == 0 {
		return true, nil
	}

	bucket, err := GetBucket(tx, bucketUUID)
	if err != nil {
		return false, err
	}

	if bucket.MaxSize == nil {
		return true, nil
	}

	return (bucket.Size + nodeSize) <= *bucket.MaxSize, err
}
