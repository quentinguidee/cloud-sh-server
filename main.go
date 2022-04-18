package main

import (
	"self-hosted-cloud/server/auth"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	auth.Routes(router)

	router.Run(":8080")
}
