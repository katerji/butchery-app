package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher implements auth.PasswordHasher using bcrypt.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new BcryptHasher with bcrypt.DefaultCost.
func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{cost: bcrypt.DefaultCost}
}

// Hash hashes a plaintext password using bcrypt.
func (h *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}
	return string(hash), nil
}

// Compare compares a bcrypt hash with a plaintext password.
func (h *BcryptHasher) Compare(hashed string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
