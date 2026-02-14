package handler

import (
	"encoding/json"
	"net/http"

	"github.com/katerji/butchery-app/backend/internal/application/admin/commands"
	"github.com/katerji/butchery-app/backend/internal/interface/http/dto"
	"github.com/katerji/butchery-app/backend/pkg/httpresponse"
)

// AdminAuthHandler handles admin authentication HTTP requests.
type AdminAuthHandler struct {
	loginHandler *commands.AdminLoginHandler
}

// NewAdminAuthHandler creates a new AdminAuthHandler.
func NewAdminAuthHandler(loginHandler *commands.AdminLoginHandler) *AdminAuthHandler {
	return &AdminAuthHandler{loginHandler: loginHandler}
}

// Login handles POST /api/v1/admin/auth/login.
//
//	@Summary		Admin login
//	@Description	Authenticate an admin with email and password. Returns JWT access and refresh tokens.
//	@Tags			Admin Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.LoginRequest			true	"Admin credentials"
//	@Success		200		{object}	dto.LoginSuccessResponse		"Successful login"
//	@Failure		400		{object}	dto.ErrorBody				"Invalid request body"
//	@Failure		401		{object}	dto.ErrorBody				"Invalid credentials"
//	@Router			/admin/auth/login [post]
func (h *AdminAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		httpresponse.Error(w, http.StatusBadRequest, "email and password are required")
		return
	}

	result, err := h.loginHandler.Handle(r.Context(), commands.AdminLoginCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		httpresponse.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	httpresponse.Success(w, dto.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	})
}
