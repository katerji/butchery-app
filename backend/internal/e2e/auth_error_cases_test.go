package e2e_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/katerji/butchery-app/backend/internal/interface/http/dto"
)

func TestIntegrationRegistration_DuplicateEmail_ReturnsConflict(t *testing.T) {
	ts := setupTestServer(t, 15*time.Minute)

	body := dto.RegisterCustomerRequest{
		Email:    "dup@example.com",
		Password: "securepassword123",
		FullName: "Jane Doe",
		Phone:    "+1234567890",
	}

	// First registration succeeds.
	resp := ts.postJSON(t, "/api/v1/auth/register", body)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// Second registration with same email returns 409.
	resp = ts.postJSON(t, "/api/v1/auth/register", body)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	errMsg := parseError(t, resp)
	assert.Equal(t, "email already exists", errMsg)
}

func TestIntegrationCustomerLogin_InvalidCredentials_ReturnsUnauthorized(t *testing.T) {
	ts := setupTestServer(t, 15*time.Minute)

	// Register a customer first.
	registerBody := dto.RegisterCustomerRequest{
		Email:    "cust@example.com",
		Password: "securepassword123",
		FullName: "Test User",
		Phone:    "+1234567890",
	}
	resp := ts.postJSON(t, "/api/v1/auth/register", registerBody)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	t.Run("wrong password", func(t *testing.T) {
		loginBody := dto.LoginRequest{Email: "cust@example.com", Password: "wrongpassword"}
		resp := ts.postJSON(t, "/api/v1/auth/login", loginBody)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		errMsg := parseError(t, resp)
		assert.Equal(t, "invalid credentials", errMsg)
	})

	t.Run("non-existent email", func(t *testing.T) {
		loginBody := dto.LoginRequest{Email: "nobody@example.com", Password: "securepassword123"}
		resp := ts.postJSON(t, "/api/v1/auth/login", loginBody)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		errMsg := parseError(t, resp)
		assert.Equal(t, "invalid credentials", errMsg)
	})
}

func TestIntegrationAdminLogin_InvalidCredentials_ReturnsUnauthorized(t *testing.T) {
	ts := setupTestServer(t, 15*time.Minute)

	t.Run("wrong password", func(t *testing.T) {
		loginBody := dto.LoginRequest{Email: testAdminEmail, Password: "wrongpassword"}
		resp := ts.postJSON(t, "/api/v1/admin/auth/login", loginBody)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		errMsg := parseError(t, resp)
		assert.Equal(t, "invalid credentials", errMsg)
	})

	t.Run("non-existent email", func(t *testing.T) {
		loginBody := dto.LoginRequest{Email: "fake@admin.com", Password: testAdminPassword}
		resp := ts.postJSON(t, "/api/v1/admin/auth/login", loginBody)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		errMsg := parseError(t, resp)
		assert.Equal(t, "invalid credentials", errMsg)
	})
}

func TestIntegrationAuth_ExpiredAccessToken_ReturnsUnauthorized(t *testing.T) {
	// Use a very short TTL so the token expires quickly.
	ts := setupTestServer(t, 1*time.Second)

	// Login as admin to get a short-lived access token.
	loginBody := dto.LoginRequest{Email: testAdminEmail, Password: testAdminPassword}
	resp := ts.postJSON(t, "/api/v1/admin/auth/login", loginBody)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp dto.LoginResponse
	parseJSON(t, resp, &loginResp)

	// Wait for the token to expire.
	time.Sleep(2 * time.Second)

	// Try to use the expired token on a protected endpoint.
	logoutBody := dto.LogoutRequest{RefreshToken: loginResp.RefreshToken}
	resp = ts.postJSONWithAuth(t, "/api/v1/auth/logout", logoutBody, loginResp.AccessToken)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	errMsg := parseError(t, resp)
	assert.Equal(t, "unauthorized", errMsg)
}

