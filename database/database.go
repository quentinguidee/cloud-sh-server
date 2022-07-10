package database

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Instance *sql.DB
}

func New(instance *sql.DB) Database {
	return Database{Instance: instance}
}

func OpenConnection(path string) (*sql.DB, error) {
	return sql.Open("sqlite3", filepath.Join(os.Getenv("DATA_PATH"), path))
}

func GetDatabase(path string) (Database, error) {
	instance, err := OpenConnection(path)
	if err != nil {
		return Database{}, errors.New("couldn't open connection to the database")
	}
	db := Database{Instance: instance}
	db.Initialize()
	return db, nil
}

func GetDatabaseFromContext(c *gin.Context) *Database {
	return c.MustGet(KeyDatabase).(*Database)
}

func (db *Database) Initialize() {
	db.CreateUsersTable()
	db.CreateSessionsTable()
	db.CreateGithubAuthTable()

	db.CreateBucketsTable()
}

func (db *Database) HardReset(path string) {
	instance, _ := OpenConnection(path)
	db.Instance = instance
	db.Initialize()
}

const KeyDatabase = "KEY_DATABASE"

func Middleware(database *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(KeyDatabase, database)
		c.Next()
	}
}
