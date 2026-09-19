package errs

import (
	"errors"
	"fmt"
)

// Error is a business error.
type Error struct {
	Code    int    // business error code
	Message string // business error message
}

func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// From extracts the business error from the error chain.
func From(err error) (*Error, bool) {
	return errors.AsType[*Error](err)
}
