package storage

import (
	"fmt"
	"github.com/google/uuid"
	. "self-hosted-cloud/server/models"

	"gorm.io/gorm"
)

func CreateUserBucket(tx *gorm.DB, userBucket UserBucket) (UserBucket, error) {
	err := tx.Create(&userBucket).Error
	return userBucket, err
}

func GetBucketUserAccess(tx *gorm.DB, bucketUUID uuid.UUID, userID int) (UserBucket, error) {
	bucketUser := UserBucket{
		Bucket: Bucket{UUID: bucketUUID},
		UserID: userID,
	}

	err := tx.Find(&bucketUser).Error

	return bucketUser, err
}

func GetBucketUserAccessType(tx *gorm.DB, bucketUUID uuid.UUID, userId int) (AccessType, error) {
	access, err := GetBucketUserAccess(tx, bucketUUID, userId)
	if err != nil {
		return Denied, err
	}

	switch access.AccessType {
	case "admin":
		return Full, nil
	default:
		err := fmt.Errorf("the access_type '%s' is not supported", access.AccessType)
		return Denied, err
	}
}
