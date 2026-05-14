package app

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/database"
	"github.com/lbrty/observer/internal/repository"
	repoaudit "github.com/lbrty/observer/internal/repository/audit"
	repoauth "github.com/lbrty/observer/internal/repository/auth"
	repohousehold "github.com/lbrty/observer/internal/repository/household"
	repomigration "github.com/lbrty/observer/internal/repository/migration"
	repoperson "github.com/lbrty/observer/internal/repository/person"
	repopet "github.com/lbrty/observer/internal/repository/pet"
	repoproj "github.com/lbrty/observer/internal/repository/project"
	reporeference "github.com/lbrty/observer/internal/repository/reference"
	reporeport "github.com/lbrty/observer/internal/repository/report"
	reposearch "github.com/lbrty/observer/internal/repository/search"
	reposupport "github.com/lbrty/observer/internal/repository/support"
	repotag "github.com/lbrty/observer/internal/repository/tag"
	repouser "github.com/lbrty/observer/internal/repository/user"
	"github.com/lbrty/observer/internal/storage"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
	ucauth "github.com/lbrty/observer/internal/usecase/auth"
	ucmy "github.com/lbrty/observer/internal/usecase/my"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
	ucreport "github.com/lbrty/observer/internal/usecase/report"
	ucsearch "github.com/lbrty/observer/internal/usecase/search"
)

// Container holds all application dependencies.
type Container struct {
	// Repositories
	UserRepo       repository.UserRepository
	CredRepo       repository.CredentialsRepository
	SessionRepo    repository.SessionRepository
	PermissionRepo repository.PermissionLoader
	PermCRUDRepo   repository.PermissionRepository
	CountryRepo    repository.CountryRepository
	StateRepo      repository.StateRepository
	PlaceRepo      repository.PlaceRepository
	OfficeRepo     repository.OfficeRepository
	CategoryRepo   repository.CategoryRepository

	// Services
	PasswordHasher crypto.PasswordHasher
	TokenGenerator crypto.TokenGenerator

	// Auth Use Case
	AuthUC *ucauth.AuthUseCase

	// Admin Use Cases
	UserUC *ucadmin.UserUseCase
	PermUC *ucadmin.PermissionUseCase

	// Reference Use Cases
	CountryUC  *ucadmin.CountryUseCase
	StateUC    *ucadmin.StateUseCase
	PlaceUC    *ucadmin.PlaceUseCase
	OfficeUC   *ucadmin.OfficeUseCase
	CategoryUC *ucadmin.CategoryUseCase

	// My Use Cases
	MyProjectsUC *ucmy.MyProjectsUseCase

	// Project Use Cases
	ProjectUC         *ucadmin.ProjectUseCase
	TagUC             *ucproject.TagUseCase
	PersonUC          *ucproject.PersonUseCase
	PersonCategoryUC  *ucproject.PersonCategoryUseCase
	PersonTagUC       *ucproject.PersonTagUseCase
	SupportRecordUC   *ucproject.SupportRecordUseCase
	MigrationRecordUC *ucproject.MigrationRecordUseCase
	HouseholdUC       *ucproject.HouseholdUseCase
	NoteUC            *ucproject.NoteUseCase
	DocumentUC        *ucproject.DocumentUseCase
	PetUC             *ucproject.PetUseCase
	PetTagUC          *ucproject.PetTagUseCase

	// Audit Use Case
	AuditUC *ucaudit.AuditUseCase

	// Report Use Cases
	ReportUC    *ucreport.ReportUseCase
	PetReportUC *ucreport.PetReportUseCase
	ReportRepo  repository.ReportRepository

	// Storage
	FileStorage storage.FileStorage

	// Search
	SearchUC *ucsearch.SearchUseCase
}

