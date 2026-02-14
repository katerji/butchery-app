package admin

import (
	"context"

	"github.com/google/uuid"
)

// Repository provides access to admin persistence.
type Repository interface {
	FindByEmail(ctx context.Context, email string) (*Admin, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Admin, error)
}
