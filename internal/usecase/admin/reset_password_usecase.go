package admin

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/repository"
)

// ResetPasswordInput holds the new password for an admin reset.
type ResetPasswordInput struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResetPasswordUseCase allows an admin to reset a user's password.
type ResetPasswordUseCase struct {
	credRepo repository.CredentialsRepository
	hasher   crypto.PasswordHasher
}

// NewResetPasswordUseCase creates a ResetPasswordUseCase.
func NewResetPasswordUseCase(
	credRepo repository.CredentialsRepository,
	hasher crypto.PasswordHasher,
) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{credRepo: credRepo, hasher: hasher}
}

// Execute resets the password for the given user.
func (uc *ResetPasswordUseCase) Execute(ctx context.Context, userID ulid.ULID, input ResetPasswordInput) error {
	cred, err := uc.credRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get credentials: %w", err)
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
