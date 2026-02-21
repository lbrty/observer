package app

import (
	"fmt"

	appauth "github.com/lbrty/observer/internal/application/auth"
	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/database"
	domainauth "github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/user"
	infraauth "github.com/lbrty/observer/internal/infrastructure/auth"
	"github.com/lbrty/observer/internal/infrastructure/persistence/postgres"
)

// Container holds all application dependencies.
type Container struct {
	// Repositories
	UserRepo    user.UserRepository
	CredRepo    user.CredentialsRepository
	SessionRepo domainauth.SessionRepository

	// Services
	PasswordHasher infraauth.PasswordHasher
	TokenGenerator infraauth.TokenGenerator

	// Use Cases
	RegisterUC     *appauth.RegisterUseCase
	LoginUC        *appauth.LoginUseCase
	RefreshTokenUC *appauth.RefreshTokenUseCase
	LogoutUC       *appauth.LogoutUseCase
}

// NewContainer wires all dependencies from config and database.
func NewContainer(cfg *config.Config, db database.DB) (*Container, error) {
	rsaKeys, err := infraauth.LoadRSAKeys(cfg.JWT.PrivateKeyPath, cfg.JWT.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load RSA keys: %w", err)
	}

	sqlxDB := db.GetDB()

	userRepo := postgres.NewUserRepository(sqlxDB)
	credRepo := postgres.NewCredentialsRepository(sqlxDB)
	sessionRepo := postgres.NewSessionRepository(sqlxDB)
	mfaRepo := postgres.NewMFARepository(sqlxDB)

	hasher := infraauth.NewArgonHasher()
	tokenGen := infraauth.NewRSATokenGenerator(
		rsaKeys,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
		cfg.JWT.MFATempTTL,
		cfg.JWT.Issuer,
	)

	registerUC := appauth.NewRegisterUseCase(userRepo, credRepo, hasher)
	loginUC := appauth.NewLoginUseCase(userRepo, credRepo, sessionRepo, mfaRepo, hasher, tokenGen)
	refreshUC := appauth.NewRefreshTokenUseCase(sessionRepo, tokenGen)
	logoutUC := appauth.NewLogoutUseCase(sessionRepo)

	return &Container{
		UserRepo:       userRepo,
		CredRepo:       credRepo,
		SessionRepo:    sessionRepo,
		PasswordHasher: hasher,
		TokenGenerator: tokenGen,
		RegisterUC:     registerUC,
		LoginUC:        loginUC,
		RefreshTokenUC: refreshUC,
		LogoutUC:       logoutUC,
	}, nil
}
