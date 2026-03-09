package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
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
	userRepo     repository.UserRepository
	credRepo     repository.CredentialsRepository
	sessionRepo  repository.SessionRepository
	mfaRepo      repository.MFARepository
	recoveryRepo repository.MFARecoveryCodeRepository
	hasher       crypto.PasswordHasher
	tokenGen     crypto.TokenGenerator
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
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:     userRepo,
		credRepo:     credRepo,
		sessionRepo:  sessionRepo,
		mfaRepo:      mfaRepo,
		recoveryRepo: recoveryRepo,
		hasher:       hasher,
		tokenGen:     tokenGen,
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

	if err := uc.sessionRepo.DeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("invalidate sessions: %w", err)
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
		// Try as a single-use recovery code.
		hash, _, hashErr := uc.hasher.Hash(input.TOTPCode)
		if hashErr != nil {
			return nil, domainauth.ErrInvalidMFACode
		}
		rc, rcErr := uc.recoveryRepo.FindUnused(ctx, userID, hash)
		if rcErr != nil {
			return nil, domainauth.ErrInvalidMFACode
		}
		_ = uc.recoveryRepo.MarkUsed(ctx, rc.ID)
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
//  3. Generates 8 single-use recovery codes and returns them (shown once only).
//
// After this call, Login will return RequiresMFA=true for this user.
func (uc *AuthUseCase) EnableMFA(ctx context.Context, userID ulid.ULID, input EnableMFAInput) (*EnableMFAOutput, error) {
	if !totp.Validate(input.TOTPCode, input.Secret) {
		return nil, domainauth.ErrInvalidMFACode
	}
	cfg := &user.MFAConfig{
		UserID:    userID,
		Method:    "totp",
		Secret:    input.Secret,
		IsEnabled: true,
		CreatedAt: time.Now().UTC(),
	}
	if err := uc.mfaRepo.Upsert(ctx, cfg); err != nil {
		return nil, err
	}

	plainCodes, records, err := generateRecoveryCodes(8, userID, uc.hasher)
	if err != nil {
		return nil, fmt.Errorf("generate recovery codes: %w", err)
	}
	if err := uc.recoveryRepo.CreateBatch(ctx, records); err != nil {
		return nil, fmt.Errorf("store recovery codes: %w", err)
	}
	return &EnableMFAOutput{RecoveryCodes: plainCodes}, nil
}

// generateRecoveryCodes returns n plain-text codes and their hashed records.
func generateRecoveryCodes(n int, userID ulid.ULID, hasher crypto.PasswordHasher) ([]string, []*user.MFARecoveryCode, error) {
	now := time.Now().UTC()
	plainCodes := make([]string, n)
	records := make([]*user.MFARecoveryCode, n)
	for i := range plainCodes {
		b := make([]byte, 8)
		if _, err := rand.Read(b); err != nil {
			return nil, nil, err
		}
		plain := hex.EncodeToString(b)
		hash, _, err := hasher.Hash(plain)
		if err != nil {
			return nil, nil, err
		}
		plainCodes[i] = plain
		records[i] = &user.MFARecoveryCode{
			ID:        iulid.NewString(),
			UserID:    userID,
			CodeHash:  hash,
			CreatedAt: now,
		}
	}
	return plainCodes, records, nil
}

// DisableMFA requires a valid TOTP code to confirm intent, then clears the secret,
// sets is_enabled=false, and deletes all recovery codes. Future logins skip the MFA challenge.
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
	if err := uc.mfaRepo.Upsert(ctx, cfg); err != nil {
		return err
	}
	_ = uc.recoveryRepo.DeleteByUserID(ctx, userID)
	return nil
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
