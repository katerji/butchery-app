package e2e_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/katerji/butchery-app/backend/internal/interface/http/dto"
)

func TestIntegrationAdminAuth_LoginAndLogout(t *testing.T) {
	ts := setupTestServer(t, 15*time.Minute)

	// Step 1: Login as seeded admin.
	loginBody := dto.LoginRequest{
		Email:    testAdminEmail,
		Password: testAdminPassword,
	}
	resp := ts.postJSON(t, "/api/v1/admin/auth/login", loginBody)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp dto.LoginResponse
	parseJSON(t, resp, &loginResp)
	assert.NotEmpty(t, loginResp.AccessToken)
	assert.NotEmpty(t, loginResp.RefreshToken)
	assert.Greater(t, loginResp.ExpiresIn, int64(0))

	// Step 2: Logout.
	logoutBody := dto.LogoutRequest{RefreshToken: loginResp.RefreshToken}
	resp = ts.postJSONWithAuth(t, "/api/v1/auth/logout", logoutBody, loginResp.AccessToken)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()

	// Step 3: Refresh with invalidated token — should be 401.
	refreshBody := dto.RefreshTokenRequest{RefreshToken: loginResp.RefreshToken}
	resp = ts.postJSON(t, "/api/v1/auth/refresh", refreshBody)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	errMsg := parseError(t, resp)
	assert.Equal(t, "invalid or expired refresh token", errMsg)
}

func TestIntegrationAdminAuth_LoginAndRefresh(t *testing.T) {
	ts := setupTestServer(t, 15*time.Minute)

	// Step 1: Login.
	loginBody := dto.LoginRequest{
		Email:    testAdminEmail,
		Password: testAdminPassword,
	}
	resp := ts.postJSON(t, "/api/v1/admin/auth/login", loginBody)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp dto.LoginResponse
	parseJSON(t, resp, &loginResp)

	// Step 2: Refresh — get new access token.
	refreshBody := dto.RefreshTokenRequest{RefreshToken: loginResp.RefreshToken}
	resp = ts.postJSON(t, "/api/v1/auth/refresh", refreshBody)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var refreshResp dto.RefreshTokenResponse
	parseJSON(t, resp, &refreshResp)
	assert.NotEmpty(t, refreshResp.AccessToken)

	// Step 3: Use new access token for logout.
	logoutBody := dto.LogoutRequest{RefreshToken: loginResp.RefreshToken}
	resp = ts.postJSONWithAuth(t, "/api/v1/auth/logout", logoutBody, refreshResp.AccessToken)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()
}
