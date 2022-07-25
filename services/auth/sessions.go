package auth

import (
	"crypto/rand"
	"fmt"
	. "self-hosted-cloud/server/models"

	"gorm.io/gorm"
)

func CreateSession(tx *gorm.DB, session *Session) error {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return err
	}

	session.Token = fmt.Sprintf("%X", token)

	return tx.Create(session).Error
}

func DeleteSession(tx *gorm.DB, session *Session) error {
	return tx.Delete(session).Error
}
