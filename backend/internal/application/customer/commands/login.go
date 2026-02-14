package commands

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/katerji/butchery-app/backend/internal/application/auth"
	domainauth "github.com/katerji/butchery-app/backend/internal/domain/auth"
	"github.com/katerji/butchery-app/backend/internal/domain/customer"
)

const refreshTokenTTL = 7 * 24 * time.Hour

// CustomerLoginCommand is the input for the customer login use case.
type CustomerLoginCommand struct {
	Email    string
	Password string
}

// CustomerLoginHandler handles customer login.
type CustomerLoginHandler struct {
	customerRepo   customer.Repository
	hasher         domainauth.PasswordHasher
	tokenGen       domainauth.TokenGenerator
	refreshRepo    domainauth.RefreshTokenRepository
	accessTokenTTL time.Duration
}

// NewCustomerLoginHandler creates a new CustomerLoginHandler.
func NewCustomerLoginHandler(
	customerRepo customer.Repository,
	hasher domainauth.PasswordHasher,
	tokenGen domainauth.TokenGenerator,
	refreshRepo domainauth.RefreshTokenRepository,
	accessTokenTTL time.Duration,
) *CustomerLoginHandler {
	return &CustomerLoginHandler{
		customerRepo:   customerRepo,
		hasher:         hasher,
		tokenGen:       tokenGen,
		refreshRepo:    refreshRepo,
		accessTokenTTL: accessTokenTTL,
	}
}

// Handle executes the customer login use case.
func (h *CustomerLoginHandler) Handle(ctx context.Context, cmd CustomerLoginCommand) (*auth.LoginResult, error) {
	email, err := customer.NewEmail(cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("%w", customer.ErrInvalidCredentials)
	}

	c, err := h.customerRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%w", customer.ErrInvalidCredentials)
	}

	if err := h.hasher.Compare(c.PasswordHash(), cmd.Password); err != nil {
		return nil, fmt.Errorf("%w", customer.ErrInvalidCredentials)
	}

	accessToken, err := h.tokenGen.GenerateAccessToken(c.ID(), domainauth.SubjectTypeCustomer)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	rawRefresh, err := h.tokenGen.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	tokenHash := hashToken(rawRefresh)
	expiresAt := time.Now().Add(refreshTokenTTL)

	refreshToken, err := domainauth.NewRefreshToken(c.ID(), domainauth.SubjectTypeCustomer, tokenHash, expiresAt)
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
