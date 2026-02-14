package auth

import "github.com/google/uuid"

// PasswordHasher hashes and compares passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed string, plain string) error
}

// TokenGenerator generates JWT access tokens and opaque refresh tokens.
type TokenGenerator interface {
	GenerateAccessToken(subjectID uuid.UUID, subjectType string) (string, error)
	GenerateRefreshToken() (string, error)
}

// TokenValidator validates JWT access tokens and extracts claims.
type TokenValidator interface {
	ValidateAccessToken(token string) (*AccessTokenClaims, error)
}

// AccessTokenClaims holds the claims extracted from a validated access token.
type AccessTokenClaims struct {
	SubjectID   uuid.UUID
	SubjectType string
}
