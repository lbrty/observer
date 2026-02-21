package auth

import (
	"context"

	"github.com/oklog/ulid/v2"
)

//go:generate mockgen -destination=mock/repository.go -package=mock github.com/lbrty/observer/internal/domain/auth SessionRepository

// SessionRepository defines persistence operations for sessions.
type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	GetByRefreshToken(ctx context.Context, token string) (*Session, error)
	Delete(ctx context.Context, id ulid.ULID) error
	DeleteByRefreshToken(ctx context.Context, token string) error
}
