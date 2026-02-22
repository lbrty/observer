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

// mfaRepo is a PostgreSQL-backed MFA config repository.
type mfaRepo struct {
	db *sqlx.DB
}

// NewMFARepository creates an MFARepository backed by the given DB.
func NewMFARepository(db *sqlx.DB) MFARepository {
	return &mfaRepo{db: db}
}

func (r *mfaRepo) Create(ctx context.Context, cfg *user.MFAConfig) error {
	const q = `
		INSERT INTO mfa_configs (user_id, method, secret, phone, is_enabled, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, q,
		cfg.UserID.String(), cfg.Method, cfg.Secret, cfg.Phone, cfg.IsEnabled, cfg.CreatedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("create mfa config: %w", err)
	}
	return nil
}

func (r *mfaRepo) GetByUserID(ctx context.Context, userID ulid.ULID) (*user.MFAConfig, error) {
	const q = `SELECT user_id, method, secret, phone, is_enabled, created_at FROM mfa_configs WHERE user_id=$1`

	var cfg user.MFAConfig
	var userIDStr string
	var createdAt time.Time

	err := r.db.QueryRowContext(ctx, q, userID.String()).Scan(
		&userIDStr, &cfg.Method, &cfg.Secret, &cfg.Phone, &cfg.IsEnabled, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("get mfa config: %w", err)
	}

	id, err := ulid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	cfg.UserID = id
	cfg.CreatedAt = createdAt.UTC()

	return &cfg, nil
}
