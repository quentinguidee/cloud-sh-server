package commands

type Command interface {
	Run() ICommandError
	Revert() ICommandError
}

type ICommandError interface {
	Error() error
	Code() int
}

type CommandError struct {
	code  int
	error error
}

func NewError(code int, error error) CommandError {
	return CommandError{
		code,
		error,
	}
}

func (c CommandError) Error() error {
	return c.error
}

func (c CommandError) Code() int {
	return c.code
}
