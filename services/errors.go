package services

type IServiceError interface {
	Code() int
	Error() error
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

func (c ServiceError) Error() error {
	return c.error
}

func (c ServiceError) Code() int {
	return c.code
}
