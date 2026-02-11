package customer_test

import (
	"testing"

	"github.com/katerji/butchery-app/backend/internal/domain/customer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Email Value Object ---

func TestNewEmail_ValidEmail_CreatesEmail(t *testing.T) {
	email, err := customer.NewEmail("user@example.com")

	require.NoError(t, err)
	assert.Equal(t, "user@example.com", email.String())
}

func TestNewEmail_NormalizesToLowercase(t *testing.T) {
	email, err := customer.NewEmail("User@EXAMPLE.COM")

	require.NoError(t, err)
	assert.Equal(t, "user@example.com", email.String())
}

func TestNewEmail_InvalidFormats_ReturnError(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{"empty", ""},
		{"no at sign", "userexample.com"},
		{"no domain", "user@"},
		{"no local part", "@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := customer.NewEmail(tt.email)
			assert.ErrorIs(t, err, customer.ErrInvalidEmail)
		})
	}
}

func TestEmail_Equals_SameValue_ReturnsTrue(t *testing.T) {
	e1, _ := customer.NewEmail("user@example.com")
	e2, _ := customer.NewEmail("USER@example.com")

	assert.True(t, e1.Equals(e2))
}

func TestEmail_Equals_DifferentValue_ReturnsFalse(t *testing.T) {
	e1, _ := customer.NewEmail("user1@example.com")
	e2, _ := customer.NewEmail("user2@example.com")

	assert.False(t, e1.Equals(e2))
}

// --- Password Value Object ---

func TestNewPassword_ValidPassword_CreatesPassword(t *testing.T) {
	pw, err := customer.NewPassword("securepass123")

	require.NoError(t, err)
	assert.Equal(t, "securepass123", pw.String())
}

func TestNewPassword_TooShort_ReturnsError(t *testing.T) {
	_, err := customer.NewPassword("short")

	assert.ErrorIs(t, err, customer.ErrInvalidPassword)
}

func TestNewPassword_ExactlyMinLength_Succeeds(t *testing.T) {
	_, err := customer.NewPassword("12345678")

	assert.NoError(t, err)
}

func TestNewPassword_Empty_ReturnsError(t *testing.T) {
	_, err := customer.NewPassword("")

	assert.ErrorIs(t, err, customer.ErrInvalidPassword)
}

// --- PhoneNumber Value Object ---

func TestNewPhoneNumber_ValidNumber_CreatesPhoneNumber(t *testing.T) {
	phone, err := customer.NewPhoneNumber("+1234567890")

	require.NoError(t, err)
	assert.Equal(t, "+1234567890", phone.String())
}

func TestNewPhoneNumber_Empty_ReturnsError(t *testing.T) {
	_, err := customer.NewPhoneNumber("")

	assert.ErrorIs(t, err, customer.ErrInvalidPhoneNumber)
}

func TestNewPhoneNumber_WhitespaceOnly_ReturnsError(t *testing.T) {
	_, err := customer.NewPhoneNumber("   ")

	assert.ErrorIs(t, err, customer.ErrInvalidPhoneNumber)
}
