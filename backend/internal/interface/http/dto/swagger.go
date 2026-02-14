package dto

// Swagger documentation-only types. These are never instantiated at runtime;
// they exist so swag can generate concrete envelope schemas for responses.

// LoginSuccessResponse wraps LoginResponse in the standard API envelope.
type LoginSuccessResponse struct {
	Data  LoginResponse `json:"data"`
	Error *string       `json:"error"`
}

// RegisterSuccessResponse wraps RegisterCustomerResponse in the standard API envelope.
type RegisterSuccessResponse struct {
	Data  RegisterCustomerResponse `json:"data"`
	Error *string                  `json:"error"`
}

// RefreshSuccessResponse wraps RefreshTokenResponse in the standard API envelope.
type RefreshSuccessResponse struct {
	Data  RefreshTokenResponse `json:"data"`
	Error *string              `json:"error"`
}

// ErrorBody is the standard error envelope returned by the API.
type ErrorBody struct {
	Data  *string `json:"data"`
	Error string  `json:"error"`
}
