package commands_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/katerji/butchery-app/backend/internal/application/admin/commands"
	"github.com/katerji/butchery-app/backend/internal/domain/admin"
	"github.com/katerji/butchery-app/backend/internal/domain/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Mocks ---

type mockAdminRepository struct {
	mock.Mock
}

func (m *mockAdminRepository) FindByEmail(ctx context.Context, email string) (*admin.Admin, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*admin.Admin), args.Error(1)
}

func (m *mockAdminRepository) FindByID(ctx context.Context, id uuid.UUID) (*admin.Admin, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*admin.Admin), args.Error(1)
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

// --- Tests ---

func TestAdminLogin_ValidCredentials_ReturnsTokens(t *testing.T) {
	adminRepo := new(mockAdminRepository)
	hasher := new(mockPasswordHasher)
	tokenGen := new(mockTokenGenerator)
	refreshRepo := new(mockRefreshTokenRepository)

	adminID := uuid.New()
	a, _ := admin.NewAdmin(adminID, "admin@butchery.com", "$2a$10$hash", "Admin")

	adminRepo.On("FindByEmail", mock.Anything, "admin@butchery.com").Return(a, nil)
	hasher.On("Compare", "$2a$10$hash", "password123").Return(nil)
	tokenGen.On("GenerateAccessToken", adminID, "admin").Return("access-token", nil)
	tokenGen.On("GenerateRefreshToken").Return("refresh-token-raw", nil)
	refreshRepo.On("Save", mock.Anything, mock.AnythingOfType("*auth.RefreshToken")).Return(nil)

	handler := commands.NewAdminLoginHandler(adminRepo, hasher, tokenGen, refreshRepo, 15*time.Minute)
	result, err := handler.Handle(context.Background(), commands.AdminLoginCommand{
		Email:    "admin@butchery.com",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.Equal(t, "access-token", result.AccessToken)
	assert.Equal(t, "refresh-token-raw", result.RefreshToken)
	assert.Greater(t, result.ExpiresIn, int64(0))
	adminRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
	tokenGen.AssertExpectations(t)
	refreshRepo.AssertExpectations(t)
}

func TestAdminLogin_AdminNotFound_ReturnsError(t *testing.T) {
	adminRepo := new(mockAdminRepository)
	hasher := new(mockPasswordHasher)
	tokenGen := new(mockTokenGenerator)
	refreshRepo := new(mockRefreshTokenRepository)

	adminRepo.On("FindByEmail", mock.Anything, "unknown@butchery.com").Return(nil, admin.ErrAdminNotFound)

	handler := commands.NewAdminLoginHandler(adminRepo, hasher, tokenGen, refreshRepo, 15*time.Minute)
	_, err := handler.Handle(context.Background(), commands.AdminLoginCommand{
		Email:    "unknown@butchery.com",
		Password: "password123",
	})

	assert.ErrorIs(t, err, admin.ErrInvalidCredentials)
}

func TestAdminLogin_WrongPassword_ReturnsError(t *testing.T) {
	adminRepo := new(mockAdminRepository)
	hasher := new(mockPasswordHasher)
	tokenGen := new(mockTokenGenerator)
	refreshRepo := new(mockRefreshTokenRepository)

	adminID := uuid.New()
	a, _ := admin.NewAdmin(adminID, "admin@butchery.com", "$2a$10$hash", "Admin")

	adminRepo.On("FindByEmail", mock.Anything, "admin@butchery.com").Return(a, nil)
	hasher.On("Compare", "$2a$10$hash", "wrongpassword").Return(errors.New("mismatch"))

	handler := commands.NewAdminLoginHandler(adminRepo, hasher, tokenGen, refreshRepo, 15*time.Minute)
	_, err := handler.Handle(context.Background(), commands.AdminLoginCommand{
		Email:    "admin@butchery.com",
		Password: "wrongpassword",
	})

	assert.ErrorIs(t, err, admin.ErrInvalidCredentials)
}
