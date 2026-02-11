package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/katerji/butchery-app/backend/internal/domain/auth"
	"github.com/katerji/butchery-app/backend/internal/interface/http/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTokenValidator struct {
	mock.Mock
}

func (m *mockTokenValidator) ValidateAccessToken(token string) (*auth.AccessTokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.AccessTokenClaims), args.Error(1)
}

func TestRequireAuth_ValidToken_PassesThrough(t *testing.T) {
	validator := new(mockTokenValidator)
	mw := middleware.NewAuthMiddleware(validator)

	claims := &auth.AccessTokenClaims{SubjectID: uuid.New(), SubjectType: "customer"}
	validator.On("ValidateAccessToken", "valid-token").Return(claims, nil)

	handler := mw.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := middleware.ClaimsFromContext(r.Context())
		assert.NotNil(t, c)
		assert.Equal(t, claims.SubjectID, c.SubjectID)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRequireAuth_MissingToken_Returns401(t *testing.T) {
	validator := new(mockTokenValidator)
	mw := middleware.NewAuthMiddleware(validator)

	handler := mw.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestRequireAuth_InvalidToken_Returns401(t *testing.T) {
	validator := new(mockTokenValidator)
	mw := middleware.NewAuthMiddleware(validator)

	validator.On("ValidateAccessToken", "bad-token").Return(nil, errors.New("invalid"))

	handler := mw.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestRequireAdmin_AdminToken_PassesThrough(t *testing.T) {
	validator := new(mockTokenValidator)
	mw := middleware.NewAuthMiddleware(validator)

	claims := &auth.AccessTokenClaims{SubjectID: uuid.New(), SubjectType: "admin"}
	validator.On("ValidateAccessToken", "admin-token").Return(claims, nil)

	handler := mw.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRequireAdmin_CustomerToken_Returns403(t *testing.T) {
	validator := new(mockTokenValidator)
	mw := middleware.NewAuthMiddleware(validator)

	claims := &auth.AccessTokenClaims{SubjectID: uuid.New(), SubjectType: "customer"}
	validator.On("ValidateAccessToken", "customer-token").Return(claims, nil)

	handler := mw.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer customer-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestRequireCustomer_CustomerToken_PassesThrough(t *testing.T) {
	validator := new(mockTokenValidator)
	mw := middleware.NewAuthMiddleware(validator)

	claims := &auth.AccessTokenClaims{SubjectID: uuid.New(), SubjectType: "customer"}
	validator.On("ValidateAccessToken", "customer-token").Return(claims, nil)

	handler := mw.RequireCustomer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer customer-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRequireCustomer_AdminToken_Returns403(t *testing.T) {
	validator := new(mockTokenValidator)
	mw := middleware.NewAuthMiddleware(validator)

	claims := &auth.AccessTokenClaims{SubjectID: uuid.New(), SubjectType: "admin"}
	validator.On("ValidateAccessToken", "admin-token").Return(claims, nil)

	handler := mw.RequireCustomer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
}
