package auth

import (
	. "self-hosted-cloud/server/models"

	"gorm.io/gorm"
)

func CreateGithubUser(tx *gorm.DB, userID int, username string) error {
	return tx.Create(&GithubUser{
		UserID:   userID,
		Username: username,
	}).Error
}

func GetGithubUser(tx *gorm.DB, username string) (User, error) {
	var user User
	err := tx.Preload("GithubUsers", "username = ?", username).Take(&user).Error
	return user, err
}
