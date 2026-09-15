package ecode

import "github.com/Wexyuan/euphrosyne/pkg/errs"

// Common error codes shared by all modules, ranged 1000-1999.
var (
	ErrInternalServer = errs.New(1000, "server internal error")
	ErrBadRequest     = errs.New(1001, "bad request error")
	ErrValidation     = errs.New(1002, "validation error")
	ErrUnauthorized   = errs.New(1003, "unauthorized error")
	ErrForbidden      = errs.New(1004, "forbidden error")
	ErrNotFound       = errs.New(1005, "not found error")
	ErrTimeout        = errs.New(1006, "timeout error")
)

// User error codes, ranged 2100-2199.
var (
	ErrUserNotFound = errs.New(2100, "user not found error")
	ErrUserExists   = errs.New(2101, "user already exists error")
	ErrUserDisabled = errs.New(2102, "user disabled error")
)
