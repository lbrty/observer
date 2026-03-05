package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/lbrty/observer/internal/crypto"
	domainauth "github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// LoginUseCase handles user login.
type LoginUseCase struct {
	userRepo    repository.UserRepository
	credRepo    repository.CredentialsRepository
	sessionRepo repository.SessionRepository
	mfaRepo     repository.MFARepository
	hasher      crypto.PasswordHasher
	tokenGen    crypto.TokenGenerator
}

// NewLoginUseCase creates a LoginUseCase.
func NewLoginUseCase(
	userRepo repository.UserRepository,
	credRepo repository.CredentialsRepository,
	sessionRepo repository.SessionRepository,
	mfaRepo repository.MFARepository,
	hasher crypto.PasswordHasher,
	tokenGen crypto.TokenGenerator,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:    userRepo,
		credRepo:    credRepo,
		sessionRepo: sessionRepo,
		mfaRepo:     mfaRepo,
		hasher:      hasher,
		tokenGen:    tokenGen,
	}
}

// Execute authenticates a user and returns tokens or an MFA challenge.
func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput, userAgent, ip string) (*LoginOutput, error) {
	u, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, user.ErrInvalidCredentials
	}

	if err := u.CanLogin(); err != nil {
		return nil, err
	}

	cred, err := uc.credRepo.GetByUserID(ctx, u.ID)
	if err != nil {
		return nil, user.ErrInvalidCredentials
	}

	if err := uc.hasher.Verify(input.Password, cred.PasswordHash, cred.Salt); err != nil {
		return nil, user.ErrInvalidCredentials
	}

	mfaCfg, err := uc.mfaRepo.GetByUserID(ctx, u.ID)
	if err == nil && mfaCfg.IsEnabled {
		mfaToken, err := uc.tokenGen.GenerateMFAToken(u.ID)
		if err != nil {
			return nil, fmt.Errorf("generate mfa token: %w", err)
		}
		return &LoginOutput{RequiresMFA: true, MFAToken: mfaToken}, nil
	}

	tokens, err := uc.createSession(ctx, u, userAgent, ip)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &LoginOutput{
		RequiresMFA: false,
		Tokens:      tokens,
		User:        toUserDTO(u),
	}, nil
}

func (uc *LoginUseCase) createSession(ctx context.Context, u *user.User, userAgent, ip string) (*TokenPair, error) {
	accessToken, expiresAt, err := uc.tokenGen.GenerateAccessToken(u.ID, string(u.Role))
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := uc.tokenGen.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	session := &domainauth.Session{
		ID:           ulid.New(),
		UserID:       u.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IP:           ip,
		ExpiresAt:    time.Now().UTC().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now().UTC(),
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("persist session: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func toUserDTO(u *user.User) *UserDTO {
	return &UserDTO{
		ID:         u.ID.String(),
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Email:      u.Email,
		Phone:      u.Phone,
		Role:       string(u.Role),
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt,
	}
}
