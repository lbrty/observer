package admin

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
)

// UpdateUserUseCase handles partial user updates.
type UpdateUserUseCase struct {
	userRepo repository.UserRepository
}

// NewUpdateUserUseCase creates an UpdateUserUseCase.
func NewUpdateUserUseCase(userRepo repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{userRepo: userRepo}
}

// Execute applies a partial update to the user and returns the updated DTO.
func (uc *UpdateUserUseCase) Execute(ctx context.Context, id ulid.ULID, input UpdateUserInput) (*UserDTO, error) {
	u, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user for update: %w", err)
	}

	if input.FirstName != nil {
		u.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		u.LastName = *input.LastName
	}
	if input.Email != nil {
		u.Email = *input.Email
	}
	if input.Phone != nil {
		u.Phone = *input.Phone
	}
	if input.OfficeID != nil {
		u.OfficeID = input.OfficeID
	}
	if input.Role != nil {
		role, err := user.ValidateRole(*input.Role)
		if err != nil {
			return nil, err
		}
		u.Role = role
	}
	if input.IsActive != nil {
		u.IsActive = *input.IsActive
	}
	if input.IsVerified != nil {
		u.IsVerified = *input.IsVerified
	}

	if err := uc.userRepo.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	dto := userToDTO(u)
	return &dto, nil
}
