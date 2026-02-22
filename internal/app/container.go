package app

import (
	"fmt"

	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/database"
	"github.com/lbrty/observer/internal/repository"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
	ucauth "github.com/lbrty/observer/internal/usecase/auth"
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

	// Auth Use Cases
	RegisterUC     *ucauth.RegisterUseCase
	LoginUC        *ucauth.LoginUseCase
	RefreshTokenUC *ucauth.RefreshTokenUseCase
	LogoutUC       *ucauth.LogoutUseCase

	// Admin Use Cases
	ListUsersUC  *ucadmin.ListUsersUseCase
	GetUserUC    *ucadmin.GetUserUseCase
	UpdateUserUC *ucadmin.UpdateUserUseCase

	// Permission Use Cases
	ListPermsUC  *ucadmin.ListPermissionsUseCase
	AssignPermUC *ucadmin.AssignPermissionUseCase
	UpdatePermUC *ucadmin.UpdatePermissionUseCase
	RevokePermUC *ucadmin.RevokePermissionUseCase

	// Reference Use Cases
	CountryUC  *ucadmin.CountryUseCase
	StateUC    *ucadmin.StateUseCase
	PlaceUC    *ucadmin.PlaceUseCase
	OfficeUC   *ucadmin.OfficeUseCase
	CategoryUC *ucadmin.CategoryUseCase
}

// NewContainer wires all dependencies from config and database.
func NewContainer(cfg *config.Config, db database.DB) (*Container, error) {
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

	hasher := crypto.NewArgonHasher()
	tokenGen := crypto.NewRSATokenGenerator(
		rsaKeys,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
		cfg.JWT.MFATempTTL,
		cfg.JWT.Issuer,
	)

	registerUC := ucauth.NewRegisterUseCase(userRepo, credRepo, hasher)
	loginUC := ucauth.NewLoginUseCase(userRepo, credRepo, sessionRepo, mfaRepo, hasher, tokenGen)
	refreshUC := ucauth.NewRefreshTokenUseCase(sessionRepo, tokenGen)
	logoutUC := ucauth.NewLogoutUseCase(sessionRepo)

	listUsersUC := ucadmin.NewListUsersUseCase(userRepo)
	getUserUC := ucadmin.NewGetUserUseCase(userRepo)
	updateUserUC := ucadmin.NewUpdateUserUseCase(userRepo)

	listPermsUC := ucadmin.NewListPermissionsUseCase(permCRUDRepo)
	assignPermUC := ucadmin.NewAssignPermissionUseCase(permCRUDRepo)
	updatePermUC := ucadmin.NewUpdatePermissionUseCase(permCRUDRepo)
	revokePermUC := ucadmin.NewRevokePermissionUseCase(permCRUDRepo)

	countryUC := ucadmin.NewCountryUseCase(countryRepo)
	stateUC := ucadmin.NewStateUseCase(stateRepo)
	placeUC := ucadmin.NewPlaceUseCase(placeRepo)
	officeUC := ucadmin.NewOfficeUseCase(officeRepo)
	categoryUC := ucadmin.NewCategoryUseCase(categoryRepo)

	return &Container{
		UserRepo:       userRepo,
		CredRepo:       credRepo,
		SessionRepo:    sessionRepo,
		PermissionRepo: permRepo,
		PermCRUDRepo:   permCRUDRepo,
		CountryRepo:    countryRepo,
		StateRepo:      stateRepo,
		PlaceRepo:      placeRepo,
		OfficeRepo:     officeRepo,
		CategoryRepo:   categoryRepo,
		PasswordHasher: hasher,
		TokenGenerator: tokenGen,
		RegisterUC:     registerUC,
		LoginUC:        loginUC,
		RefreshTokenUC: refreshUC,
		LogoutUC:       logoutUC,
		ListUsersUC:    listUsersUC,
		GetUserUC:      getUserUC,
		UpdateUserUC:   updateUserUC,
		ListPermsUC:    listPermsUC,
		AssignPermUC:   assignPermUC,
		UpdatePermUC:   updatePermUC,
		RevokePermUC:   revokePermUC,
		CountryUC:      countryUC,
		StateUC:        stateUC,
		PlaceUC:        placeUC,
		OfficeUC:       officeUC,
		CategoryUC:     categoryUC,
	}, nil
}
