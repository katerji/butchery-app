package auth_test

import (
	"testing"

	infraauth "github.com/katerji/butchery-app/backend/internal/infrastructure/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBcryptHasher_Hash_ReturnsHashedPassword(t *testing.T) {
	hasher := infraauth.NewBcryptHasher()

	hash, err := hasher.Hash("password123")

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "password123", hash)
}

func TestBcryptHasher_Compare_MatchingPassword_ReturnsNil(t *testing.T) {
	hasher := infraauth.NewBcryptHasher()

	hash, _ := hasher.Hash("password123")
	err := hasher.Compare(hash, "password123")

	assert.NoError(t, err)
}

func TestBcryptHasher_Compare_WrongPassword_ReturnsError(t *testing.T) {
	hasher := infraauth.NewBcryptHasher()

	hash, _ := hasher.Hash("password123")
	err := hasher.Compare(hash, "wrongpassword")

	assert.Error(t, err)
}
