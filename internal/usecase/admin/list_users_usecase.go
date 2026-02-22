package admin

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
)

// ListUsersUseCase lists users with pagination and filtering.
type ListUsersUseCase struct {
	userRepo repository.UserRepository
}

// NewListUsersUseCase creates a ListUsersUseCase.
func NewListUsersUseCase(userRepo repository.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{userRepo: userRepo}
}

// Execute returns a paginated list of users.
func (uc *ListUsersUseCase) Execute(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PerPage < 1 {
		input.PerPage = 20
	}
	if input.PerPage > 100 {
		input.PerPage = 100
	}

	filter := user.UserListFilter{
		Page:     input.Page,
		PerPage:  input.PerPage,
		Search:   input.Search,
		Role:     input.Role,
		IsActive: input.IsActive,
	}

	users, total, err := uc.userRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	dtos := make([]UserDTO, len(users))
	for i, u := range users {
		dtos[i] = userToDTO(u)
	}

	return &ListUsersOutput{
		Users:   dtos,
		Total:   total,
		Page:    input.Page,
		PerPage: input.PerPage,
	}, nil
}

func userToDTO(u *user.User) UserDTO {
	return UserDTO{
		ID:         u.ID.String(),
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Email:      u.Email,
		Phone:      u.Phone,
		OfficeID:   u.OfficeID,
		Role:       string(u.Role),
		IsVerified: u.IsVerified,
		IsActive:   u.IsActive,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}
