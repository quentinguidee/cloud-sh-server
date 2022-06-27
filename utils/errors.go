package utils

import "github.com/gin-gonic/gin"

func ErrorMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
		err := context.Errors.Last()
		context.JSON(-1, gin.H{
			"message": err.Error(),
		})
	}
}
