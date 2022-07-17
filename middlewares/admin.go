package middlewares

import (
	"errors"
	"net/http"
	"self-hosted-cloud/server/utils"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := utils.GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if user.Role != "admin" {
			err := errors.New("you must be admin to access this route")
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		c.Next()
	}
}
