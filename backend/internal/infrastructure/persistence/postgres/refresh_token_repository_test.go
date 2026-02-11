package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/katerji/butchery-app/backend/internal/domain/auth"
	pgstore "github.com/katerji/butchery-app/backend/internal/infrastructure/persistence/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRefreshToken(t *testing.T, subjectID uuid.UUID, subjectType string) *auth.RefreshToken {
	t.Helper()
	token, err := auth.NewRefreshToken(subjectID, subjectType, "token-hash-"+uuid.NewString(), time.Now().Add(7*24*time.Hour))
	require.NoError(t, err)
	return token
}

func TestIntegrationRefreshTokenRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := setupTestDB(t)
	repo := pgstore.NewRefreshTokenRepository(pool)
	ctx := context.Background()

	t.Run("saves and retrieves refresh token", func(t *testing.T) {
		truncateAll(t, pool)
		subjectID := uuid.New()
		token := newTestRefreshToken(t, subjectID, "customer")

		err := repo.Save(ctx, token)
		require.NoError(t, err)

		found, err := repo.FindByTokenHash(ctx, token.TokenHash())
		require.NoError(t, err)
		assert.Equal(t, token.ID(), found.ID())
		assert.Equal(t, subjectID, found.SubjectID())
		assert.Equal(t, "customer", found.SubjectType())
		assert.Equal(t, token.TokenHash(), found.TokenHash())
	})
}

func TestIntegrationRefreshTokenRepository_FindByTokenHash(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := setupTestDB(t)
	repo := pgstore.NewRefreshTokenRepository(pool)
	ctx := context.Background()
	truncateAll(t, pool)

	subjectID := uuid.New()
	token := newTestRefreshToken(t, subjectID, "admin")
	require.NoError(t, repo.Save(ctx, token))

	t.Run("existing hash returns token", func(t *testing.T) {
		found, err := repo.FindByTokenHash(ctx, token.TokenHash())

		require.NoError(t, err)
		assert.Equal(t, token.ID(), found.ID())
		assert.Equal(t, "admin", found.SubjectType())
	})

	t.Run("non-existing hash returns ErrRefreshTokenNotFound", func(t *testing.T) {
		_, err := repo.FindByTokenHash(ctx, "non-existing-hash")

		assert.ErrorIs(t, err, auth.ErrRefreshTokenNotFound)
	})
}

func TestIntegrationRefreshTokenRepository_DeleteBySubjectID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := setupTestDB(t)
	repo := pgstore.NewRefreshTokenRepository(pool)
	ctx := context.Background()
	truncateAll(t, pool)

	subjectID := uuid.New()
	token1 := newTestRefreshToken(t, subjectID, "customer")
	token2 := newTestRefreshToken(t, subjectID, "customer")
	require.NoError(t, repo.Save(ctx, token1))
	require.NoError(t, repo.Save(ctx, token2))

	err := repo.DeleteBySubjectID(ctx, subjectID)
	require.NoError(t, err)

	_, err = repo.FindByTokenHash(ctx, token1.TokenHash())
	assert.ErrorIs(t, err, auth.ErrRefreshTokenNotFound)

	_, err = repo.FindByTokenHash(ctx, token2.TokenHash())
	assert.ErrorIs(t, err, auth.ErrRefreshTokenNotFound)
}

func TestIntegrationRefreshTokenRepository_DeleteByTokenHash(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := setupTestDB(t)
	repo := pgstore.NewRefreshTokenRepository(pool)
	ctx := context.Background()
	truncateAll(t, pool)

	subjectID := uuid.New()
	token := newTestRefreshToken(t, subjectID, "admin")
	require.NoError(t, repo.Save(ctx, token))

	t.Run("existing token is deleted", func(t *testing.T) {
		err := repo.DeleteByTokenHash(ctx, token.TokenHash())
		require.NoError(t, err)

		_, err = repo.FindByTokenHash(ctx, token.TokenHash())
		assert.ErrorIs(t, err, auth.ErrRefreshTokenNotFound)
	})

	t.Run("non-existing token returns ErrRefreshTokenNotFound", func(t *testing.T) {
		err := repo.DeleteByTokenHash(ctx, "non-existing-hash")

		assert.ErrorIs(t, err, auth.ErrRefreshTokenNotFound)
	})
}
