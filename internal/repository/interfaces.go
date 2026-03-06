package repository

import (
	"context"

	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/document"
	"github.com/lbrty/observer/internal/domain/household"
	"github.com/lbrty/observer/internal/domain/migration"
	"github.com/lbrty/observer/internal/domain/note"
	"github.com/lbrty/observer/internal/domain/person"
	"github.com/lbrty/observer/internal/domain/pet"
	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/reference"
	"github.com/lbrty/observer/internal/domain/report"
	"github.com/lbrty/observer/internal/domain/support"
	"github.com/lbrty/observer/internal/domain/tag"
	"github.com/lbrty/observer/internal/domain/user"
)

//go:generate mockgen -destination=mock/repository.go -package=mock github.com/lbrty/observer/internal/repository UserRepository,CredentialsRepository,MFARepository,VerificationTokenRepository,SessionRepository,PermissionLoader,PermissionRepository,ProjectRepository,CountryRepository,StateRepository,PlaceRepository,OfficeRepository,CategoryRepository,TagRepository,PersonRepository,PersonCategoryRepository,PersonTagRepository,SupportRecordRepository,MigrationRecordRepository,HouseholdRepository,HouseholdMemberRepository,PersonNoteRepository,DocumentRepository,PetRepository,ReportRepository

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Create(ctx context.Context, u *user.User) error
	GetByID(ctx context.Context, id ulid.ULID) (*user.User, error)
	GetByIDs(ctx context.Context, ids []ulid.ULID) ([]*user.User, error)
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
	Update(ctx context.Context, cred *user.Credentials) error
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
	ListByUserID(ctx context.Context, userID string) ([]*project.ProjectPermission, error)
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
	ListAll(ctx context.Context) ([]*reference.State, error)
	List(ctx context.Context, countryID string) ([]*reference.State, error)
	GetByID(ctx context.Context, id string) (*reference.State, error)
	Create(ctx context.Context, s *reference.State) error
	Update(ctx context.Context, s *reference.State) error
	Delete(ctx context.Context, id string) error
}

// PlaceRepository defines persistence operations for places.
type PlaceRepository interface {
	ListAll(ctx context.Context) ([]*reference.Place, error)
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

// ProjectRepository defines persistence operations for projects.
type ProjectRepository interface {
	List(ctx context.Context, filter project.ProjectListFilter) ([]*project.Project, int, error)
	GetByID(ctx context.Context, id string) (*project.Project, error)
	Create(ctx context.Context, p *project.Project) error
	Update(ctx context.Context, p *project.Project) error
}

// TagRepository defines persistence operations for tags.
type TagRepository interface {
	List(ctx context.Context, projectID string) ([]*tag.Tag, error)
	GetByID(ctx context.Context, id string) (*tag.Tag, error)
	Create(ctx context.Context, t *tag.Tag) error
	Update(ctx context.Context, t *tag.Tag) error
	Delete(ctx context.Context, id string) error
}

// PersonRepository defines persistence operations for people.
type PersonRepository interface {
	List(ctx context.Context, filter person.PersonListFilter) ([]*person.Person, int, error)
	GetByID(ctx context.Context, id string) (*person.Person, error)
	Create(ctx context.Context, p *person.Person) error
	Update(ctx context.Context, p *person.Person) error
	Delete(ctx context.Context, id string) error
}

// PersonCategoryRepository manages person-category associations.
type PersonCategoryRepository interface {
	List(ctx context.Context, personID string) ([]string, error)
	ReplaceAll(ctx context.Context, personID string, categoryIDs []string) error
}

// PersonTagRepository manages person-tag associations.
type PersonTagRepository interface {
	List(ctx context.Context, personID string) ([]string, error)
	ReplaceAll(ctx context.Context, personID string, tagIDs []string) error
}

// SupportRecordRepository defines persistence operations for support records.
type SupportRecordRepository interface {
	List(ctx context.Context, filter support.RecordListFilter) ([]*support.Record, int, error)
	GetByID(ctx context.Context, id string) (*support.Record, error)
	Create(ctx context.Context, r *support.Record) error
	Update(ctx context.Context, r *support.Record) error
	Delete(ctx context.Context, id string) error
}

// MigrationRecordRepository defines persistence operations for migration records.
type MigrationRecordRepository interface {
	ListByPerson(ctx context.Context, personID string) ([]*migration.Record, error)
	GetByID(ctx context.Context, id string) (*migration.Record, error)
	Create(ctx context.Context, r *migration.Record) error
	Update(ctx context.Context, r *migration.Record) error
}

// HouseholdRepository defines persistence operations for households.
type HouseholdRepository interface {
	List(ctx context.Context, projectID string, page, perPage int) ([]*household.Household, int, error)
	GetByID(ctx context.Context, id string) (*household.Household, error)
	Create(ctx context.Context, h *household.Household) error
	Update(ctx context.Context, h *household.Household) error
	Delete(ctx context.Context, id string) error
}

// HouseholdMemberRepository manages household membership.
type HouseholdMemberRepository interface {
	List(ctx context.Context, householdID string) ([]*household.Member, error)
	Add(ctx context.Context, m *household.Member) error
	Remove(ctx context.Context, householdID, personID string) error
}

// PersonNoteRepository defines persistence operations for person notes.
type PersonNoteRepository interface {
	List(ctx context.Context, personID string) ([]*note.Note, error)
	GetByID(ctx context.Context, id string) (*note.Note, error)
	Create(ctx context.Context, n *note.Note) error
	Update(ctx context.Context, n *note.Note) error
	Delete(ctx context.Context, id string) error
}

// DocumentRepository defines persistence operations for document metadata.
type DocumentRepository interface {
	List(ctx context.Context, personID string) ([]*document.Document, error)
	GetByID(ctx context.Context, id string) (*document.Document, error)
	Create(ctx context.Context, d *document.Document) error
	Update(ctx context.Context, d *document.Document) error
	Delete(ctx context.Context, id string) error
}

// PetRepository defines persistence operations for pets.
type PetRepository interface {
	List(ctx context.Context, projectID string, page, perPage int) ([]*pet.Pet, int, error)
	GetByID(ctx context.Context, id string) (*pet.Pet, error)
	Create(ctx context.Context, p *pet.Pet) error
	Update(ctx context.Context, p *pet.Pet) error
	Delete(ctx context.Context, id string) error
}

// ReportRepository provides aggregation queries for ADR-005 reports.
type ReportRepository interface {
	CountConsultations(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	CountBySex(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	CountByIDPStatus(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	CountByCategory(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	CountByCurrentRegion(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	CountBySphere(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	CountByOffice(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	CountByAgeGroup(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	CountByTag(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	CountFamilyUnits(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	CountByCaseStatus(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error)
	StatusFlowReport(ctx context.Context, f report.ReportFilter) ([]report.StatusFlow, error)
}
