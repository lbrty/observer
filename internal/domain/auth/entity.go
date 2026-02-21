package auth

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

// Session represents an authenticated user session.
type Session struct {
	ID           ulid.ULID
	UserID       ulid.ULID
	RefreshToken string
	UserAgent    string
	IP           string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
