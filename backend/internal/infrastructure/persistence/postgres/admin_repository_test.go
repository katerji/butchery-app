package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/katerji/butchery-app/backend/internal/domain/admin"
	pgstore "github.com/katerji/butchery-app/backend/internal/infrastructure/persistence/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAdminRepository_FindByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := setupTestDB(t)
	repo := pgstore.NewAdminRepository(pool)
	ctx := context.Background()

	// Seed an admin
	adminID := uuid.New()
	_, err := pool.Exec(ctx,
		"INSERT INTO admins (id, email, password_hash, full_name) VALUES ($1, $2, $3, $4)",
		adminID, "admin@butchery.com", "$2a$10$hashvalue", "Butchery Admin",
	)
	require.NoError(t, err)

	t.Run("existing email returns admin", func(t *testing.T) {
		a, err := repo.FindByEmail(ctx, "admin@butchery.com")

		require.NoError(t, err)
		assert.Equal(t, adminID, a.ID())
		assert.Equal(t, "admin@butchery.com", a.Email())
		assert.Equal(t, "$2a$10$hashvalue", a.PasswordHash())
		assert.Equal(t, "Butchery Admin", a.FullName())
	})

	t.Run("non-existing email returns ErrAdminNotFound", func(t *testing.T) {
		_, err := repo.FindByEmail(ctx, "unknown@butchery.com")

		assert.ErrorIs(t, err, admin.ErrAdminNotFound)
	})
}

func TestIntegrationAdminRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := setupTestDB(t)
	repo := pgstore.NewAdminRepository(pool)
	ctx := context.Background()

	adminID := uuid.New()
	_, err := pool.Exec(ctx,
		"INSERT INTO admins (id, email, password_hash, full_name) VALUES ($1, $2, $3, $4)",
		adminID, "admin@butchery.com", "$2a$10$hashvalue", "Butchery Admin",
	)
	require.NoError(t, err)

	t.Run("existing ID returns admin", func(t *testing.T) {
		a, err := repo.FindByID(ctx, adminID)

		require.NoError(t, err)
		assert.Equal(t, adminID, a.ID())
		assert.Equal(t, "admin@butchery.com", a.Email())
	})

	t.Run("non-existing ID returns ErrAdminNotFound", func(t *testing.T) {
		_, err := repo.FindByID(ctx, uuid.New())

		assert.ErrorIs(t, err, admin.ErrAdminNotFound)
	})
}
