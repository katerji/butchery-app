package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katerji/butchery-app/backend/internal/domain/auth"
)

// RefreshTokenRepository implements auth.RefreshTokenRepository using PostgreSQL.
type RefreshTokenRepository struct {
	pool *pgxpool.Pool
}

// NewRefreshTokenRepository creates a new RefreshTokenRepository.
func NewRefreshTokenRepository(pool *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{pool: pool}
}

// Save persists a refresh token.
func (r *RefreshTokenRepository) Save(ctx context.Context, token *auth.RefreshToken) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO refresh_tokens (id, subject_id, subject_type, token_hash, expires_at, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		token.ID(), token.SubjectID(), token.SubjectType(),
		token.TokenHash(), token.ExpiresAt(), token.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf("inserting refresh token: %w", err)
	}
	return nil
}

// FindByTokenHash finds a refresh token by its hash.
func (r *RefreshTokenRepository) FindByTokenHash(ctx context.Context, hash string) (*auth.RefreshToken, error) {
	var id, subjectID uuid.UUID
	var subjectType, tokenHash string
	var expiresAt, createdAt time.Time

	err := r.pool.QueryRow(ctx,
		"SELECT id, subject_id, subject_type, token_hash, expires_at, created_at FROM refresh_tokens WHERE token_hash = $1",
		hash,
	).Scan(&id, &subjectID, &subjectType, &tokenHash, &expiresAt, &createdAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, auth.ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("querying refresh token by hash: %w", err)
	}

	return auth.ReconstructRefreshToken(id, subjectID, subjectType, tokenHash, expiresAt, createdAt), nil
}

// DeleteBySubjectID deletes all refresh tokens for a given subject.
func (r *RefreshTokenRepository) DeleteBySubjectID(ctx context.Context, subjectID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		"DELETE FROM refresh_tokens WHERE subject_id = $1",
		subjectID,
	)
	if err != nil {
		return fmt.Errorf("deleting refresh tokens by subject: %w", err)
	}
	return nil
}

// DeleteByTokenHash deletes a refresh token by its hash.
func (r *RefreshTokenRepository) DeleteByTokenHash(ctx context.Context, hash string) error {
	result, err := r.pool.Exec(ctx,
		"DELETE FROM refresh_tokens WHERE token_hash = $1",
		hash,
	)
	if err != nil {
		return fmt.Errorf("deleting refresh token by hash: %w", err)
	}
	if result.RowsAffected() == 0 {
		return auth.ErrRefreshTokenNotFound
	}
	return nil
}
