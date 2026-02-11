package commands_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/katerji/butchery-app/backend/internal/application/customer/commands"
	"github.com/katerji/butchery-app/backend/internal/domain/customer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCustomerLogin_ValidCredentials_ReturnsTokens(t *testing.T) {
	custRepo := new(mockCustomerRepository)
	hasher := new(mockPasswordHasher)
	tokenGen := new(mockTokenGenerator)
	refreshRepo := new(mockRefreshTokenRepository)

	customerID := uuid.New()
	email, _ := customer.NewEmail("user@example.com")
	phone, _ := customer.NewPhoneNumber("+1234567890")
	c := customer.ReconstructCustomer(customerID, email, "$2a$10$hash", "John Doe", phone)

	custRepo.On("FindByEmail", mock.Anything, email).Return(c, nil)
	hasher.On("Compare", "$2a$10$hash", "password123").Return(nil)
	tokenGen.On("GenerateAccessToken", customerID, "customer").Return("access-token", nil)
	tokenGen.On("GenerateRefreshToken").Return("refresh-token-raw", nil)
	refreshRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	handler := commands.NewCustomerLoginHandler(custRepo, hasher, tokenGen, refreshRepo)
	result, err := handler.Handle(context.Background(), commands.CustomerLoginCommand{
		Email:    "user@example.com",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.Equal(t, "access-token", result.AccessToken)
	assert.Equal(t, "refresh-token-raw", result.RefreshToken)
	assert.Greater(t, result.ExpiresIn, int64(0))
}

func TestCustomerLogin_CustomerNotFound_ReturnsError(t *testing.T) {
	custRepo := new(mockCustomerRepository)
	hasher := new(mockPasswordHasher)
	tokenGen := new(mockTokenGenerator)
	refreshRepo := new(mockRefreshTokenRepository)

	email, _ := customer.NewEmail("unknown@example.com")
	custRepo.On("FindByEmail", mock.Anything, email).Return(nil, customer.ErrCustomerNotFound)

	handler := commands.NewCustomerLoginHandler(custRepo, hasher, tokenGen, refreshRepo)
	_, err := handler.Handle(context.Background(), commands.CustomerLoginCommand{
		Email:    "unknown@example.com",
		Password: "password123",
	})

	assert.ErrorIs(t, err, customer.ErrInvalidCredentials)
}

func TestCustomerLogin_WrongPassword_ReturnsError(t *testing.T) {
	custRepo := new(mockCustomerRepository)
	hasher := new(mockPasswordHasher)
	tokenGen := new(mockTokenGenerator)
	refreshRepo := new(mockRefreshTokenRepository)

	customerID := uuid.New()
	email, _ := customer.NewEmail("user@example.com")
	phone, _ := customer.NewPhoneNumber("+1234567890")
	c := customer.ReconstructCustomer(customerID, email, "$2a$10$hash", "John Doe", phone)

	custRepo.On("FindByEmail", mock.Anything, email).Return(c, nil)
	hasher.On("Compare", "$2a$10$hash", "wrongpassword").Return(errors.New("mismatch"))

	handler := commands.NewCustomerLoginHandler(custRepo, hasher, tokenGen, refreshRepo)
	_, err := handler.Handle(context.Background(), commands.CustomerLoginCommand{
		Email:    "user@example.com",
		Password: "wrongpassword",
	})

	assert.ErrorIs(t, err, customer.ErrInvalidCredentials)
}
