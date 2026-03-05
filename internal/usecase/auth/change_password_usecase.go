package auth

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
)

// ChangePasswordUseCase lets an authenticated user change their password.
type ChangePasswordUseCase struct {
	credRepo repository.CredentialsRepository
	hasher   crypto.PasswordHasher
}

// NewChangePasswordUseCase creates a ChangePasswordUseCase.
func NewChangePasswordUseCase(
	credRepo repository.CredentialsRepository,
	hasher crypto.PasswordHasher,
) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		credRepo: credRepo,
		hasher:   hasher,
	}
}

// Execute verifies the current password and replaces it with a new one.
func (uc *ChangePasswordUseCase) Execute(ctx context.Context, userID ulid.ULID, input ChangePasswordInput) error {
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
