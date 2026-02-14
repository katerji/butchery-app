package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/katerji/butchery-app/backend/internal/domain/customer"
	pgstore "github.com/katerji/butchery-app/backend/internal/infrastructure/persistence/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCustomer(t *testing.T) *customer.Customer {
	t.Helper()
	email, err := customer.NewEmail("john@example.com")
	require.NoError(t, err)
	phone, err := customer.NewPhoneNumber("+1234567890")
	require.NoError(t, err)
	c, err := customer.NewCustomer(uuid.New(), email, "$2a$10$hash", "John Doe", phone)
	require.NoError(t, err)
	return c
}

func TestIntegrationCustomerRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := setupTestDB(t)
	repo := pgstore.NewCustomerRepository(pool)
	ctx := context.Background()

	t.Run("saves and retrieves customer", func(t *testing.T) {
		truncateAll(t, pool)
		c := newTestCustomer(t)

		err := repo.Save(ctx, c)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, c.ID())
		require.NoError(t, err)
		assert.Equal(t, c.ID(), found.ID())
		assert.Equal(t, c.Email().String(), found.Email().String())
		assert.Equal(t, c.FullName(), found.FullName())
		assert.Equal(t, c.Phone().String(), found.Phone().String())
	})

	t.Run("duplicate email returns error", func(t *testing.T) {
		truncateAll(t, pool)
		c1 := newTestCustomer(t)
		err := repo.Save(ctx, c1)
		require.NoError(t, err)

		// Create second customer with same email but different ID
		email, _ := customer.NewEmail("john@example.com")
		phone, _ := customer.NewPhoneNumber("+9876543210")
		c2, _ := customer.NewCustomer(uuid.New(), email, "$2a$10$hash2", "Jane Doe", phone)

		err = repo.Save(ctx, c2)
		assert.Error(t, err)
	})
}

func TestIntegrationCustomerRepository_FindByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := setupTestDB(t)
	repo := pgstore.NewCustomerRepository(pool)
	ctx := context.Background()
	truncateAll(t, pool)

	c := newTestCustomer(t)
	require.NoError(t, repo.Save(ctx, c))

	t.Run("existing email returns customer", func(t *testing.T) {
		email, _ := customer.NewEmail("john@example.com")
		found, err := repo.FindByEmail(ctx, email)

		require.NoError(t, err)
		assert.Equal(t, c.ID(), found.ID())
		assert.Equal(t, "john@example.com", found.Email().String())
	})

	t.Run("non-existing email returns ErrCustomerNotFound", func(t *testing.T) {
		email, _ := customer.NewEmail("unknown@example.com")
		_, err := repo.FindByEmail(ctx, email)

		assert.ErrorIs(t, err, customer.ErrCustomerNotFound)
	})
}

func TestIntegrationCustomerRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := setupTestDB(t)
	repo := pgstore.NewCustomerRepository(pool)
	ctx := context.Background()
	truncateAll(t, pool)

	c := newTestCustomer(t)
	require.NoError(t, repo.Save(ctx, c))

	t.Run("existing ID returns customer", func(t *testing.T) {
		found, err := repo.FindByID(ctx, c.ID())

		require.NoError(t, err)
		assert.Equal(t, c.ID(), found.ID())
	})

	t.Run("non-existing ID returns ErrCustomerNotFound", func(t *testing.T) {
		_, err := repo.FindByID(ctx, uuid.New())

		assert.ErrorIs(t, err, customer.ErrCustomerNotFound)
	})
}

func TestIntegrationCustomerRepository_ExistsByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool := setupTestDB(t)
	repo := pgstore.NewCustomerRepository(pool)
	ctx := context.Background()
	truncateAll(t, pool)

	c := newTestCustomer(t)
	require.NoError(t, repo.Save(ctx, c))

	t.Run("existing email returns true", func(t *testing.T) {
		email, _ := customer.NewEmail("john@example.com")
		exists, err := repo.ExistsByEmail(ctx, email)

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("non-existing email returns false", func(t *testing.T) {
		email, _ := customer.NewEmail("unknown@example.com")
		exists, err := repo.ExistsByEmail(ctx, email)

		require.NoError(t, err)
		assert.False(t, exists)
	})
}
