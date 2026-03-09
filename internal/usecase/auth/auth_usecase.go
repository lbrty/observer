package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/pquerna/otp/totp"

	"github.com/lbrty/observer/internal/crypto"
	domainauth "github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
	iulid "github.com/lbrty/observer/internal/ulid"
)

// AuthUseCase handles authentication, session management, and user profile operations.
type AuthUseCase struct {
	userRepo    repository.UserRepository
	credRepo    repository.CredentialsRepository
	sessionRepo repository.SessionRepository
	mfaRepo     repository.MFARepository
	hasher      crypto.PasswordHasher
	tokenGen    crypto.TokenGenerator
}

// NewAuthUseCase creates an AuthUseCase.
func NewAuthUseCase(
	userRepo repository.UserRepository,
	credRepo repository.CredentialsRepository,
	sessionRepo repository.SessionRepository,
	mfaRepo repository.MFARepository,
	hasher crypto.PasswordHasher,
	tokenGen crypto.TokenGenerator,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		credRepo:    credRepo,
		sessionRepo: sessionRepo,
		mfaRepo:     mfaRepo,
		hasher:      hasher,
		tokenGen:    tokenGen,
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

	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	cred := &user.Credentials{
		UserID:       userID,
		PasswordHash: hash,
		Salt:         salt,
		UpdatedAt:    now,
	}

	if err := uc.credRepo.Create(ctx, cred); err != nil {
		return nil, fmt.Errorf("create credentials: %w", err)
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
func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput, userAgent, ip string) (*LoginOutput, error) {
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
		ExpiresAt:    time.Now().UTC().Add(7 * 24 * time.Hour),
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

// ChangePassword verifies the current password and replaces it with a new one.
func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID ulid.ULID, input ChangePasswordInput) error {
	cred, err := uc.credRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get credentials: %w", err)
	}

	if err := uc.hasher.Verify(input.CurrentPassword, cred.PasswordHash, cred.Salt); err != nil {
		return user.ErrInvalidCredentials
	}

	hash, salt, err := uc.hasher.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	cred.PasswordHash = hash
	cred.Salt = salt
	if err := uc.credRepo.Update(ctx, cred); err != nil {
		return fmt.Errorf("update credentials: %w", err)
	}

	return nil
}

// UpdateProfile applies profile changes for the given user.
func (uc *AuthUseCase) UpdateProfile(ctx context.Context, userID ulid.ULID, input UpdateProfileInput) (*UserDTO, error) {
	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if input.FirstName != nil {
		u.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		u.LastName = *input.LastName
	}
	if input.Phone != nil {
		u.Phone = *input.Phone
	}

	if err := uc.userRepo.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	dto := toUserDTO(u)
	return dto, nil
}

// IsMFAEnabled reports whether the user has MFA enabled.
func (uc *AuthUseCase) IsMFAEnabled(ctx context.Context, userID ulid.ULID) bool {
	cfg, err := uc.mfaRepo.GetByUserID(ctx, userID)
	return err == nil && cfg.IsEnabled
}

// VerifyMFA completes the two-step login flow:
//  1. Validates the short-lived MFA JWT issued by Login (proves the password was correct).
//  2. Checks the 6-digit TOTP code against the stored secret (RFC 6238, 30-second window).
//  3. On success, creates a full session and returns access + refresh tokens.
//
// The TOTP secret is never re-transmitted here — only the ephemeral 6-digit code is sent.
func (uc *AuthUseCase) VerifyMFA(ctx context.Context, input VerifyMFAInput, userAgent, ip string) (*LoginOutput, error) {
	claims, err := uc.tokenGen.ValidateMFAToken(input.MFAToken)
	if err != nil {
		return nil, domainauth.ErrInvalidMFACode
	}

	userID, err := iulid.Parse(claims.UserID)
	if err != nil {
		return nil, domainauth.ErrInvalidMFACode
	}

	mfaCfg, err := uc.mfaRepo.GetByUserID(ctx, userID)
	if err != nil || !mfaCfg.IsEnabled {
		return nil, domainauth.ErrInvalidMFACode
	}

	if !totp.Validate(input.TOTPCode, mfaCfg.Secret) {
		return nil, domainauth.ErrInvalidMFACode
	}

	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := u.CanLogin(); err != nil {
		return nil, err
	}

	tokens, err := uc.createSession(ctx, u, userAgent, ip)
	if err != nil {
		return nil, fmt.Errorf("create session after MFA: %w", err)
	}

	return &LoginOutput{
		RequiresMFA: false,
		Tokens:      tokens,
		User:        toUserDTO(u),
	}, nil
}

// SetupMFA generates a new TOTP secret and OTPAuth URL for the given user.
// It does NOT save anything — the secret is returned to the frontend exactly once so
// the user can scan the QR code. EnableMFA must be called with the secret + a valid
// TOTP code to persist and activate it.
func (uc *AuthUseCase) SetupMFA(ctx context.Context, userID ulid.ULID) (*MFASetupOutput, error) {
	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Observer",
		AccountName: u.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("generate totp key: %w", err)
	}
	return &MFASetupOutput{
		Secret:     key.Secret(),
		OTPAuthURL: key.URL(),
	}, nil
}

// EnableMFA activates TOTP-based MFA for the user:
//  1. Verifies the provided TOTP code against the secret (confirms the user's app is synced).
//  2. Upserts the mfa_configs row with is_enabled=true.
//
// After this call, Login will return RequiresMFA=true for this user.
func (uc *AuthUseCase) EnableMFA(ctx context.Context, userID ulid.ULID, input EnableMFAInput) error {
	if !totp.Validate(input.TOTPCode, input.Secret) {
		return domainauth.ErrInvalidMFACode
	}
	cfg := &user.MFAConfig{
		UserID:    userID,
		Method:    "totp",
		Secret:    input.Secret,
		IsEnabled: true,
		CreatedAt: time.Now().UTC(),
	}
	return uc.mfaRepo.Upsert(ctx, cfg)
}

// DisableMFA requires a valid TOTP code to confirm intent, then clears the secret
// and sets is_enabled=false. Future logins will skip the MFA challenge.
func (uc *AuthUseCase) DisableMFA(ctx context.Context, userID ulid.ULID, input DisableMFAInput) error {
	cfg, err := uc.mfaRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if !totp.Validate(input.TOTPCode, cfg.Secret) {
		return domainauth.ErrInvalidMFACode
	}
	cfg.IsEnabled = false
	cfg.Secret = ""
	return uc.mfaRepo.Upsert(ctx, cfg)
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
