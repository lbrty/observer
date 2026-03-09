package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/pquerna/otp/totp"

	"github.com/lbrty/observer/internal/crypto"
	domainauth "github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/user"
	iulid "github.com/lbrty/observer/internal/ulid"
)

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
