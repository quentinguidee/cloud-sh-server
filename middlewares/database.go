package middlewares

import (
	. "self-hosted-cloud/server/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DatabaseMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(KeyDatabase, db)
		c.Next()
	}
}
