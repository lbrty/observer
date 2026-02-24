package admin

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
)

// ListPermissionsUseCase lists project permissions with enriched user details.
type ListPermissionsUseCase struct {
	permRepo repository.PermissionRepository
	userRepo repository.UserRepository
}

// NewListPermissionsUseCase creates a ListPermissionsUseCase.
func NewListPermissionsUseCase(permRepo repository.PermissionRepository, userRepo repository.UserRepository) *ListPermissionsUseCase {
	return &ListPermissionsUseCase{permRepo: permRepo, userRepo: userRepo}
}

// Execute returns all permissions for a project with user details.
func (uc *ListPermissionsUseCase) Execute(ctx context.Context, projectID string) ([]PermissionMemberDTO, error) {
	perms, err := uc.permRepo.List(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list permissions: %w", err)
	}

	if len(perms) == 0 {
		return []PermissionMemberDTO{}, nil
	}

	ids := make([]ulid.ULID, 0, len(perms))
	seen := make(map[string]bool, len(perms))
	for _, p := range perms {
		if !seen[p.UserID] {
			seen[p.UserID] = true
			id, err := ulid.Parse(p.UserID)
			if err != nil {
				continue
			}
			ids = append(ids, id)
		}
	}

	users, err := uc.userRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("get users by ids: %w", err)
	}

	userMap := make(map[string]*user.User, len(users))
	for _, u := range users {
		userMap[u.ID.String()] = u
	}

	dtos := make([]PermissionMemberDTO, len(perms))
	for i, p := range perms {
		dtos[i] = permToMemberDTO(p, userMap[p.UserID])
	}
	return dtos, nil
}

func permToMemberDTO(p *project.ProjectPermission, u *user.User) PermissionMemberDTO {
	dto := PermissionMemberDTO{
		ID:               p.ID,
		ProjectID:        p.ProjectID,
		UserID:           p.UserID,
		Role:             string(p.Role),
		CanViewContact:   p.CanViewContact,
		CanViewPersonal:  p.CanViewPersonal,
		CanViewDocuments: p.CanViewDocuments,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}
	if u != nil {
		dto.UserFirstName = u.FirstName
		dto.UserLastName = u.LastName
		dto.UserEmail = u.Email
		dto.UserRole = string(u.Role)
	}
	return dto
}
