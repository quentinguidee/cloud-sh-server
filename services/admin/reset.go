package admin

import (
	"log"
	"os"
	"path/filepath"
	"self-hosted-cloud/server/database"

	"gorm.io/gorm"
)

func ResetServer(tx *gorm.DB) error {
	log.Println("SERVER RESET…")
	paths := []string{
		"buckets",
		"demo.json",
	}

	for _, path := range paths {
		err := os.RemoveAll(filepath.Join(os.Getenv("DATA_PATH"), path))
		if err != nil {
			return err
		}
	}

	err := database.HardReset(tx)
	if err != nil {
		return err
	}

	log.Println("SERVER RESET: DONE")
	return nil
}
