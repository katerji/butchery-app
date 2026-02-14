package auth

import (
	"context"

	"github.com/google/uuid"
)

// RefreshTokenRepository provides access to refresh token persistence.
type RefreshTokenRepository interface {
	Save(ctx context.Context, token *RefreshToken) error
	FindByTokenHash(ctx context.Context, hash string) (*RefreshToken, error)
	DeleteBySubjectID(ctx context.Context, subjectID uuid.UUID) error
	DeleteByTokenHash(ctx context.Context, hash string) error
}
