package admin

import (
	"strings"

	"github.com/google/uuid"
)

// Admin represents a butchery administrator (back-office user).
type Admin struct {
	id           uuid.UUID
	email        string
	passwordHash string
	fullName     string
}

// NewAdmin creates an Admin entity with validation.
func NewAdmin(id uuid.UUID, email, passwordHash, fullName string) (*Admin, error) {
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}
	if strings.TrimSpace(fullName) == "" {
		return nil, ErrEmptyFullName
	}

	return &Admin{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		fullName:     fullName,
	}, nil
}

func (a *Admin) ID() uuid.UUID       { return a.id }
func (a *Admin) Email() string        { return a.email }
func (a *Admin) PasswordHash() string { return a.passwordHash }
func (a *Admin) FullName() string     { return a.fullName }

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return false
	}
	return parts[0] != "" && parts[1] != ""
}
