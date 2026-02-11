package commands_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/katerji/butchery-app/backend/internal/application/customer/commands"
	"github.com/katerji/butchery-app/backend/internal/domain/customer"
	"github.com/katerji/butchery-app/backend/internal/domain/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Mocks ---

type mockCustomerRepository struct {
	mock.Mock
}

func (m *mockCustomerRepository) Save(ctx context.Context, c *customer.Customer) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *mockCustomerRepository) FindByEmail(ctx context.Context, email customer.Email) (*customer.Customer, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customer.Customer), args.Error(1)
}

func (m *mockCustomerRepository) FindByID(ctx context.Context, id uuid.UUID) (*customer.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customer.Customer), args.Error(1)
}

func (m *mockCustomerRepository) ExistsByEmail(ctx context.Context, email customer.Email) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

type mockPasswordHasher struct {
	mock.Mock
}

func (m *mockPasswordHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *mockPasswordHasher) Compare(hashed string, plain string) error {
	args := m.Called(hashed, plain)
	return args.Error(0)
}

type mockTokenGenerator struct {
	mock.Mock
}

func (m *mockTokenGenerator) GenerateAccessToken(subjectID uuid.UUID, subjectType string) (string, error) {
	args := m.Called(subjectID, subjectType)
	return args.String(0), args.Error(1)
}

func (m *mockTokenGenerator) GenerateRefreshToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type mockRefreshTokenRepository struct {
	mock.Mock
}

func (m *mockRefreshTokenRepository) Save(ctx context.Context, token *auth.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *mockRefreshTokenRepository) FindByTokenHash(ctx context.Context, hash string) (*auth.RefreshToken, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.RefreshToken), args.Error(1)
}

func (m *mockRefreshTokenRepository) DeleteBySubjectID(ctx context.Context, subjectID uuid.UUID) error {
	args := m.Called(ctx, subjectID)
	return args.Error(0)
}

func (m *mockRefreshTokenRepository) DeleteByTokenHash(ctx context.Context, hash string) error {
	args := m.Called(ctx, hash)
	return args.Error(0)
}

// --- RegisterCustomer Tests ---

func TestRegisterCustomer_ValidInputs_CreatesCustomer(t *testing.T) {
	custRepo := new(mockCustomerRepository)
	hasher := new(mockPasswordHasher)

	email, _ := customer.NewEmail("user@example.com")
	custRepo.On("ExistsByEmail", mock.Anything, email).Return(false, nil)
	hasher.On("Hash", "password123").Return("$2a$10$hashed", nil)
	custRepo.On("Save", mock.Anything, mock.AnythingOfType("*customer.Customer")).Return(nil)

	handler := commands.NewRegisterCustomerHandler(custRepo, hasher)
	result, err := handler.Handle(context.Background(), commands.RegisterCustomerCommand{
		Email:    "user@example.com",
		Password: "password123",
		FullName: "John Doe",
		Phone:    "+1234567890",
	})

	require.NoError(t, err)
	assert.Equal(t, "user@example.com", result.Email)
	assert.Equal(t, "John Doe", result.FullName)
	assert.NotEqual(t, uuid.Nil, result.CustomerID)
	custRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
}

func TestRegisterCustomer_EmailAlreadyExists_ReturnsError(t *testing.T) {
	custRepo := new(mockCustomerRepository)
	hasher := new(mockPasswordHasher)

	email, _ := customer.NewEmail("existing@example.com")
	custRepo.On("ExistsByEmail", mock.Anything, email).Return(true, nil)

	handler := commands.NewRegisterCustomerHandler(custRepo, hasher)
	_, err := handler.Handle(context.Background(), commands.RegisterCustomerCommand{
		Email:    "existing@example.com",
		Password: "password123",
		FullName: "John Doe",
		Phone:    "+1234567890",
	})

	assert.ErrorIs(t, err, customer.ErrEmailAlreadyExists)
}

func TestRegisterCustomer_InvalidEmail_ReturnsError(t *testing.T) {
	custRepo := new(mockCustomerRepository)
	hasher := new(mockPasswordHasher)

	handler := commands.NewRegisterCustomerHandler(custRepo, hasher)
	_, err := handler.Handle(context.Background(), commands.RegisterCustomerCommand{
		Email:    "invalid-email",
		Password: "password123",
		FullName: "John Doe",
		Phone:    "+1234567890",
	})

	assert.ErrorIs(t, err, customer.ErrInvalidEmail)
}

func TestRegisterCustomer_PasswordTooShort_ReturnsError(t *testing.T) {
	custRepo := new(mockCustomerRepository)
	hasher := new(mockPasswordHasher)

	handler := commands.NewRegisterCustomerHandler(custRepo, hasher)
	_, err := handler.Handle(context.Background(), commands.RegisterCustomerCommand{
		Email:    "user@example.com",
		Password: "short",
		FullName: "John Doe",
		Phone:    "+1234567890",
	})

	assert.ErrorIs(t, err, customer.ErrInvalidPassword)
}

func TestRegisterCustomer_InvalidPhoneNumber_ReturnsError(t *testing.T) {
	custRepo := new(mockCustomerRepository)
	hasher := new(mockPasswordHasher)

	handler := commands.NewRegisterCustomerHandler(custRepo, hasher)
	_, err := handler.Handle(context.Background(), commands.RegisterCustomerCommand{
		Email:    "user@example.com",
		Password: "password123",
		FullName: "John Doe",
		Phone:    "",
	})

	assert.ErrorIs(t, err, customer.ErrInvalidPhoneNumber)
}
