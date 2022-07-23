package auth

import (
	"crypto/rand"
	"fmt"
	. "self-hosted-cloud/server/models"

	"gorm.io/gorm"
)

func CreateSession(tx *gorm.DB, userID int) (Session, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return Session{}, err
	}

	session := Session{
		UserID: userID,
		Token:  fmt.Sprintf("%X", token),
	}

	return session, tx.Create(&session).Error
}

func DeleteSession(tx *gorm.DB, session *Session) error {
	return tx.Delete(session).Error
}
