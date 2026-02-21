package user

import (
	"context"

	"github.com/oklog/ulid/v2"
)

//go:generate mockgen -destination=mock/repository.go -package=mock github.com/lbrty/observer/internal/domain/user UserRepository,CredentialsRepository,MFARepository,VerificationTokenRepository

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id ulid.ULID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
	Update(ctx context.Context, u *User) error
	UpdateVerified(ctx context.Context, id ulid.ULID, verified bool) error
}

// CredentialsRepository defines persistence operations for credentials.
type CredentialsRepository interface {
	Create(ctx context.Context, cred *Credentials) error
	GetByUserID(ctx context.Context, userID ulid.ULID) (*Credentials, error)
}

// MFARepository defines persistence operations for MFA configs.
type MFARepository interface {
	Create(ctx context.Context, cfg *MFAConfig) error
	GetByUserID(ctx context.Context, userID ulid.ULID) (*MFAConfig, error)
}

// VerificationTokenRepository defines persistence operations for verification tokens.
type VerificationTokenRepository interface {
	Create(ctx context.Context, token *VerificationToken) error
	GetByToken(ctx context.Context, token string) (*VerificationToken, error)
	MarkUsed(ctx context.Context, id ulid.ULID) error
}
