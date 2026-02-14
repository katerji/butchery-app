package admin

import "errors"

var (
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrEmptyFullName      = errors.New("full name must not be empty")
	ErrAdminNotFound      = errors.New("admin not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
