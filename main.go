package main

import (
	"fmt"
	"log"
	"os"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/routes/auth"
	"self-hosted-cloud/server/routes/storage"
	"self-hosted-cloud/server/routes/user"
	"self-hosted-cloud/server/utils"

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

	db, err := database.GetDatabase(fmt.Sprintf("%s/database.sqlite", dataPath))
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	router := gin.Default()

	router.Use(cors.AllowAll())
	router.Use(database.Middleware(db))
	router.Use(utils.ErrorMiddleware())

	auth.LoadRoutes(router)
	user.LoadRoutes(router)
	storage.LoadRoutes(router)

	err = router.Run("localhost:" + os.Getenv("SERVER_PORT"))
	if err != nil {
		return
	}
}
