package main

import (
	"log"
	"os"
	"self-hosted-cloud/server/database"
	"self-hosted-cloud/server/middlewares"
	"self-hosted-cloud/server/routes/admin"
	"self-hosted-cloud/server/routes/auth"
	"self-hosted-cloud/server/routes/storage"
	"self-hosted-cloud/server/routes/user"
	adminservice "self-hosted-cloud/server/services/admin"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	if godotenv.Load() != nil {
		log.Fatal(".env Couldn't be loaded.")
	}

	err := os.Mkdir(os.Getenv("DATA_PATH"), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err.Error())
		return
	}

	db, err := database.GetDatabase("database.sqlite")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	appIsInDemoMode, serviceError := adminservice.AppIsInDemoMode()
	if serviceError != nil {
		log.Fatal(serviceError.Error())
	}
	if appIsInDemoMode {
		adminservice.StartDemoMode(&db)
	}

	router := gin.Default()
	router.Use(cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		MaxAge:         int(12 * time.Hour),
	}))
	router.Use(middlewares.DatabaseMiddleware(&db))
	router.Use(middlewares.ErrorMiddleware())

	auth.LoadRoutes(router)
	user.LoadRoutes(router)
	storage.LoadRoutes(router)
	admin.LoadRoutes(router)

	err = router.Run(os.Getenv("SERVER_URI"))
	if err != nil {
		return
	}
}
