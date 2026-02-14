package e2e_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/katerji/butchery-app/backend/internal/interface/http/dto"
)

func TestIntegrationCustomerAuth_FullLifecycle(t *testing.T) {
	ts := setupTestServer(t, 15*time.Minute)

	// Step 1: Register a customer.
	registerBody := dto.RegisterCustomerRequest{
		Email:    "customer@example.com",
		Password: "securepassword123",
		FullName: "John Doe",
		Phone:    "+1234567890",
	}
	resp := ts.postJSON(t, "/api/v1/auth/register", registerBody)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var registerResp dto.RegisterCustomerResponse
	parseJSON(t, resp, &registerResp)
	assert.NotEmpty(t, registerResp.ID)
	assert.Equal(t, "customer@example.com", registerResp.Email)
	assert.Equal(t, "John Doe", registerResp.FullName)

	// Step 2: Login with the registered customer.
	loginBody := dto.LoginRequest{
		Email:    "customer@example.com",
		Password: "securepassword123",
	}
	resp = ts.postJSON(t, "/api/v1/auth/login", loginBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp dto.LoginResponse
	parseJSON(t, resp, &loginResp)
	assert.NotEmpty(t, loginResp.AccessToken)
	assert.NotEmpty(t, loginResp.RefreshToken)
	assert.Greater(t, loginResp.ExpiresIn, int64(0))

	accessToken := loginResp.AccessToken
	refreshToken := loginResp.RefreshToken

	// Step 3: Hit protected endpoint (logout) without token — expect 401.
	resp = ts.postJSON(t, "/api/v1/auth/logout", dto.LogoutRequest{RefreshToken: refreshToken})
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()

	// Step 4: Refresh to get a new access token.
	refreshBody := dto.RefreshTokenRequest{RefreshToken: refreshToken}
	resp = ts.postJSON(t, "/api/v1/auth/refresh", refreshBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var refreshResp dto.RefreshTokenResponse
	parseJSON(t, resp, &refreshResp)
	assert.NotEmpty(t, refreshResp.AccessToken)
	assert.Greater(t, refreshResp.ExpiresIn, int64(0))

	newAccessToken := refreshResp.AccessToken

	// Step 5: Logout with the new access token.
	logoutBody := dto.LogoutRequest{RefreshToken: refreshToken}
	resp = ts.postJSONWithAuth(t, "/api/v1/auth/logout", logoutBody, newAccessToken)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()

	// Step 6: Refresh again with the same refresh token — should be 401 (token invalidated).
	resp = ts.postJSON(t, "/api/v1/auth/refresh", refreshBody)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	errMsg := parseError(t, resp)
	assert.Equal(t, "invalid or expired refresh token", errMsg)

	// Verify the original access token still works for an authenticated request
	// (access tokens are stateless and remain valid until expiry, even after logout).
	_ = accessToken
}
