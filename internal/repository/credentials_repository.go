package repository

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

// credentialsRepo is a PostgreSQL-backed credentials repository.
type credentialsRepo struct {
	db *sqlx.DB
}

// NewCredentialsRepository creates a CredentialsRepository backed by the given DB.
func NewCredentialsRepository(db *sqlx.DB) CredentialsRepository {
	return &credentialsRepo{db: db}
}

func (r *credentialsRepo) Create(ctx context.Context, cred *user.Credentials) error {
	const q = `
		INSERT INTO credentials (user_id, password_hash, salt, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, q,
		cred.UserID.String(), cred.PasswordHash, cred.Salt, cred.UpdatedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("create credentials: %w", err)
	}
	return nil
}

func (r *credentialsRepo) Update(ctx context.Context, cred *user.Credentials) error {
	const q = `
		UPDATE credentials
		SET password_hash = $1, salt = $2, updated_at = $3
		WHERE user_id = $4
	`
	_, err := r.db.ExecContext(ctx, q,
		cred.PasswordHash, cred.Salt, time.Now().UTC(), cred.UserID.String(),
	)
	if err != nil {
		return fmt.Errorf("update credentials: %w", err)
	}
	return nil
}

func (r *credentialsRepo) GetByUserID(ctx context.Context, userID ulid.ULID) (*user.Credentials, error) {
	const q = `SELECT user_id, password_hash, salt, updated_at FROM credentials WHERE user_id=$1`

	var cred user.Credentials
	var userIDStr string
	var updatedAt time.Time

	err := r.db.QueryRowContext(ctx, q, userID.String()).Scan(
		&userIDStr, &cred.PasswordHash, &cred.Salt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("get credentials: %w", err)
	}

	id, err := ulid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	cred.UserID = id
	cred.UpdatedAt = updatedAt.UTC()

	return &cred, nil
}
