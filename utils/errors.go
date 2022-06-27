package utils

import "github.com/gin-gonic/gin"

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		c.JSON(-1, gin.H{
			"message": err.Error(),
		})
	}
}
