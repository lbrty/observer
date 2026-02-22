package repository

import (
	"context"

	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/reference"
	"github.com/lbrty/observer/internal/domain/user"
)

//go:generate mockgen -destination=mock/repository.go -package=mock github.com/lbrty/observer/internal/repository UserRepository,CredentialsRepository,MFARepository,VerificationTokenRepository,SessionRepository,PermissionLoader,PermissionRepository,CountryRepository,StateRepository,PlaceRepository,OfficeRepository,CategoryRepository

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Create(ctx context.Context, u *user.User) error
	GetByID(ctx context.Context, id ulid.ULID) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	GetByPhone(ctx context.Context, phone string) (*user.User, error)
	Update(ctx context.Context, u *user.User) error
	UpdateVerified(ctx context.Context, id ulid.ULID, verified bool) error
	List(ctx context.Context, filter user.UserListFilter) ([]*user.User, int, error)
}

// CredentialsRepository defines persistence operations for credentials.
type CredentialsRepository interface {
	Create(ctx context.Context, cred *user.Credentials) error
	GetByUserID(ctx context.Context, userID ulid.ULID) (*user.Credentials, error)
}

// MFARepository defines persistence operations for MFA configs.
type MFARepository interface {
	Create(ctx context.Context, cfg *user.MFAConfig) error
	GetByUserID(ctx context.Context, userID ulid.ULID) (*user.MFAConfig, error)
}

// VerificationTokenRepository defines persistence operations for verification tokens.
type VerificationTokenRepository interface {
	Create(ctx context.Context, token *user.VerificationToken) error
	GetByToken(ctx context.Context, token string) (*user.VerificationToken, error)
	MarkUsed(ctx context.Context, id ulid.ULID) error
}

// SessionRepository defines persistence operations for sessions.
type SessionRepository interface {
	Create(ctx context.Context, session *auth.Session) error
	GetByRefreshToken(ctx context.Context, token string) (*auth.Session, error)
	Delete(ctx context.Context, id ulid.ULID) error
	DeleteByRefreshToken(ctx context.Context, token string) error
}

// PermissionLoader loads project-level permissions for authorization.
type PermissionLoader interface {
	GetPermission(ctx context.Context, userID ulid.ULID, projectID string) (*project.Permission, error)
	IsProjectOwner(ctx context.Context, userID ulid.ULID, projectID string) (bool, error)
}

// PermissionRepository defines CRUD operations for project permissions (admin use).
type PermissionRepository interface {
	List(ctx context.Context, projectID string) ([]*project.ProjectPermission, error)
	GetByID(ctx context.Context, id string) (*project.ProjectPermission, error)
	Create(ctx context.Context, p *project.ProjectPermission) error
	Update(ctx context.Context, p *project.ProjectPermission) error
	Delete(ctx context.Context, id string) error
}

// CountryRepository defines persistence operations for countries.
type CountryRepository interface {
	List(ctx context.Context) ([]*reference.Country, error)
	GetByID(ctx context.Context, id string) (*reference.Country, error)
	Create(ctx context.Context, c *reference.Country) error
	Update(ctx context.Context, c *reference.Country) error
	Delete(ctx context.Context, id string) error
}

// StateRepository defines persistence operations for states.
type StateRepository interface {
	List(ctx context.Context, countryID string) ([]*reference.State, error)
	GetByID(ctx context.Context, id string) (*reference.State, error)
	Create(ctx context.Context, s *reference.State) error
	Update(ctx context.Context, s *reference.State) error
	Delete(ctx context.Context, id string) error
}

// PlaceRepository defines persistence operations for places.
type PlaceRepository interface {
	List(ctx context.Context, stateID string) ([]*reference.Place, error)
	GetByID(ctx context.Context, id string) (*reference.Place, error)
	Create(ctx context.Context, p *reference.Place) error
	Update(ctx context.Context, p *reference.Place) error
	Delete(ctx context.Context, id string) error
}

// OfficeRepository defines persistence operations for offices.
type OfficeRepository interface {
	List(ctx context.Context) ([]*reference.Office, error)
	GetByID(ctx context.Context, id string) (*reference.Office, error)
	Create(ctx context.Context, o *reference.Office) error
	Update(ctx context.Context, o *reference.Office) error
	Delete(ctx context.Context, id string) error
}

// CategoryRepository defines persistence operations for categories.
type CategoryRepository interface {
	List(ctx context.Context) ([]*reference.Category, error)
	GetByID(ctx context.Context, id string) (*reference.Category, error)
	Create(ctx context.Context, c *reference.Category) error
	Update(ctx context.Context, c *reference.Category) error
	Delete(ctx context.Context, id string) error
}
