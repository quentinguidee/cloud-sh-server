package main

import (
	"self-hosted-cloud/server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", routes.GetAbout)

	router.Run(":8080")
}
