package commands

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/katerji/butchery-app/backend/internal/application/auth"
	"github.com/katerji/butchery-app/backend/internal/domain/admin"
	domainauth "github.com/katerji/butchery-app/backend/internal/domain/auth"
)

const refreshTokenTTL = 7 * 24 * time.Hour

// AdminLoginCommand is the input for the admin login use case.
type AdminLoginCommand struct {
	Email    string
	Password string
}

// AdminLoginHandler handles admin login.
type AdminLoginHandler struct {
	adminRepo       admin.Repository
	hasher          domainauth.PasswordHasher
	tokenGen        domainauth.TokenGenerator
	refreshRepo     domainauth.RefreshTokenRepository
	accessTokenTTL  time.Duration
}

// NewAdminLoginHandler creates a new AdminLoginHandler with its dependencies.
func NewAdminLoginHandler(
	adminRepo admin.Repository,
	hasher domainauth.PasswordHasher,
	tokenGen domainauth.TokenGenerator,
	refreshRepo domainauth.RefreshTokenRepository,
	accessTokenTTL time.Duration,
) *AdminLoginHandler {
	return &AdminLoginHandler{
		adminRepo:      adminRepo,
		hasher:         hasher,
		tokenGen:       tokenGen,
		refreshRepo:    refreshRepo,
		accessTokenTTL: accessTokenTTL,
	}
}

// Handle executes the admin login use case.
func (h *AdminLoginHandler) Handle(ctx context.Context, cmd AdminLoginCommand) (*auth.LoginResult, error) {
	a, err := h.adminRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("%w", admin.ErrInvalidCredentials)
	}

	if err := h.hasher.Compare(a.PasswordHash(), cmd.Password); err != nil {
		return nil, fmt.Errorf("%w", admin.ErrInvalidCredentials)
	}

	accessToken, err := h.tokenGen.GenerateAccessToken(a.ID(), domainauth.SubjectTypeAdmin)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	rawRefresh, err := h.tokenGen.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	tokenHash := hashToken(rawRefresh)
	expiresAt := time.Now().Add(refreshTokenTTL)

	refreshToken, err := domainauth.NewRefreshToken(a.ID(), domainauth.SubjectTypeAdmin, tokenHash, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("creating refresh token: %w", err)
	}

	if err := h.refreshRepo.Save(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("saving refresh token: %w", err)
	}

	return &auth.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresIn:    int64(h.accessTokenTTL / time.Second),
	}, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
