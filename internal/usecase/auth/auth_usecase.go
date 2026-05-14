package auth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lbrty/observer/internal/crypto"
	domainauth "github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
	iulid "github.com/lbrty/observer/internal/ulid"
)

// AuthUseCase handles authentication, session management, and user profile operations.
type AuthUseCase struct {
	userRepo      repository.UserRepository
	credRepo      repository.CredentialsRepository
	sessionRepo   repository.SessionRepository
	mfaRepo       repository.MFARepository
	recoveryRepo  repository.MFARecoveryCodeRepository
	hasher        crypto.PasswordHasher
	tokenGen      crypto.TokenGenerator
	loginAttempts repository.LoginAttemptStore
	refreshTTL    time.Duration
}

// NewAuthUseCase creates an AuthUseCase.
func NewAuthUseCase(
	userRepo repository.UserRepository,
	credRepo repository.CredentialsRepository,
	sessionRepo repository.SessionRepository,
	mfaRepo repository.MFARepository,
	recoveryRepo repository.MFARecoveryCodeRepository,
	hasher crypto.PasswordHasher,
	tokenGen crypto.TokenGenerator,
	loginAttempts repository.LoginAttemptStore,
	refreshTTL time.Duration,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:      userRepo,
		credRepo:      credRepo,
		sessionRepo:   sessionRepo,
		mfaRepo:       mfaRepo,
		recoveryRepo:  recoveryRepo,
		hasher:        hasher,
		tokenGen:      tokenGen,
		loginAttempts: loginAttempts,
		refreshTTL:    refreshTTL,
	}
}

// Register registers a new user.
// Email enumeration is prevented by returning success silently when the address is already registered.
// Role is always set to RoleGuest — admin promotes users via /admin/users/:id.
func (uc *AuthUseCase) Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	// Email already registered — return success silently to prevent enumeration.
	if _, err := uc.userRepo.GetByEmail(ctx, input.Email); err == nil {
		return &RegisterOutput{
			Message: "Registration successful. Your account is pending admin approval.",
		}, nil
	}

	hash, salt, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	userID := iulid.New()
	now := time.Now().UTC()

	newUser := &user.User{
		ID:         userID,
		Email:      input.Email,
		Role:       user.RoleGuest, // always guest; admin promotes via /admin/users/:id
		IsVerified: false,
		IsActive:   false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	cred := &user.Credentials{
		UserID:       userID,
		PasswordHash: hash,
		Salt:         salt,
		UpdatedAt:    now,
	}

	if err := uc.userRepo.CreateWithCredentials(ctx, newUser, cred); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &RegisterOutput{
		UserID:  userID.String(),
		Message: "Registration successful. Your account is pending admin approval.",
	}, nil
}

// Login authenticates a user with email + password and returns one of two outcomes:
//   - MFA enabled: issues a short-lived MFA token (no session yet) and sets RequiresMFA=true.
//     The client must follow up with POST /auth/mfa { mfa_token, totp_code } to complete login.
//   - MFA disabled: creates a full session and returns access + refresh tokens immediately.
//
// Lockout checks are enforced here: the call fails closed if the rate limiter store is unavailable.
func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput, userAgent, ip string) (*LoginOutput, error) {
	remaining, err := uc.loginAttempts.IsLocked(ctx, input.Email)
	if err != nil {
		slog.ErrorContext(ctx, "check login lockout", slog.Any("err", err))
		return nil, domainauth.ErrRateLimiterUnavailable
	}
	if remaining != 0 {
		if remaining > 0 {
			return nil, domainauth.ErrAccountTemporarilyLocked
		}
		return nil, domainauth.ErrAccountLocked
	}

	u, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		uc.recordFailure(ctx, input.Email)
		return nil, user.ErrInvalidCredentials
	}

	if err := u.CanLogin(); err != nil {
		return nil, err
	}

	cred, err := uc.credRepo.GetByUserID(ctx, u.ID)
	if err != nil {
		uc.recordFailure(ctx, input.Email)
		return nil, user.ErrInvalidCredentials
	}

	if err := uc.hasher.Verify(input.Password, cred.PasswordHash, cred.Salt); err != nil {
		return nil, uc.recordFailureOrLock(ctx, input.Email)
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

	if clearErr := uc.loginAttempts.ClearAttempts(ctx, input.Email); clearErr != nil {
		slog.ErrorContext(ctx, "clear login attempts", slog.Any("err", clearErr))
	}

	return &LoginOutput{
		RequiresMFA: false,
		Tokens:      tokens,
		User:        toUserDTO(u),
	}, nil
}

