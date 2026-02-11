package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	domainauth "github.com/katerji/butchery-app/backend/internal/domain/auth"
)

// TokenService implements auth.TokenGenerator and auth.TokenValidator.
type TokenService struct {
	secret         []byte
	accessTokenTTL time.Duration
}

// NewTokenService creates a new TokenService.
func NewTokenService(secret string, accessTokenTTL time.Duration) *TokenService {
	return &TokenService{
		secret:         []byte(secret),
		accessTokenTTL: accessTokenTTL,
	}
}

// GenerateAccessToken generates a signed JWT access token.
func (s *TokenService) GenerateAccessToken(subjectID uuid.UUID, subjectType string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  subjectID.String(),
		"type": subjectType,
		"exp":  now.Add(s.accessTokenTTL).Unix(),
		"iat":  now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("signing access token: %w", err)
	}
	return signed, nil
}

// GenerateRefreshToken generates a cryptographically random opaque refresh token.
func (s *TokenService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generating refresh token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// ValidateAccessToken validates a JWT access token and returns the claims.
func (s *TokenService) ValidateAccessToken(tokenString string) (*domainauth.AccessTokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing access token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("missing sub claim")
	}

	subjectID, err := uuid.Parse(sub)
	if err != nil {
		return nil, fmt.Errorf("parsing subject ID: %w", err)
	}

	subjectType, ok := claims["type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing type claim")
	}

	return &domainauth.AccessTokenClaims{
		SubjectID:   subjectID,
		SubjectType: subjectType,
	}, nil
}

// HashRefreshToken hashes a raw refresh token using SHA256.
func HashRefreshToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
