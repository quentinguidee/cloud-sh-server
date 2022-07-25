package storage

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	. "self-hosted-cloud/server/models"

	"gorm.io/gorm"
)

func CreateUserBucket(tx *gorm.DB, userBucket *UserBucket) error {
	return tx.Create(userBucket).Error
}

func GetBucketUserAccess(tx *gorm.DB, bucketUUID uuid.UUID, userID int) (UserBucket, error) {
	bucketUser := UserBucket{
		Bucket: Bucket{UUID: bucketUUID},
		UserID: userID,
	}

	err := tx.Find(&bucketUser).Error

	return bucketUser, err
}

// AuthorizeAccess will return an error if the user cannot access this bucket ; nil error otherwise.
func AuthorizeAccess(tx *gorm.DB, desiredAccessType AccessType, bucketUUID uuid.UUID, userID int) error {
	maxAccessType, err := GetBucketUserAccessType(tx, bucketUUID, userID)
	if err != nil {
		return err
	}
	if maxAccessType < desiredAccessType {
		return errors.New("insufficient permissions")
	}
	return nil
}

func GetBucketUserAccessType(tx *gorm.DB, bucketUUID uuid.UUID, userID int) (AccessType, error) {
	access, err := GetBucketUserAccess(tx, bucketUUID, userID)
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
