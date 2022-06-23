package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	cors "github.com/rs/cors/wrapper/gin"
	"log"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/routes/auth"
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

	err = router.Run("localhost:8080")
	if err != nil {
		return
	}
}
