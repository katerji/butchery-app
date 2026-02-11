package customer_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/katerji/butchery-app/backend/internal/domain/customer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCustomer_ValidInputs_CreatesCustomer(t *testing.T) {
	id := uuid.New()
	email, _ := customer.NewEmail("user@example.com")
	phone, _ := customer.NewPhoneNumber("+1234567890")

	c, err := customer.NewCustomer(id, email, "$2a$10$hashedpassword", "John Doe", phone)

	require.NoError(t, err)
	assert.Equal(t, id, c.ID())
	assert.Equal(t, "user@example.com", c.Email().String())
	assert.Equal(t, "$2a$10$hashedpassword", c.PasswordHash())
	assert.Equal(t, "John Doe", c.FullName())
	assert.Equal(t, "+1234567890", c.Phone().String())
	assert.False(t, c.CreatedAt().IsZero())
	assert.False(t, c.UpdatedAt().IsZero())
}

func TestNewCustomer_EmptyFullName_ReturnsError(t *testing.T) {
	id := uuid.New()
	email, _ := customer.NewEmail("user@example.com")
	phone, _ := customer.NewPhoneNumber("+1234567890")

	_, err := customer.NewCustomer(id, email, "$2a$10$hash", "", phone)

	assert.ErrorIs(t, err, customer.ErrEmptyFullName)
}

func TestNewCustomer_WhitespaceFullName_ReturnsError(t *testing.T) {
	id := uuid.New()
	email, _ := customer.NewEmail("user@example.com")
	phone, _ := customer.NewPhoneNumber("+1234567890")

	_, err := customer.NewCustomer(id, email, "$2a$10$hash", "   ", phone)

	assert.ErrorIs(t, err, customer.ErrEmptyFullName)
}

func TestReconstructCustomer_RestoresAllFields(t *testing.T) {
	id := uuid.New()
	email, _ := customer.NewEmail("user@example.com")
	phone, _ := customer.NewPhoneNumber("+1234567890")

	c := customer.ReconstructCustomer(id, email, "$2a$10$hash", "John Doe", phone)

	assert.Equal(t, id, c.ID())
	assert.Equal(t, "user@example.com", c.Email().String())
	assert.Equal(t, "$2a$10$hash", c.PasswordHash())
	assert.Equal(t, "John Doe", c.FullName())
	assert.Equal(t, "+1234567890", c.Phone().String())
}
