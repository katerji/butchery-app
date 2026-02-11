package customer

import "errors"

var (
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters")
	ErrInvalidPhoneNumber = errors.New("phone number must not be empty")
	ErrEmptyFullName      = errors.New("full name must not be empty")
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
