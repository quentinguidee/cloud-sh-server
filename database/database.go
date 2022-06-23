package database

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
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
	db := Database{instance: instance}
	db.Initialize()
	return db, nil
}

func (db *Database) Initialize() {
	_, _ = db.CreateUsersTable()
	_, _ = db.CreateGitHubAuthTable()
}

const KeyDatabase = "KEY_DATABASE"

func Middleware(database Database) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set(KeyDatabase, database)
		context.Next()
	}
}
