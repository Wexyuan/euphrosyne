package errs

import (
	"errors"
	"fmt"
)

// Error is a business error carrying a code and a message.
type Error struct {
	// Code identifies the failure.
	Code int
	// Message describes the failure.
	Message string

	// cause keeps the original error for the error chain.
	cause error
}

// New creates an error with the given code and message.
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Newf creates an error with a formatted message.
func Newf(code int, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap attaches a code and message to an existing error.
func Wrap(err error, code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		cause:   err,
	}
}

// Error returns the message of the error.
func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.cause)
	}
	return e.Message
}

// Unwrap returns the wrapped cause.
func (e *Error) Unwrap() error {
	return e.cause
}

// Is reports whether target carries the same code.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && t.Code == e.Code
}

// From extracts the business error from the error chain.
func From(err error) (*Error, bool) {
	var target *Error
	if errors.As(err, &target) {
		return target, true
	}
	return nil, false
}
