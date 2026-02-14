package customer

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Customer is the aggregate root for the customer bounded context.
type Customer struct {
	id           uuid.UUID
	email        Email
	passwordHash string
	fullName     string
	phone        PhoneNumber
	createdAt    time.Time
	updatedAt    time.Time
}

// NewCustomer creates a new Customer aggregate with invariant validation.
func NewCustomer(id uuid.UUID, email Email, passwordHash string, fullName string, phone PhoneNumber) (*Customer, error) {
	if strings.TrimSpace(fullName) == "" {
		return nil, ErrEmptyFullName
	}

	now := time.Now()
	return &Customer{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		fullName:     fullName,
		phone:        phone,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// ReconstructCustomer reconstructs a Customer from persistence without validation.
func ReconstructCustomer(id uuid.UUID, email Email, passwordHash string, fullName string, phone PhoneNumber) *Customer {
	return &Customer{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		fullName:     fullName,
		phone:        phone,
	}
}

func (c *Customer) ID() uuid.UUID       { return c.id }
func (c *Customer) Email() Email         { return c.email }
func (c *Customer) PasswordHash() string { return c.passwordHash }
func (c *Customer) FullName() string     { return c.fullName }
func (c *Customer) Phone() PhoneNumber   { return c.phone }
func (c *Customer) CreatedAt() time.Time { return c.createdAt }
func (c *Customer) UpdatedAt() time.Time { return c.updatedAt }
