package dto

// LoginRequest is the request body for login endpoints.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is the response body for login endpoints.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RegisterCustomerRequest is the request body for customer registration.
type RegisterCustomerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

// RegisterCustomerResponse is the response body for customer registration.
type RegisterCustomerResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

// RefreshTokenRequest is the request body for token refresh.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse is the response body for token refresh.
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// LogoutRequest is the request body for logout.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}
