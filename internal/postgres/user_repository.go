package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/domain/user"
)

// UserRepository is a PostgreSQL-backed user repository.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a UserRepository backed by the given DB.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	const q = `
		INSERT INTO users (id, email, phone, role, is_verified, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, q,
		u.ID.String(), u.Email, u.Phone, string(u.Role),
		u.IsVerified, u.IsActive, u.CreatedAt.UTC(), u.UpdatedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id ulid.ULID) (*user.User, error) {
	const q = `
		SELECT id, email, phone, role, is_verified, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`
	return r.scanUser(r.db.QueryRowContext(ctx, q, id.String()))
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	const q = `
		SELECT id, email, phone, role, is_verified, is_active, created_at, updated_at
		FROM users WHERE email = $1
	`
	return r.scanUser(r.db.QueryRowContext(ctx, q, email))
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
	const q = `
		SELECT id, email, phone, role, is_verified, is_active, created_at, updated_at
		FROM users WHERE phone = $1
	`
	return r.scanUser(r.db.QueryRowContext(ctx, q, phone))
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	const q = `
		UPDATE users
		SET email=$2, phone=$3, role=$4, is_verified=$5, is_active=$6, updated_at=$7
		WHERE id=$1
	`
	res, err := r.db.ExecContext(ctx, q,
		u.ID.String(), u.Email, u.Phone, string(u.Role),
		u.IsVerified, u.IsActive, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return r.checkRowsAffected(res, user.ErrUserNotFound)
}

func (r *UserRepository) UpdateVerified(ctx context.Context, id ulid.ULID, verified bool) error {
	const q = `UPDATE users SET is_verified=$2, updated_at=$3 WHERE id=$1`
	res, err := r.db.ExecContext(ctx, q, id.String(), verified, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("update verified: %w", err)
	}
	return r.checkRowsAffected(res, user.ErrUserNotFound)
}

func (r *UserRepository) scanUser(row *sql.Row) (*user.User, error) {
	var u user.User
	var idStr string
	var role string
	var createdAt, updatedAt time.Time

	err := row.Scan(&idStr, &u.Email, &u.Phone, &role, &u.IsVerified, &u.IsActive, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("scan user: %w", err)
	}

	id, err := ulid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	u.ID = id
	u.Role = user.Role(role)
	u.CreatedAt = createdAt.UTC()
	u.UpdatedAt = updatedAt.UTC()

	return &u, nil
}

func (r *UserRepository) checkRowsAffected(res sql.Result, notFoundErr error) error {
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return notFoundErr
	}
	return nil
}
