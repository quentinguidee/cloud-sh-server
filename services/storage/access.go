package storage

import (
	"fmt"
	. "self-hosted-cloud/server/models"

	"gorm.io/gorm"
)

func CreateBucketUser(tx *gorm.DB, bucketID int, userID int) (BucketUser, error) {
	access := BucketUser{
		BucketID:   bucketID,
		UserID:     userID,
		AccessType: "admin",
	}

	err := tx.Create(&access).Error

	return access, err
}

func GetBucketUserAccess(tx *gorm.DB, bucketID int, userID int) (BucketUser, error) {
	bucketUser := BucketUser{
		BucketID: bucketID,
		UserID:   userID,
	}

	err := tx.Find(&bucketUser).Error

	return bucketUser, err
}

func GetBucketUserAccessType(tx *gorm.DB, bucketId int, userId int) (AccessType, error) {
	access, err := GetBucketUserAccess(tx, bucketId, userId)
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
