package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/katerji/butchery-app/backend/internal/domain/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRefreshToken_ValidInputs_CreatesToken(t *testing.T) {
	subjectID := uuid.New()
	tokenHash := "hashed-token-value"
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	token, err := auth.NewRefreshToken(subjectID, "customer", tokenHash, expiresAt)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, token.ID())
	assert.Equal(t, subjectID, token.SubjectID())
	assert.Equal(t, "customer", token.SubjectType())
	assert.Equal(t, tokenHash, token.TokenHash())
	assert.Equal(t, expiresAt, token.ExpiresAt())
	assert.False(t, token.CreatedAt().IsZero())
}

func TestNewRefreshToken_AdminSubjectType_CreatesToken(t *testing.T) {
	subjectID := uuid.New()
	tokenHash := "hashed-token-value"
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	token, err := auth.NewRefreshToken(subjectID, "admin", tokenHash, expiresAt)

	require.NoError(t, err)
	assert.Equal(t, "admin", token.SubjectType())
}

func TestNewRefreshToken_InvalidSubjectType_ReturnsError(t *testing.T) {
	subjectID := uuid.New()
	tokenHash := "hashed-token-value"
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err := auth.NewRefreshToken(subjectID, "unknown", tokenHash, expiresAt)

	assert.ErrorIs(t, err, auth.ErrInvalidSubjectType)
}

func TestNewRefreshToken_EmptyTokenHash_ReturnsError(t *testing.T) {
	subjectID := uuid.New()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err := auth.NewRefreshToken(subjectID, "customer", "", expiresAt)

	assert.ErrorIs(t, err, auth.ErrEmptyTokenHash)
}

func TestRefreshToken_IsExpired_WhenExpired_ReturnsTrue(t *testing.T) {
	subjectID := uuid.New()
	tokenHash := "hashed-token-value"
	expiresAt := time.Now().Add(-1 * time.Hour)

	token, err := auth.NewRefreshToken(subjectID, "customer", tokenHash, expiresAt)

	require.NoError(t, err)
	assert.True(t, token.IsExpired())
}

func TestRefreshToken_IsExpired_WhenNotExpired_ReturnsFalse(t *testing.T) {
	subjectID := uuid.New()
	tokenHash := "hashed-token-value"
	expiresAt := time.Now().Add(1 * time.Hour)

	token, err := auth.NewRefreshToken(subjectID, "customer", tokenHash, expiresAt)

	require.NoError(t, err)
	assert.False(t, token.IsExpired())
}

func TestReconstructRefreshToken_RestoresAllFields(t *testing.T) {
	id := uuid.New()
	subjectID := uuid.New()
	tokenHash := "hashed-token-value"
	expiresAt := time.Now().Add(1 * time.Hour)
	createdAt := time.Now().Add(-1 * time.Hour)

	token := auth.ReconstructRefreshToken(id, subjectID, "admin", tokenHash, expiresAt, createdAt)

	assert.Equal(t, id, token.ID())
	assert.Equal(t, subjectID, token.SubjectID())
	assert.Equal(t, "admin", token.SubjectType())
	assert.Equal(t, tokenHash, token.TokenHash())
	assert.Equal(t, expiresAt, token.ExpiresAt())
	assert.Equal(t, createdAt, token.CreatedAt())
}
