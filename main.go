package main

import (
	"log"
	"os"
	"path/filepath"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/routes/auth"
	"self-hosted-cloud/server/routes/storage"
	"self-hosted-cloud/server/routes/user"
	"self-hosted-cloud/server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	if godotenv.Load() != nil {
		log.Fatal(".env Couldn't be loaded.")
	}

	dataPath := os.Getenv("DATA_PATH")

	err := os.Mkdir(dataPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err.Error())
		return
	}

	db, err := database.GetDatabase(filepath.Join(dataPath, "database.sqlite"))
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	router := gin.Default()

	router.Use(cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		MaxAge:         int(12 * time.Hour),
	}))
	router.Use(database.Middleware(db))
	router.Use(utils.ErrorMiddleware())

	auth.LoadRoutes(router)
	user.LoadRoutes(router)
	storage.LoadRoutes(router)

	err = router.Run(os.Getenv("SERVER_URI"))
	if err != nil {
		return
	}
}
