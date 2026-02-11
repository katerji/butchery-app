package handler

import (
	"encoding/json"
	"net/http"

	authcmd "github.com/katerji/butchery-app/backend/internal/application/auth/commands"
	"github.com/katerji/butchery-app/backend/internal/interface/http/dto"
	"github.com/katerji/butchery-app/backend/pkg/httpresponse"
)

// AuthHandler handles shared authentication HTTP requests (refresh, logout).
type AuthHandler struct {
	refreshHandler *authcmd.RefreshTokenHandler
	logoutHandler  *authcmd.LogoutHandler
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(
	refreshHandler *authcmd.RefreshTokenHandler,
	logoutHandler *authcmd.LogoutHandler,
) *AuthHandler {
	return &AuthHandler{
		refreshHandler: refreshHandler,
		logoutHandler:  logoutHandler,
	}
}

// Refresh handles POST /api/v1/auth/refresh.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		httpresponse.Error(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	result, err := h.refreshHandler.Handle(r.Context(), authcmd.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		httpresponse.Error(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	httpresponse.Success(w, dto.RefreshTokenResponse{
		AccessToken: result.AccessToken,
		ExpiresIn:   result.ExpiresIn,
	})
}

// Logout handles POST /api/v1/auth/logout.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		httpresponse.Error(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	if err := h.logoutHandler.Handle(r.Context(), authcmd.LogoutCommand{
		RefreshToken: req.RefreshToken,
	}); err != nil {
		httpresponse.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	httpresponse.NoContent(w)
}
