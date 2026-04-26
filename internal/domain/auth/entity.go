package auth

import (
	"time"

	"github.com/oklog/ulid/v2"
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
