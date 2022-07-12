package middlewares

import (
	. "self-hosted-cloud/server/database"

	"github.com/gin-gonic/gin"
)

func DatabaseMiddleware(database *Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(KeyDatabase, database)
		c.Next()
	}
}
