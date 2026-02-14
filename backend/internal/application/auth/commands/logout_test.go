package commands_test

import (
	"context"
	"testing"

	"github.com/katerji/butchery-app/backend/internal/application/auth/commands"
	"github.com/katerji/butchery-app/backend/internal/domain/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogout_ValidToken_DeletesRefreshToken(t *testing.T) {
	refreshRepo := new(mockRefreshTokenRepository)

	refreshRepo.On("DeleteByTokenHash", mock.Anything, mock.AnythingOfType("string")).Return(nil)

	handler := commands.NewLogoutHandler(refreshRepo)
	err := handler.Handle(context.Background(), commands.LogoutCommand{
		RefreshToken: "raw-refresh-token",
	})

	assert.NoError(t, err)
	refreshRepo.AssertExpectations(t)
}

func TestLogout_NonExistentToken_NoError(t *testing.T) {
	refreshRepo := new(mockRefreshTokenRepository)

	// Even if the token doesn't exist, DeleteByTokenHash should not error
	refreshRepo.On("DeleteByTokenHash", mock.Anything, mock.AnythingOfType("string")).Return(auth.ErrRefreshTokenNotFound)

	handler := commands.NewLogoutHandler(refreshRepo)
	err := handler.Handle(context.Background(), commands.LogoutCommand{
		RefreshToken: "unknown-token",
	})

	// Logout is idempotent — no error even if token was already deleted
	assert.NoError(t, err)
	refreshRepo.AssertExpectations(t)
}
