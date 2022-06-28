package database

import (
	"crypto/rand"
	"errors"
	"fmt"
	. "self-hosted-cloud/server/models"
)

func (db *Database) CreateSessionsTable() {
	_, _ = db.Instance.Exec(`
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

	_, err = db.Instance.Exec(request, session.UserId, session.Token)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

func (db *Database) CloseSession(session Session) error {
	request := "DELETE FROM sessions WHERE token = ? AND user_id = ?"

	res, err := db.Instance.Exec(request, session.Token, session.UserId)
	if err != nil {
		return err
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		return errors.New("the session doesn't exists")
	}

	return nil
}

func (db *Database) GetUserFromSession(token string) (User, error) {
	request := `
		SELECT users.id, users.username, users.name, users.profile_picture
		FROM users, sessions
		WHERE sessions.user_id = users.id
		  AND sessions.token = ?
	`

	var user User
	err := db.Instance.QueryRow(request, token).Scan(
		&user.Id,
		&user.Username,
		&user.Name,
		&user.ProfilePicture)

	if err != nil {
		return User{}, err
	}

	return user, err
}

func (db *Database) ValidateToken(token string, userId int) (bool, error) {
	request := "SELECT id, user_id, token FROM sessions WHERE token = ?"

	var session Session
	err := db.Instance.QueryRow(request, token).Scan(&session.Id, &session.UserId, &session.Token)
	if err != nil {
		// Token not found
		return false, err
	}

	return userId == session.Id && token == session.Token, nil
}
