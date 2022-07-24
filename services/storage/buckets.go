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
	userBucket, err := CreateUserBucket(tx, UserBucket{
		Bucket: Bucket{
			UUID: uuid.New(),
			Name: "Main bucket",
			Type: "user_bucket",
		},
		UserID:     userID,
		AccessType: "admin",
	})
	if err != nil {
		return err
	}

	_, err = CreateRootNode(tx, userID, userBucket.Bucket.UUID)
	if err != nil {
		return err
	}

	return CreateBucketInFileSystem(userBucket.Bucket.UUID)
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
	err := tx.Take(&bucket, "uuid = ?", bucketUUID).Error
	if err != nil {
		return bucket, err
	}

	err = tx.Where("bucket_uuid = ?", bucketUUID).Where("parent_uuid IS NULL").Take(&bucket.RootNode).Error
	return bucket, err
}

func GetUserBucket(tx *gorm.DB, userID int) (Bucket, error) {
	var bucketUser UserBucket
	err := tx.Preload("Bucket").Take(&bucketUser, "user_id = ?", userID).Error
	if err != nil {
		return Bucket{}, err
	}
	return GetBucket(tx, bucketUser.Bucket.UUID)
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
