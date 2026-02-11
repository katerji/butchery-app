package commands

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domainauth "github.com/katerji/butchery-app/backend/internal/domain/auth"
	"github.com/katerji/butchery-app/backend/internal/domain/customer"
)

// RegisterCustomerCommand is the input for the customer registration use case.
type RegisterCustomerCommand struct {
	Email    string
	Password string
	FullName string
	Phone    string
}

// RegisterCustomerResult is the output of the customer registration use case.
type RegisterCustomerResult struct {
	CustomerID uuid.UUID
	Email      string
	FullName   string
}

// RegisterCustomerHandler handles customer registration.
type RegisterCustomerHandler struct {
	customerRepo customer.Repository
	hasher       domainauth.PasswordHasher
}

// NewRegisterCustomerHandler creates a new RegisterCustomerHandler.
func NewRegisterCustomerHandler(
	customerRepo customer.Repository,
	hasher domainauth.PasswordHasher,
) *RegisterCustomerHandler {
	return &RegisterCustomerHandler{
		customerRepo: customerRepo,
		hasher:       hasher,
	}
}

// Handle executes the customer registration use case.
func (h *RegisterCustomerHandler) Handle(ctx context.Context, cmd RegisterCustomerCommand) (*RegisterCustomerResult, error) {
	email, err := customer.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}

	if _, err := customer.NewPassword(cmd.Password); err != nil {
		return nil, err
	}

	phone, err := customer.NewPhoneNumber(cmd.Phone)
	if err != nil {
		return nil, err
	}

	exists, err := h.customerRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("checking email existence: %w", err)
	}
	if exists {
		return nil, customer.ErrEmailAlreadyExists
	}

	hashedPassword, err := h.hasher.Hash(cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	id := uuid.New()
	c, err := customer.NewCustomer(id, email, hashedPassword, cmd.FullName, phone)
	if err != nil {
		return nil, err
	}

	if err := h.customerRepo.Save(ctx, c); err != nil {
		return nil, fmt.Errorf("saving customer: %w", err)
	}

	return &RegisterCustomerResult{
		CustomerID: c.ID(),
		Email:      c.Email().String(),
		FullName:   c.FullName(),
	}, nil
}
