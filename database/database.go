package database

import (
	_ "embed"
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//go:embed create_db.sql
var createDatabaseRequest string

//go:embed reset_db.sql
var resetDatabaseRequest string

type Database struct {
	Instance *sqlx.DB
}

func New(instance *sqlx.DB) Database {
	return Database{Instance: instance}
}

func OpenConnection() (*sqlx.DB, error) {
	source := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASS"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_SSL"))

	return sqlx.Open("postgres", source)
}

func GetDatabase() (Database, error) {
	instance, err := OpenConnection()
	if err != nil {
		return Database{}, errors.New("couldn't open connection to the database")
	}
	db := Database{Instance: instance}

	var serverId int
	err = db.Instance.QueryRowx("SELECT id FROM servers WHERE id = 1").Scan(&serverId)
	if err != nil {
		// The servers table doesn't exist, so, the database is not initialized.
		err = db.Initialize()
		return db, err
	}
	err = db.Update()
	return db, err
}

func GetDatabaseFromContext(c *gin.Context) *Database {
	return c.MustGet(KeyDatabase).(*Database)
}

func (db *Database) Initialize() error {
	_, err := db.Instance.Exec(createDatabaseRequest)
	return err
}

func (db *Database) HardReset() error {
	_, err := db.Instance.Exec(resetDatabaseRequest)
	if err != nil {
		return err
	}
	return db.Initialize()
}

const KeyDatabase = "KEY_DATABASE"
