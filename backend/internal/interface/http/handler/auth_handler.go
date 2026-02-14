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
//
//	@Summary		Refresh access token
//	@Description	Exchange a valid refresh token for a new access token.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.RefreshTokenRequest		true	"Refresh token"
//	@Success		200		{object}	dto.RefreshSuccessResponse	"New access token"
//	@Failure		400		{object}	dto.ErrorBody				"Invalid request body"
//	@Failure		401		{object}	dto.ErrorBody				"Invalid or expired refresh token"
//	@Router			/auth/refresh [post]
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
//
//	@Summary		Logout
//	@Description	Revoke a refresh token, effectively logging the user out.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.LogoutRequest	true	"Refresh token to revoke"
//	@Success		204		"Successfully logged out"
//	@Failure		400		{object}	dto.ErrorBody		"Invalid request body"
//	@Failure		401		{object}	dto.ErrorBody		"Unauthorized"
//	@Failure		500		{object}	dto.ErrorBody		"Internal server error"
//	@Router			/auth/logout [post]
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
