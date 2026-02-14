package commands_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/katerji/butchery-app/backend/internal/application/auth/commands"
	"github.com/katerji/butchery-app/backend/internal/domain/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Mocks ---

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

// --- RefreshToken Tests ---

func TestRefreshToken_ValidToken_ReturnsNewAccessToken(t *testing.T) {
	refreshRepo := new(mockRefreshTokenRepository)
	tokenGen := new(mockTokenGenerator)

	subjectID := uuid.New()
	storedToken := auth.ReconstructRefreshToken(
		uuid.New(), subjectID, "customer", "hashed-value",
		time.Now().Add(7*24*time.Hour), time.Now(),
	)

	refreshRepo.On("FindByTokenHash", mock.Anything, mock.AnythingOfType("string")).Return(storedToken, nil)
	tokenGen.On("GenerateAccessToken", subjectID, "customer").Return("new-access-token", nil)

	handler := commands.NewRefreshTokenHandler(refreshRepo, tokenGen, 15*time.Minute)
	result, err := handler.Handle(context.Background(), commands.RefreshTokenCommand{
		RefreshToken: "raw-refresh-token",
	})

	require.NoError(t, err)
	assert.Equal(t, "new-access-token", result.AccessToken)
	assert.Greater(t, result.ExpiresIn, int64(0))
}

func TestRefreshToken_ExpiredToken_ReturnsError(t *testing.T) {
	refreshRepo := new(mockRefreshTokenRepository)
	tokenGen := new(mockTokenGenerator)

	subjectID := uuid.New()
	storedToken := auth.ReconstructRefreshToken(
		uuid.New(), subjectID, "customer", "hashed-value",
		time.Now().Add(-1*time.Hour), time.Now().Add(-8*24*time.Hour),
	)

	refreshRepo.On("FindByTokenHash", mock.Anything, mock.AnythingOfType("string")).Return(storedToken, nil)

	handler := commands.NewRefreshTokenHandler(refreshRepo, tokenGen, 15*time.Minute)
	_, err := handler.Handle(context.Background(), commands.RefreshTokenCommand{
		RefreshToken: "raw-refresh-token",
	})

	assert.ErrorIs(t, err, auth.ErrRefreshTokenExpired)
}

func TestRefreshToken_UnknownToken_ReturnsError(t *testing.T) {
	refreshRepo := new(mockRefreshTokenRepository)
	tokenGen := new(mockTokenGenerator)

	refreshRepo.On("FindByTokenHash", mock.Anything, mock.AnythingOfType("string")).Return(nil, auth.ErrRefreshTokenNotFound)

	handler := commands.NewRefreshTokenHandler(refreshRepo, tokenGen, 15*time.Minute)
	_, err := handler.Handle(context.Background(), commands.RefreshTokenCommand{
		RefreshToken: "unknown-token",
	})

	assert.ErrorIs(t, err, auth.ErrRefreshTokenNotFound)
}