// recordFailure increments the failure counter; logs but does not propagate store errors.
func (uc *AuthUseCase) recordFailure(ctx context.Context, email string) {
	if _, err := uc.loginAttempts.RecordFailure(ctx, email); err != nil {
		slog.ErrorContext(ctx, "record login failure", slog.Any("err", err))
	}
}

// recordFailureOrLock increments the failure counter and returns a lockout error if the
// threshold was just reached, or ErrInvalidCredentials otherwise.
func (uc *AuthUseCase) recordFailureOrLock(ctx context.Context, email string) error {
	lockDur, err := uc.loginAttempts.RecordFailure(ctx, email)
	if err != nil {
		slog.ErrorContext(ctx, "record login failure", slog.Any("err", err))
		return user.ErrInvalidCredentials
	}
	if lockDur != 0 {
		if lockDur > 0 {
			return domainauth.ErrAccountTemporarilyLocked
		}
		return domainauth.ErrAccountLocked
	}
	return user.ErrInvalidCredentials
}

func (uc *AuthUseCase) createSession(ctx context.Context, u *user.User, userAgent, ip string) (*TokenPair, error) {
	accessToken, expiresAt, err := uc.tokenGen.GenerateAccessToken(u.ID, string(u.Role))
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := uc.tokenGen.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	session := &domainauth.Session{
		ID:           iulid.New(),
		UserID:       u.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IP:           ip,
		ExpiresAt:    time.Now().UTC().Add(uc.refreshTTL),
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

// RefreshToken rotates the refresh token and issues a new access token.
func (uc *AuthUseCase) RefreshToken(ctx context.Context, input RefreshTokenInput) (*TokenPair, error) {
	session, err := uc.sessionRepo.GetByRefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return nil, domainauth.ErrSessionNotFound
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		_ = uc.sessionRepo.Delete(ctx, session.ID)
		return nil, domainauth.ErrSessionExpired
	}

	if input.IP != "" && session.IP != input.IP {
		slog.Warn("token refresh from new IP",
			slog.String("session_ip", session.IP),
			slog.String("request_ip", input.IP),
			slog.String("user_id", session.UserID.String()),
		)
	}
	if input.UserAgent != "" && session.UserAgent != input.UserAgent {
		slog.Warn("token refresh from new user-agent",
			slog.String("user_id", session.UserID.String()),
		)
	}

	if err := uc.sessionRepo.Delete(ctx, session.ID); err != nil {
		return nil, fmt.Errorf("delete old session: %w", err)
	}

	u, err := uc.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if err := u.CanLogin(); err != nil {
		return nil, err
	}

	accessToken, expiresAt, err := uc.tokenGen.GenerateAccessToken(u.ID, string(u.Role))
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	newRefreshToken, err := uc.tokenGen.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	newSession := &domainauth.Session{
		ID:           iulid.New(),
		UserID:       session.UserID,
		RefreshToken: newRefreshToken,
		UserAgent:    session.UserAgent,
		IP:           session.IP,
		ExpiresAt:    time.Now().UTC().Add(uc.refreshTTL),
		CreatedAt:    time.Now().UTC(),
	}

	if err := uc.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, fmt.Errorf("create new session: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Logout deletes the session associated with the given refresh token.
func (uc *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	return uc.sessionRepo.DeleteByRefreshToken(ctx, refreshToken)
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
