package auth

import (
	"time"

	"github.com/google/uuid"
)

const (
	SubjectTypeAdmin    = "admin"
	SubjectTypeCustomer = "customer"
)

// RefreshToken represents a stored refresh token used to issue new access tokens.
type RefreshToken struct {
	id          uuid.UUID
	subjectID   uuid.UUID
	subjectType string
	tokenHash   string
	expiresAt   time.Time
	createdAt   time.Time
}

// NewRefreshToken creates a new RefreshToken with validation.
func NewRefreshToken(subjectID uuid.UUID, subjectType string, tokenHash string, expiresAt time.Time) (*RefreshToken, error) {
	if subjectType != SubjectTypeAdmin && subjectType != SubjectTypeCustomer {
		return nil, ErrInvalidSubjectType
	}
	if tokenHash == "" {
		return nil, ErrEmptyTokenHash
	}

	return &RefreshToken{
		id:          uuid.New(),
		subjectID:   subjectID,
		subjectType: subjectType,
		tokenHash:   tokenHash,
		expiresAt:   expiresAt,
		createdAt:   time.Now(),
	}, nil
}

// ReconstructRefreshToken reconstructs a RefreshToken from persistence without validation.
func ReconstructRefreshToken(id, subjectID uuid.UUID, subjectType, tokenHash string, expiresAt, createdAt time.Time) *RefreshToken {
	return &RefreshToken{
		id:          id,
		subjectID:   subjectID,
		subjectType: subjectType,
		tokenHash:   tokenHash,
		expiresAt:   expiresAt,
		createdAt:   createdAt,
	}
}

func (rt *RefreshToken) ID() uuid.UUID       { return rt.id }
func (rt *RefreshToken) SubjectID() uuid.UUID { return rt.subjectID }
func (rt *RefreshToken) SubjectType() string  { return rt.subjectType }
func (rt *RefreshToken) TokenHash() string    { return rt.tokenHash }
func (rt *RefreshToken) ExpiresAt() time.Time { return rt.expiresAt }
func (rt *RefreshToken) CreatedAt() time.Time { return rt.createdAt }

// IsExpired returns true if the refresh token has passed its expiration time.
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.expiresAt)
}
