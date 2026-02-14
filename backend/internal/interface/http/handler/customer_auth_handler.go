package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	custcmd "github.com/katerji/butchery-app/backend/internal/application/customer/commands"
	"github.com/katerji/butchery-app/backend/internal/domain/customer"
	"github.com/katerji/butchery-app/backend/internal/interface/http/dto"
	"github.com/katerji/butchery-app/backend/pkg/httpresponse"
)

// CustomerAuthHandler handles customer authentication HTTP requests.
type CustomerAuthHandler struct {
	registerHandler *custcmd.RegisterCustomerHandler
	loginHandler    *custcmd.CustomerLoginHandler
}

// NewCustomerAuthHandler creates a new CustomerAuthHandler.
func NewCustomerAuthHandler(
	registerHandler *custcmd.RegisterCustomerHandler,
	loginHandler *custcmd.CustomerLoginHandler,
) *CustomerAuthHandler {
	return &CustomerAuthHandler{
		registerHandler: registerHandler,
		loginHandler:    loginHandler,
	}
}

// Register handles POST /api/v1/auth/register.
func (h *CustomerAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.FullName == "" || req.Phone == "" {
		httpresponse.Error(w, http.StatusBadRequest, "all fields are required")
		return
	}

	result, err := h.registerHandler.Handle(r.Context(), custcmd.RegisterCustomerCommand{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Phone:    req.Phone,
	})
	if err != nil {
		switch {
		case errors.Is(err, customer.ErrEmailAlreadyExists):
			httpresponse.Error(w, http.StatusConflict, "email already exists")
		case errors.Is(err, customer.ErrInvalidEmail):
			httpresponse.Error(w, http.StatusUnprocessableEntity, "invalid email format")
		case errors.Is(err, customer.ErrInvalidPassword):
			httpresponse.Error(w, http.StatusUnprocessableEntity, "password must be at least 8 characters")
		case errors.Is(err, customer.ErrInvalidPhoneNumber):
			httpresponse.Error(w, http.StatusUnprocessableEntity, "invalid phone number")
		default:
			httpresponse.Error(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	httpresponse.Created(w, dto.RegisterCustomerResponse{
		ID:       result.CustomerID.String(),
		Email:    result.Email,
		FullName: result.FullName,
	})
}

// Login handles POST /api/v1/auth/login.
func (h *CustomerAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		httpresponse.Error(w, http.StatusBadRequest, "email and password are required")
		return
	}

	result, err := h.loginHandler.Handle(r.Context(), custcmd.CustomerLoginCommand{
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
