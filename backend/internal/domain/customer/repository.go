package customer

import (
	"context"

	"github.com/google/uuid"
)

// Repository provides access to customer persistence.
type Repository interface {
	Save(ctx context.Context, customer *Customer) error
	FindByEmail(ctx context.Context, email Email) (*Customer, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Customer, error)
	ExistsByEmail(ctx context.Context, email Email) (bool, error)
}
