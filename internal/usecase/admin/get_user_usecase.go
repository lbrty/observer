package admin

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/repository"
)

// GetUserUseCase retrieves a single user by ID.
type GetUserUseCase struct {
	userRepo repository.UserRepository
}

// NewGetUserUseCase creates a GetUserUseCase.
func NewGetUserUseCase(userRepo repository.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{userRepo: userRepo}
}

// Execute returns a user by ID.
func (uc *GetUserUseCase) Execute(ctx context.Context, id ulid.ULID) (*UserDTO, error) {
	u, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	dto := userToDTO(u)
	return &dto, nil
}
