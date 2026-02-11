package auth

// LoginResult is the output of any login use case (admin or customer).
type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // seconds until access token expires
}

// RefreshTokenResult is the output of the refresh token use case.
type RefreshTokenResult struct {
	AccessToken string
	ExpiresIn   int64
}
