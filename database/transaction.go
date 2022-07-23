package database

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewTX(c *gin.Context) *gorm.DB {
	db := GetDatabaseFromContext(c)
	return db.Begin()
}
