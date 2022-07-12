package database

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Instance *sqlx.DB
}

func New(instance *sqlx.DB) Database {
	return Database{Instance: instance}
}

func OpenConnection(path string) (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", filepath.Join(os.Getenv("DATA_PATH"), path))
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
