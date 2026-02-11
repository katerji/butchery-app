package admin_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/katerji/butchery-app/backend/internal/domain/admin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAdmin_ValidInputs_CreatesAdmin(t *testing.T) {
	id := uuid.New()

	a, err := admin.NewAdmin(id, "admin@butchery.com", "$2a$10$hashedpassword", "Butchery Admin")

	require.NoError(t, err)
	assert.Equal(t, id, a.ID())
	assert.Equal(t, "admin@butchery.com", a.Email())
	assert.Equal(t, "$2a$10$hashedpassword", a.PasswordHash())
	assert.Equal(t, "Butchery Admin", a.FullName())
}

func TestNewAdmin_InvalidEmail_ReturnsError(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{"empty email", ""},
		{"no at sign", "adminbutchery.com"},
		{"no domain", "admin@"},
		{"no local part", "@butchery.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := admin.NewAdmin(uuid.New(), tt.email, "$2a$10$hash", "Admin")
			assert.ErrorIs(t, err, admin.ErrInvalidEmail)
		})
	}
}

func TestNewAdmin_EmptyFullName_ReturnsError(t *testing.T) {
	_, err := admin.NewAdmin(uuid.New(), "admin@butchery.com", "$2a$10$hash", "")
	assert.ErrorIs(t, err, admin.ErrEmptyFullName)
}

func TestNewAdmin_WhitespaceOnlyFullName_ReturnsError(t *testing.T) {
	_, err := admin.NewAdmin(uuid.New(), "admin@butchery.com", "$2a$10$hash", "   ")
	assert.ErrorIs(t, err, admin.ErrEmptyFullName)
}
