package database

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	. "self-hosted-cloud/server/models"
)

func (db *Database) CreateSessionsTable() (sql.Result, error) {
	return db.instance.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id      INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			token   VARCHAR(255) UNIQUE
		)
	`)
}

func (db *Database) CreateSession(userId int) (Session, error) {
	request := "INSERT INTO sessions(user_id, token) VALUES (?, ?)"

	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return Session{}, err
	}

	session := Session{
		UserId: userId,
		Token:  fmt.Sprintf("%X", token),
	}

	_, err = db.instance.Exec(request, session.UserId, session.Token)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

func (db *Database) ValidateToken(token string, userId int) (bool, error) {
	request := "SELECT id, user_id, token FROM sessions WHERE token = ?"

	var session Session
	err := db.instance.QueryRow(request, token).Scan(&session.Id, &session.UserId, &session.Token)
	if err != nil {
		// Token not found
		return false, err
	}

	return userId == session.Id && token == session.Token, nil
}
