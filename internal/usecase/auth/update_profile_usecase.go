package auth

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/repository"
)

// UpdateProfileUseCase lets an authenticated user update their own profile.
type UpdateProfileUseCase struct {
	userRepo repository.UserRepository
}

// NewUpdateProfileUseCase creates an UpdateProfileUseCase.
func NewUpdateProfileUseCase(userRepo repository.UserRepository) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{userRepo: userRepo}
}

// Execute applies profile changes for the given user.
func (uc *UpdateProfileUseCase) Execute(ctx context.Context, userID ulid.ULID, input UpdateProfileInput) (*UserDTO, error) {
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
