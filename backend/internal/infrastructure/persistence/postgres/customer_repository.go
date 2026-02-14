package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katerji/butchery-app/backend/internal/domain/customer"
)

// CustomerRepository implements customer.Repository using PostgreSQL.
type CustomerRepository struct {
	pool *pgxpool.Pool
}

// NewCustomerRepository creates a new CustomerRepository.
func NewCustomerRepository(pool *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{pool: pool}
}

// Save persists a new customer.
func (r *CustomerRepository) Save(ctx context.Context, c *customer.Customer) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO customers (id, email, password_hash, full_name, phone_number, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		c.ID(), c.Email().String(), c.PasswordHash(), c.FullName(), c.Phone().String(),
		c.CreatedAt(), c.UpdatedAt(),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return customer.ErrEmailAlreadyExists
		}
		return fmt.Errorf("inserting customer: %w", err)
	}
	return nil
}

// FindByEmail finds a customer by email.
func (r *CustomerRepository) FindByEmail(ctx context.Context, email customer.Email) (*customer.Customer, error) {
	var id uuid.UUID
	var dbEmail, passwordHash, fullName, phoneNumber string

	err := r.pool.QueryRow(ctx,
		"SELECT id, email, password_hash, full_name, phone_number FROM customers WHERE email = $1",
		email.String(),
	).Scan(&id, &dbEmail, &passwordHash, &fullName, &phoneNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customer.ErrCustomerNotFound
		}
		return nil, fmt.Errorf("querying customer by email: %w", err)
	}

	return reconstructCustomer(id, dbEmail, passwordHash, fullName, phoneNumber)
}

// FindByID finds a customer by ID.
func (r *CustomerRepository) FindByID(ctx context.Context, id uuid.UUID) (*customer.Customer, error) {
	var dbEmail, passwordHash, fullName, phoneNumber string

	err := r.pool.QueryRow(ctx,
		"SELECT email, password_hash, full_name, phone_number FROM customers WHERE id = $1",
		id,
	).Scan(&dbEmail, &passwordHash, &fullName, &phoneNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customer.ErrCustomerNotFound
		}
		return nil, fmt.Errorf("querying customer by id: %w", err)
	}

	return reconstructCustomer(id, dbEmail, passwordHash, fullName, phoneNumber)
}

// ExistsByEmail checks if a customer with the given email already exists.
func (r *CustomerRepository) ExistsByEmail(ctx context.Context, email customer.Email) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM customers WHERE email = $1)",
		email.String(),
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking customer existence: %w", err)
	}
	return exists, nil
}

func reconstructCustomer(id uuid.UUID, emailStr, passwordHash, fullName, phoneStr string) (*customer.Customer, error) {
	email, err := customer.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("reconstructing email: %w", err)
	}
	phone, err := customer.NewPhoneNumber(phoneStr)
	if err != nil {
		return nil, fmt.Errorf("reconstructing phone: %w", err)
	}
	return customer.ReconstructCustomer(id, email, passwordHash, fullName, phone), nil
}
