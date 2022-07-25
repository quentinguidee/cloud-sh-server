package database

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"self-hosted-cloud/server/models"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var tables = []interface{}{
	&models.Server{},
	&models.User{},
	&models.Bucket{},
	&models.Node{},
	&models.Session{},
	&models.GithubUser{},
	&models.NodeUser{},
	&models.UserBucket{},
}

func OpenConnection() (*gorm.DB, error) {
	source := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASS"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_SSL"))

	return gorm.Open(postgres.Open(source), &gorm.Config{})
}

func GetDatabase() (*gorm.DB, error) {
	db, err := OpenConnection()
	if err != nil {
		return nil, errors.New("couldn't open connection to the database")
	}

	return db, Initialize(db)
}

func GetDatabaseFromContext(c *gin.Context) *gorm.DB {
	return c.MustGet(KeyDatabase).(*gorm.DB)
}

func Initialize(db *gorm.DB) error {
	err := db.AutoMigrate(tables...)

	if err != nil {
		return err
	}

	var server models.Server
	err = db.Take(&server, 1).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = db.Create(&models.Server{
			VersionMajor:    0,
			VersionMinor:    0,
			VersionPatch:    0,
			DatabaseVersion: Version,
		}).Error
		return err
	}
	return err
}

func HardReset(db *gorm.DB) error {
	for _, table := range tables {
		err := db.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}
	return Initialize(db)
}

const KeyDatabase = "KEY_DATABASE"
