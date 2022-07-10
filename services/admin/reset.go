package admin

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/services"
)

func ResetServer(db *Database) IServiceError {
	log.Println("SERVER RESET…")
	paths := []string{
		"buckets",
		"database.sqlite",
	}

	for _, path := range paths {
		err := os.RemoveAll(filepath.Join(os.Getenv("DATA_PATH"), path))
		if err != nil {
			return NewServiceError(http.StatusInternalServerError, err)
		}
	}

	db.HardReset("database.sqlite")

	log.Println("SERVER RESET: DONE")
	return nil
}
