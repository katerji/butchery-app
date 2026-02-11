package commands

import (
	"context"
	"errors"

	"github.com/katerji/butchery-app/backend/internal/domain/auth"
)

// LogoutCommand is the input for the logout use case.
type LogoutCommand struct {
	RefreshToken string
}

// LogoutHandler handles user logout by invalidating refresh tokens.
type LogoutHandler struct {
	refreshRepo auth.RefreshTokenRepository
}

// NewLogoutHandler creates a new LogoutHandler.
func NewLogoutHandler(refreshRepo auth.RefreshTokenRepository) *LogoutHandler {
	return &LogoutHandler{refreshRepo: refreshRepo}
}

// Handle executes the logout use case. It is idempotent.
func (h *LogoutHandler) Handle(ctx context.Context, cmd LogoutCommand) error {
	tokenHash := hashToken(cmd.RefreshToken)

	err := h.refreshRepo.DeleteByTokenHash(ctx, tokenHash)
	if err != nil && !errors.Is(err, auth.ErrRefreshTokenNotFound) {
		return err
	}

	return nil
}