// NewContainer wires all dependencies from config, database, and redis.
func NewContainer(cfg *config.Config, db database.DB, redisClient *redis.Client) (*Container, error) {
	rsaKeys, err := crypto.LoadRSAKeys(cfg.JWT.PrivateKeyPath, cfg.JWT.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load RSA keys: %w", err)
	}

	sqlxDB := db.GetDB()

	userRepo := repouser.New(sqlxDB)
	credRepo := repouser.NewCredentials(sqlxDB)
	sessionRepo := repoauth.NewSession(sqlxDB)
	mfaRepo := repouser.NewMFA(sqlxDB)
	mfaRecoveryRepo := repouser.NewMFARecovery(sqlxDB)
	permRepo := repoproj.NewLoader(sqlxDB)
	permCRUDRepo := repoproj.NewCRUD(sqlxDB)
	countryRepo := reporeference.NewCountry(sqlxDB)
	stateRepo := reporeference.NewState(sqlxDB)
	placeRepo := reporeference.NewPlace(sqlxDB)
	officeRepo := reporeference.NewOffice(sqlxDB)
	categoryRepo := reporeference.NewCategory(sqlxDB)
	projectRepo := repoproj.New(sqlxDB)
	tagRepo := repotag.New(sqlxDB)
	personRepo := repoperson.New(sqlxDB)
	personCatRepo := repoperson.NewCategory(sqlxDB)
	personTagRepo := repoperson.NewTag(sqlxDB)
	supportRepo := reposupport.New(sqlxDB)
	migrationRepo := repomigration.New(sqlxDB)
	householdRepo := repohousehold.New(sqlxDB)
	householdMemberRepo := repohousehold.NewMember(sqlxDB)
	noteRepo := repoperson.NewNote(sqlxDB)
	documentRepo := repoperson.NewDocument(sqlxDB)
	petRepo := repopet.New(sqlxDB)
	petTagRepo := repopet.NewTag(sqlxDB)
	searchRepo := reposearch.New(sqlxDB)
	auditRepo := repoaudit.New(sqlxDB)
	reportRepo := reporeport.New(sqlxDB)
	petReportRepo := repopet.NewReport(sqlxDB)

	var fileStorage storage.FileStorage
	switch cfg.Storage.Backend {
	case "s3":
		fileStorage, err = storage.NewS3Storage(cfg.Storage)
	default:
		fileStorage, err = storage.NewLocalStorage(cfg.Storage.Path)
	}
	if err != nil {
		return nil, fmt.Errorf("init file storage: %w", err)
	}

	hasher := crypto.NewArgonHasher()
	tokenGen := crypto.NewRSATokenGenerator(
		rsaKeys,
		cfg.JWT.AccessTTL,
		cfg.JWT.MFATempTTL,
		cfg.JWT.Issuer,
	)

	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	loginAttemptStore := repoauth.NewLoginAttemptStore(redisClient, userRepo)
	authUC := ucauth.NewAuthUseCase(
		userRepo, credRepo, sessionRepo, mfaRepo, mfaRecoveryRepo, hasher, tokenGen, loginAttemptStore,
	)
	userUC := ucadmin.NewUserUseCase(userRepo, credRepo, hasher, sessionRepo, loginAttemptStore, auditUC)
	permUC := ucadmin.NewPermissionUseCase(permCRUDRepo, userRepo, auditUC)
	countryUC := ucadmin.NewCountryUseCase(countryRepo)
	stateUC := ucadmin.NewStateUseCase(stateRepo)
	placeUC := ucadmin.NewPlaceUseCase(placeRepo)
	officeUC := ucadmin.NewOfficeUseCase(officeRepo)
	categoryUC := ucadmin.NewCategoryUseCase(categoryRepo)
	myProjectsUC := ucmy.NewMyProjectsUseCase(permCRUDRepo, projectRepo)
	projectUC := ucadmin.NewProjectUseCase(projectRepo, permCRUDRepo, auditUC)
	tagUC := ucproject.NewTagUseCase(tagRepo, auditUC)
	personUC := ucproject.NewPersonUseCase(personRepo, personTagRepo, auditUC)
	personCategoryUC := ucproject.NewPersonCategoryUseCase(personCatRepo, personRepo)
	personTagUC := ucproject.NewPersonTagUseCase(personTagRepo, personRepo)
	supportRecordUC := ucproject.NewSupportRecordUseCase(supportRepo, personRepo, auditUC)
	migrationRecordUC := ucproject.NewMigrationRecordUseCase(migrationRepo, personRepo, auditUC)
	householdUC := ucproject.NewHouseholdUseCase(householdRepo, householdMemberRepo, personRepo, auditUC)
	noteUC := ucproject.NewNoteUseCase(noteRepo, personRepo, auditUC)
	documentUC := ucproject.NewDocumentUseCase(documentRepo, personRepo, fileStorage, auditUC)
	petUC := ucproject.NewPetUseCase(petRepo, petTagRepo, auditUC)
	petTagUC := ucproject.NewPetTagUseCase(petTagRepo, petRepo)
	reportUC := ucreport.NewReportUseCase(reportRepo)
	petReportUC := ucreport.NewPetReportUseCase(petReportRepo)
	searchUC := ucsearch.NewSearchUseCase(searchRepo)

	return &Container{
		UserRepo:          userRepo,
		CredRepo:          credRepo,
		SessionRepo:       sessionRepo,
		PermissionRepo:    permRepo,
		PermCRUDRepo:      permCRUDRepo,
		CountryRepo:       countryRepo,
		StateRepo:         stateRepo,
		PlaceRepo:         placeRepo,
		OfficeRepo:        officeRepo,
		CategoryRepo:      categoryRepo,
		PasswordHasher:    hasher,
		TokenGenerator:    tokenGen,
		AuthUC:            authUC,
		UserUC:            userUC,
		PermUC:            permUC,
		CountryUC:         countryUC,
		StateUC:           stateUC,
		PlaceUC:           placeUC,
		OfficeUC:          officeUC,
		CategoryUC:        categoryUC,
		MyProjectsUC:      myProjectsUC,
		ProjectUC:         projectUC,
		TagUC:             tagUC,
		PersonUC:          personUC,
		PersonCategoryUC:  personCategoryUC,
		PersonTagUC:       personTagUC,
		SupportRecordUC:   supportRecordUC,
		MigrationRecordUC: migrationRecordUC,
		HouseholdUC:       householdUC,
		NoteUC:            noteUC,
		FileStorage:       fileStorage,
		DocumentUC:        documentUC,
		PetUC:             petUC,
		PetTagUC:          petTagUC,
		AuditUC:           auditUC,
		ReportUC:          reportUC,
		PetReportUC:       petReportUC,
		ReportRepo:        reportRepo,
		SearchUC:          searchUC,
	}, nil
}
