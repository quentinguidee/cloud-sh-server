package auth

import (
	. "self-hosted-cloud/server/models"

	"gorm.io/gorm"
)

func CreateUser(tx *gorm.DB, user *User) error {
	return tx.Create(user).Error
}

func GetUser(tx *gorm.DB, username string) (User, error) {
	var user User
	err := tx.Where(&User{Username: username}).Find(&user).Error
	return user, err
}

func GetUserFromToken(tx *gorm.DB, token string) (User, error) {
	var user User
	err := tx.Preload("Sessions", "token = ?", token).Find(&user).Error
	return user, err
}

func GetUsersByRole(tx *gorm.DB, role string) ([]User, error) {
	var users []User
	err := tx.Where(&User{Role: &role}).Find(&users).Error
	return users, err
}
