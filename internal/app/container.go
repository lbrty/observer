package app

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/database"
	"github.com/lbrty/observer/internal/repository"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
	ucauth "github.com/lbrty/observer/internal/usecase/auth"
	ucmy "github.com/lbrty/observer/internal/usecase/my"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
	ucreport "github.com/lbrty/observer/internal/usecase/report"
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

	// Report Use Cases
	ReportUC   *ucreport.ReportUseCase
	ReportRepo repository.ReportRepository

	// Login attempt tracking
	LoginAttemptStore repository.LoginAttemptStore
}

// NewContainer wires all dependencies from config, database, and redis.
func NewContainer(cfg *config.Config, db database.DB, redisClient *redis.Client) (*Container, error) {
	rsaKeys, err := crypto.LoadRSAKeys(cfg.JWT.PrivateKeyPath, cfg.JWT.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load RSA keys: %w", err)
	}

	sqlxDB := db.GetDB()

	userRepo := repository.NewUserRepository(sqlxDB)
	credRepo := repository.NewCredentialsRepository(sqlxDB)
	sessionRepo := repository.NewSessionRepository(sqlxDB)
	mfaRepo := repository.NewMFARepository(sqlxDB)
	permRepo := repository.NewPermissionRepository(sqlxDB)
	permCRUDRepo := repository.NewProjectPermissionRepository(sqlxDB)
	countryRepo := repository.NewCountryRepository(sqlxDB)
	stateRepo := repository.NewStateRepository(sqlxDB)
	placeRepo := repository.NewPlaceRepository(sqlxDB)
	officeRepo := repository.NewOfficeRepository(sqlxDB)
	categoryRepo := repository.NewCategoryRepository(sqlxDB)
	projectRepo := repository.NewProjectRepository(sqlxDB)
	tagRepo := repository.NewTagRepository(sqlxDB)
	personRepo := repository.NewPersonRepository(sqlxDB)
	personCatRepo := repository.NewPersonCategoryRepository(sqlxDB)
	personTagRepo := repository.NewPersonTagRepository(sqlxDB)
	supportRepo := repository.NewSupportRecordRepository(sqlxDB)
	migrationRepo := repository.NewMigrationRecordRepository(sqlxDB)
	householdRepo := repository.NewHouseholdRepository(sqlxDB)
	householdMemberRepo := repository.NewHouseholdMemberRepository(sqlxDB)
	noteRepo := repository.NewPersonNoteRepository(sqlxDB)
	documentRepo := repository.NewDocumentRepository(sqlxDB)
	petRepo := repository.NewPetRepository(sqlxDB)
	reportRepo := repository.NewReportRepository(sqlxDB)

	hasher := crypto.NewArgonHasher()
	tokenGen := crypto.NewRSATokenGenerator(
		rsaKeys,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
		cfg.JWT.MFATempTTL,
		cfg.JWT.Issuer,
	)

	authUC := ucauth.NewAuthUseCase(userRepo, credRepo, sessionRepo, mfaRepo, hasher, tokenGen)
	userUC := ucadmin.NewUserUseCase(userRepo, credRepo, hasher)
	permUC := ucadmin.NewPermissionUseCase(permCRUDRepo, userRepo)

	countryUC := ucadmin.NewCountryUseCase(countryRepo)
	stateUC := ucadmin.NewStateUseCase(stateRepo)
	placeUC := ucadmin.NewPlaceUseCase(placeRepo)
	officeUC := ucadmin.NewOfficeUseCase(officeRepo)
	categoryUC := ucadmin.NewCategoryUseCase(categoryRepo)

	myProjectsUC := ucmy.NewMyProjectsUseCase(permCRUDRepo, projectRepo)

	projectUC := ucadmin.NewProjectUseCase(projectRepo)
	tagUC := ucproject.NewTagUseCase(tagRepo)
	personUC := ucproject.NewPersonUseCase(personRepo)
	personCategoryUC := ucproject.NewPersonCategoryUseCase(personCatRepo)
	personTagUC := ucproject.NewPersonTagUseCase(personTagRepo)
	supportRecordUC := ucproject.NewSupportRecordUseCase(supportRepo)
	migrationRecordUC := ucproject.NewMigrationRecordUseCase(migrationRepo)
	householdUC := ucproject.NewHouseholdUseCase(householdRepo, householdMemberRepo)
	noteUC := ucproject.NewNoteUseCase(noteRepo)
	documentUC := ucproject.NewDocumentUseCase(documentRepo)
	petUC := ucproject.NewPetUseCase(petRepo)
	reportUC := ucreport.NewReportUseCase(reportRepo)
	loginAttemptStore := repository.NewLoginAttemptStore(redisClient)

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
		DocumentUC:        documentUC,
		PetUC:             petUC,
		ReportUC:          reportUC,
		ReportRepo:        reportRepo,
		LoginAttemptStore: loginAttemptStore,
	}, nil
}
