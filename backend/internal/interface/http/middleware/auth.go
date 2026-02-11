package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/katerji/butchery-app/backend/internal/domain/auth"
	"github.com/katerji/butchery-app/backend/pkg/httpresponse"
)

type contextKey string

const claimsKey contextKey = "claims"

// AuthMiddleware provides JWT authentication middleware.
type AuthMiddleware struct {
	validator auth.TokenValidator
}

// NewAuthMiddleware creates a new AuthMiddleware.
func NewAuthMiddleware(validator auth.TokenValidator) *AuthMiddleware {
	return &AuthMiddleware{validator: validator}
}

// RequireAuth validates the JWT and injects claims into context.
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := m.extractClaims(r)
		if err != nil {
			httpresponse.Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin validates the JWT and ensures the subject is an admin.
func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := m.extractClaims(r)
		if err != nil {
			httpresponse.Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if claims.SubjectType != auth.SubjectTypeAdmin {
			httpresponse.Error(w, http.StatusForbidden, "forbidden")
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireCustomer validates the JWT and ensures the subject is a customer.
func (m *AuthMiddleware) RequireCustomer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := m.extractClaims(r)
		if err != nil {
			httpresponse.Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if claims.SubjectType != auth.SubjectTypeCustomer {
			httpresponse.Error(w, http.StatusForbidden, "forbidden")
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) extractClaims(r *http.Request) (*auth.AccessTokenClaims, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return nil, fmt.Errorf("missing authorization header")
	}

	token := strings.TrimPrefix(header, "Bearer ")
	if token == header {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	return m.validator.ValidateAccessToken(token)
}

// ClaimsFromContext extracts access token claims from the request context.
func ClaimsFromContext(ctx context.Context) *auth.AccessTokenClaims {
	claims, _ := ctx.Value(claimsKey).(*auth.AccessTokenClaims)
	return claims
}
