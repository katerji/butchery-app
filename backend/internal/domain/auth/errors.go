package auth

import "errors"

var (
	ErrInvalidSubjectType = errors.New("subject type must be 'admin' or 'customer'")
	ErrEmptyTokenHash     = errors.New("token hash must not be empty")
	ErrRefreshTokenExpired = errors.New("refresh token has expired")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)
