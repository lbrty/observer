package user

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// Role represents a user role.
type Role string

const (
	RoleAdmin      Role = "admin"
	RoleStaff      Role = "staff"
	RoleConsultant Role = "consultant"
	RoleGuest      Role = "guest"
)

// User is the core user domain entity.
type User struct {
	ID         ulid.ULID
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	OfficeID   *string
	Role       Role
	IsVerified bool
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// UserListFilter controls pagination and filtering for user listing.
type UserListFilter struct {
	Page     int
	PerPage  int
	Search   string
	Role     string
	IsActive *bool
}

// CanLogin returns an error if the user is not allowed to log in.
func (u *User) CanLogin() error {
	if !u.IsActive {
		return ErrUserNotActive
	}
	return nil
}

// Credentials stores hashed password and salt for a user.
type Credentials struct {
	UserID       ulid.ULID
	PasswordHash string
	Salt         string
	UpdatedAt    time.Time
}

// MFAConfig stores MFA settings for a user.
type MFAConfig struct {
	UserID    ulid.ULID
	Method    string
	Secret    string
	Phone     string
	IsEnabled bool
	CreatedAt time.Time
}

// VerificationToken represents an email/password-reset verification token.
type VerificationToken struct {
	ID        ulid.ULID
	UserID    ulid.ULID
	Token     string
	Type      string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

// ValidateRole parses and validates a role string.
func ValidateRole(role string) (Role, error) {
	switch Role(role) {
	case RoleAdmin, RoleStaff, RoleConsultant, RoleGuest:
		return Role(role), nil
	default:
		return "", ErrInvalidRole
	}
}
