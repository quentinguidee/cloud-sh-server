package main

import (
	"log"
	"self-hosted-cloud/server/auth"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	if godotenv.Load() != nil {
		log.Fatal(".env Couldn't be loaded.")
	}

	router := gin.Default()
	router.Use(cors.Default())

	auth.LoadRoutes(router)

	router.Run(":8080")
}
