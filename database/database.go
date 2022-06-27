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

func GetDatabaseFromContext(c *gin.Context) Database {
	return c.MustGet(KeyDatabase).(Database)
}

func (db *Database) Initialize() {
	db.CreateUsersTable()
	db.CreateSessionsTable()
	db.CreateGithubAuthTable()

	db.CreateBucketsTable()
}

const KeyDatabase = "KEY_DATABASE"

func Middleware(database Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(KeyDatabase, database)
		c.Next()
	}
}
