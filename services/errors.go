package services

import "github.com/gin-gonic/gin"

type IServiceError interface {
	Code() int
	Error() error
	Throws(c *gin.Context)
}

type ServiceError struct {
	code  int
	error error
}

func NewServiceError(code int, error error) ServiceError {
	return ServiceError{
		code,
		error,
	}
}

func (s ServiceError) Error() error {
	return s.error
}

func (s ServiceError) Code() int {
	return s.code
}

func (s ServiceError) Throws(c *gin.Context) {
	c.AbortWithError(s.code, s.error)
}
