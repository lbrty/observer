package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/domain/auth"
)

// SessionRepository is a PostgreSQL-backed session repository.
type SessionRepository struct {
	db *sqlx.DB
}

// NewSessionRepository creates a SessionRepository backed by the given DB.
func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, s *auth.Session) error {
	const q = `
		INSERT INTO sessions (id, user_id, refresh_token, user_agent, ip, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, q,
		s.ID.String(), s.UserID.String(), s.RefreshToken,
		s.UserAgent, s.IP, s.ExpiresAt.UTC(), s.CreatedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (r *SessionRepository) GetByRefreshToken(ctx context.Context, token string) (*auth.Session, error) {
	const q = `
		SELECT id, user_id, refresh_token, user_agent, ip, expires_at, created_at
		FROM sessions WHERE refresh_token=$1
	`
	return r.scanSession(r.db.QueryRowContext(ctx, q, token))
}

func (r *SessionRepository) Delete(ctx context.Context, id ulid.ULID) error {
	const q = `DELETE FROM sessions WHERE id=$1`
	_, err := r.db.ExecContext(ctx, q, id.String())
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (r *SessionRepository) DeleteByRefreshToken(ctx context.Context, token string) error {
	const q = `DELETE FROM sessions WHERE refresh_token=$1`
	_, err := r.db.ExecContext(ctx, q, token)
	if err != nil {
		return fmt.Errorf("delete session by token: %w", err)
	}
	return nil
}

func (r *SessionRepository) scanSession(row *sql.Row) (*auth.Session, error) {
	var s auth.Session
	var idStr, userIDStr string
	var expiresAt, createdAt time.Time

	err := row.Scan(
		&idStr, &userIDStr, &s.RefreshToken,
		&s.UserAgent, &s.IP, &expiresAt, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, auth.ErrSessionNotFound
		}
		return nil, fmt.Errorf("scan session: %w", err)
	}

	id, err := ulid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("parse session id: %w", err)
	}

	userID, err := ulid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	s.ID = id
	s.UserID = userID
	s.ExpiresAt = expiresAt.UTC()
	s.CreatedAt = createdAt.UTC()

	return &s, nil
}
