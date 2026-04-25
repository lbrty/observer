package auth

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/domain/user"
)

// Me returns the current user's profile with MFA status.
func (uc *AuthUseCase) Me(ctx context.Context, userID ulid.ULID) (*MeDTO, error) {
	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &MeDTO{
		ID:         u.ID.String(),
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Email:      u.Email,
		Phone:      u.Phone,
		Role:       string(u.Role),
		IsVerified: u.IsVerified,
		MFAEnabled: uc.IsMFAEnabled(ctx, userID),
		OfficeID:   u.OfficeID,
		CreatedAt:  u.CreatedAt,
	}, nil
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
