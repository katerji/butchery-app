package commands

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	appauth "github.com/katerji/butchery-app/backend/internal/application/auth"
	"github.com/katerji/butchery-app/backend/internal/domain/auth"
)

// RefreshTokenCommand is the input for the refresh token use case.
type RefreshTokenCommand struct {
	RefreshToken string
}

// RefreshTokenHandler handles refreshing access tokens.
type RefreshTokenHandler struct {
	refreshRepo auth.RefreshTokenRepository
	tokenGen    auth.TokenGenerator
}

// NewRefreshTokenHandler creates a new RefreshTokenHandler.
func NewRefreshTokenHandler(
	refreshRepo auth.RefreshTokenRepository,
	tokenGen auth.TokenGenerator,
) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		refreshRepo: refreshRepo,
		tokenGen:    tokenGen,
	}
}

// Handle executes the refresh token use case.
func (h *RefreshTokenHandler) Handle(ctx context.Context, cmd RefreshTokenCommand) (*appauth.RefreshTokenResult, error) {
	tokenHash := hashToken(cmd.RefreshToken)

	storedToken, err := h.refreshRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	if storedToken.IsExpired() {
		return nil, fmt.Errorf("%w", auth.ErrRefreshTokenExpired)
	}

	accessToken, err := h.tokenGen.GenerateAccessToken(storedToken.SubjectID(), storedToken.SubjectType())
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	return &appauth.RefreshTokenResult{
		AccessToken: accessToken,
		ExpiresIn:   int64(15 * time.Minute / time.Second),
	}, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
