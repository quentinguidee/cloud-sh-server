package storage

import (
	"fmt"
	"github.com/google/uuid"
	. "self-hosted-cloud/server/models"

	"gorm.io/gorm"
)

func CreateBucketUser(tx *gorm.DB, bucketUUID uuid.UUID, userID int) (BucketUser, error) {
	access := BucketUser{
		BucketUUID: bucketUUID,
		UserID:     userID,
		AccessType: "admin",
	}

	err := tx.Create(&access).Error

	return access, err
}

func GetBucketUserAccess(tx *gorm.DB, bucketUUID uuid.UUID, userID int) (BucketUser, error) {
	bucketUser := BucketUser{
		BucketUUID: bucketUUID,
		UserID:     userID,
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
