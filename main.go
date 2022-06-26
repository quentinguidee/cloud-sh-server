package main

import (
	"log"
	"os"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/routes/auth"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	if godotenv.Load() != nil {
		log.Fatal(".env Couldn't be loaded.")
	}

	db, err := database.GetDatabase("database.sqlite")
	if err != nil {
		return
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.Use(database.Middleware(db))

	auth.LoadRoutes(router)
	err = router.Run("localhost:" + os.Getenv("SERVER_PORT"))
	if err != nil {
		return
	}
}