func TestIntegrationAuth_InvalidToken_ReturnsUnauthorized(t *testing.T) {
	ts := setupTestServer(t, 15*time.Minute)

	// Login to get a valid refresh token for the logout body.
	loginBody := dto.LoginRequest{Email: testAdminEmail, Password: testAdminPassword}
	resp := ts.postJSON(t, "/api/v1/admin/auth/login", loginBody)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp dto.LoginResponse
	parseJSON(t, resp, &loginResp)

	logoutBody := dto.LogoutRequest{RefreshToken: loginResp.RefreshToken}

	t.Run("garbage token", func(t *testing.T) {
		resp := ts.postJSONWithAuth(t, "/api/v1/auth/logout", logoutBody, "not-a-real-jwt-token")
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		errMsg := parseError(t, resp)
		assert.Equal(t, "unauthorized", errMsg)
	})

	t.Run("wrong secret", func(t *testing.T) {
		// Create a token signed with a different secret.
		wrongToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":  "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			"type": "admin",
			"exp":  time.Now().Add(15 * time.Minute).Unix(),
			"iat":  time.Now().Unix(),
		})
		signed, err := wrongToken.SignedString([]byte("completely-different-secret-key"))
		require.NoError(t, err)

		resp := ts.postJSONWithAuth(t, "/api/v1/auth/logout", logoutBody, signed)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		errMsg := parseError(t, resp)
		assert.Equal(t, "unauthorized", errMsg)
	})
}

func TestIntegrationAuth_RefreshWithInvalidToken_ReturnsUnauthorized(t *testing.T) {
	ts := setupTestServer(t, 15*time.Minute)

	refreshBody := dto.RefreshTokenRequest{RefreshToken: "nonexistent-refresh-token-value"}
	resp := ts.postJSON(t, "/api/v1/auth/refresh", refreshBody)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	errMsg := parseError(t, resp)
	assert.Equal(t, "invalid or expired refresh token", errMsg)
}

func TestIntegrationAuth_LogoutIdempotency(t *testing.T) {
	ts := setupTestServer(t, 15*time.Minute)

	// Login to get tokens.
	loginBody := dto.LoginRequest{Email: testAdminEmail, Password: testAdminPassword}
	resp := ts.postJSON(t, "/api/v1/admin/auth/login", loginBody)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp dto.LoginResponse
	parseJSON(t, resp, &loginResp)

	logoutBody := dto.LogoutRequest{RefreshToken: loginResp.RefreshToken}

	// First logout — should succeed.
	resp = ts.postJSONWithAuth(t, "/api/v1/auth/logout", logoutBody, loginResp.AccessToken)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()

	// Second logout with same token — should also succeed (idempotent).
	resp = ts.postJSONWithAuth(t, "/api/v1/auth/logout", logoutBody, loginResp.AccessToken)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()
}

func TestIntegrationRegistration_ValidationErrors(t *testing.T) {
	ts := setupTestServer(t, 15*time.Minute)

	t.Run("invalid email", func(t *testing.T) {
		body := dto.RegisterCustomerRequest{
			Email:    "not-an-email",
			Password: "securepassword123",
			FullName: "Test User",
			Phone:    "+1234567890",
		}
		resp := ts.postJSON(t, "/api/v1/auth/register", body)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		errMsg := parseError(t, resp)
		assert.Equal(t, "invalid email format", errMsg)
	})

	t.Run("short password", func(t *testing.T) {
		body := dto.RegisterCustomerRequest{
			Email:    "valid@example.com",
			Password: "short",
			FullName: "Test User",
			Phone:    "+1234567890",
		}
		resp := ts.postJSON(t, "/api/v1/auth/register", body)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		errMsg := parseError(t, resp)
		assert.Equal(t, "password must be at least 8 characters", errMsg)
	})

	t.Run("missing fields", func(t *testing.T) {
		body := map[string]string{
			"email": "valid@example.com",
			// password, full_name, phone missing
		}
		resp := ts.postJSON(t, "/api/v1/auth/register", body)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		errMsg := parseError(t, resp)
		assert.Equal(t, "all fields are required", errMsg)
	})
}
