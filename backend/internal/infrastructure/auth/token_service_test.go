package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	infraauth "github.com/katerji/butchery-app/backend/internal/infrastructure/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenService_GenerateAccessToken_ReturnsValidToken(t *testing.T) {
	svc := infraauth.NewTokenService("test-secret", 15*time.Minute)
	subjectID := uuid.New()

	token, err := svc.GenerateAccessToken(subjectID, "admin")

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestTokenService_ValidateAccessToken_ValidToken_ReturnsClaims(t *testing.T) {
	svc := infraauth.NewTokenService("test-secret", 15*time.Minute)
	subjectID := uuid.New()

	token, _ := svc.GenerateAccessToken(subjectID, "customer")
	claims, err := svc.ValidateAccessToken(token)

	require.NoError(t, err)
	assert.Equal(t, subjectID, claims.SubjectID)
	assert.Equal(t, "customer", claims.SubjectType)
}

func TestTokenService_ValidateAccessToken_ExpiredToken_ReturnsError(t *testing.T) {
	svc := infraauth.NewTokenService("test-secret", -1*time.Minute)
	subjectID := uuid.New()

	token, _ := svc.GenerateAccessToken(subjectID, "admin")
	_, err := svc.ValidateAccessToken(token)

	assert.Error(t, err)
}

func TestTokenService_ValidateAccessToken_WrongSecret_ReturnsError(t *testing.T) {
	svc1 := infraauth.NewTokenService("secret-1", 15*time.Minute)
	svc2 := infraauth.NewTokenService("secret-2", 15*time.Minute)
	subjectID := uuid.New()

	token, _ := svc1.GenerateAccessToken(subjectID, "admin")
	_, err := svc2.ValidateAccessToken(token)

	assert.Error(t, err)
}

func TestTokenService_GenerateRefreshToken_ReturnsUniqueTokens(t *testing.T) {
	svc := infraauth.NewTokenService("test-secret", 15*time.Minute)

	token1, err1 := svc.GenerateRefreshToken()
	token2, err2 := svc.GenerateRefreshToken()

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, token2)
	assert.NotEqual(t, token1, token2)
}

func TestTokenService_ValidateAccessToken_InvalidString_ReturnsError(t *testing.T) {
	svc := infraauth.NewTokenService("test-secret", 15*time.Minute)

	_, err := svc.ValidateAccessToken("not-a-jwt")

	assert.Error(t, err)
}
