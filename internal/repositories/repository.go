package repositories

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
)

var (
	// PostgreSQL error codes
	UniqueViolationErrorCode = "23505"
)
