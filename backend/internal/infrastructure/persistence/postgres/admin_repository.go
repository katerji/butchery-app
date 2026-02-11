package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katerji/butchery-app/backend/internal/domain/admin"
)

// AdminRepository implements admin.Repository using PostgreSQL.
type AdminRepository struct {
	pool *pgxpool.Pool
}

// NewAdminRepository creates a new AdminRepository.
func NewAdminRepository(pool *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{pool: pool}
}

// FindByEmail finds an admin by email address.
func (r *AdminRepository) FindByEmail(ctx context.Context, email string) (*admin.Admin, error) {
	var id uuid.UUID
	var dbEmail, passwordHash, fullName string

	err := r.pool.QueryRow(ctx,
		"SELECT id, email, password_hash, full_name FROM admins WHERE email = $1",
		email,
	).Scan(&id, &dbEmail, &passwordHash, &fullName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, admin.ErrAdminNotFound
		}
		return nil, fmt.Errorf("querying admin by email: %w", err)
	}

	return admin.NewAdmin(id, dbEmail, passwordHash, fullName)
}

// FindByID finds an admin by ID.
func (r *AdminRepository) FindByID(ctx context.Context, id uuid.UUID) (*admin.Admin, error) {
	var dbEmail, passwordHash, fullName string

	err := r.pool.QueryRow(ctx,
		"SELECT email, password_hash, full_name FROM admins WHERE id = $1",
		id,
	).Scan(&dbEmail, &passwordHash, &fullName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, admin.ErrAdminNotFound
		}
		return nil, fmt.Errorf("querying admin by id: %w", err)
	}

	return admin.NewAdmin(id, dbEmail, passwordHash, fullName)
}
