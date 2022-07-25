package main

import (
	"log"
	"os"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/middlewares"
	"self-hosted-cloud/server/routes"
	adminservice "self-hosted-cloud/server/services/admin"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(".env Couldn't be loaded: %s", err.Error())
	}

	err := os.Mkdir(os.Getenv("DATA_PATH"), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("failed to create DATA_PATH: %s", err.Error())
		return
	}

	db, err := database.GetDatabase()
	if err != nil {
		log.Fatalf("failed to get the database: %s", err.Error())
		return
	}

	appIsInDemoMode, err := adminservice.AppIsInDemoMode()
	if err != nil {
		log.Fatal(err.Error())
	}
	if appIsInDemoMode {
		err := adminservice.StartDemoMode(db)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	}

	router := gin.Default()
	router.Use(cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		MaxAge:         int(12 * time.Hour),
	}))
	router.Use(middlewares.DatabaseMiddleware(db))
	router.Use(middlewares.ErrorMiddleware())

	routes.LoadRoutes(router)

	if err := router.Run(os.Getenv("SERVER_URI")); err != nil {
		log.Fatal(err.Error())
		return
	}
}
