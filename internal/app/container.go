package app

import (
	"fmt"

	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/database"
	domainauth "github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/postgres"
	ucauth "github.com/lbrty/observer/internal/usecase/auth"
)

// Container holds all application dependencies.
type Container struct {
	// Repositories
	UserRepo       user.UserRepository
	CredRepo       user.CredentialsRepository
	SessionRepo    domainauth.SessionRepository
	PermissionRepo project.PermissionLoader

	// Services
	PasswordHasher crypto.PasswordHasher
	TokenGenerator crypto.TokenGenerator

	// Use Cases
	RegisterUC     *ucauth.RegisterUseCase
	LoginUC        *ucauth.LoginUseCase
	RefreshTokenUC *ucauth.RefreshTokenUseCase
	LogoutUC       *ucauth.LogoutUseCase
}

// NewContainer wires all dependencies from config and database.
func NewContainer(cfg *config.Config, db database.DB) (*Container, error) {
	rsaKeys, err := crypto.LoadRSAKeys(cfg.JWT.PrivateKeyPath, cfg.JWT.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load RSA keys: %w", err)
	}

	sqlxDB := db.GetDB()

	userRepo := postgres.NewUserRepository(sqlxDB)
	credRepo := postgres.NewCredentialsRepository(sqlxDB)
	sessionRepo := postgres.NewSessionRepository(sqlxDB)
	mfaRepo := postgres.NewMFARepository(sqlxDB)
	permRepo := postgres.NewPermissionRepository(sqlxDB)

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

	return &Container{
		UserRepo:       userRepo,
		CredRepo:       credRepo,
		SessionRepo:    sessionRepo,
		PermissionRepo: permRepo,
		PasswordHasher: hasher,
		TokenGenerator: tokenGen,
		RegisterUC:     registerUC,
		LoginUC:        loginUC,
		RefreshTokenUC: refreshUC,
		LogoutUC:       logoutUC,
	}, nil
}
