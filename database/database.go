package database

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	. "self-hosted-cloud/server/models"
)

type Database struct {
	instance *sql.DB
}

func New(instance *sql.DB) Database {
	return Database{instance: instance}
}

func GetDatabase(path string) (Database, error) {
	instance, err := sql.Open("sqlite3", path)
	if err != nil {
		return Database{}, errors.New("couldn't open connection to the database")
	}
	return Database{instance: instance}, nil
}

const KeyDatabase = "KEY_DATABASE"

func Middleware(database Database) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set(KeyDatabase, database)
		context.Next()
	}
}

func (db *Database) GetUser(username string) (User, error) {
	statement, err := db.instance.Prepare("SELECT id, username, name FROM users WHERE username = ?")
	if err != nil {
		return User{}, errors.New("failed to prepare statement")
	}

	var user User
	err = statement.QueryRow(username).Scan(&user.Id, &user.Username, &user.Name)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *Database) CreateUser(user User) error {
	statement, err := db.instance.Prepare("INSERT INTO users(username, name) VALUES (?, ?)")
	if err != nil {
		return errors.New("failed to prepare statement")
	}

	_, err = statement.Exec(user.Username, user.Name)
	if err != nil {
		return err
	}

	return nil
}
