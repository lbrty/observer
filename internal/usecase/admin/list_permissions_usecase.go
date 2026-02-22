package admin

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/repository"
)

// ListPermissionsUseCase lists project permissions.
type ListPermissionsUseCase struct {
	permRepo repository.PermissionRepository
}

// NewListPermissionsUseCase creates a ListPermissionsUseCase.
func NewListPermissionsUseCase(permRepo repository.PermissionRepository) *ListPermissionsUseCase {
	return &ListPermissionsUseCase{permRepo: permRepo}
}

// Execute returns all permissions for a project.
func (uc *ListPermissionsUseCase) Execute(ctx context.Context, projectID string) ([]PermissionDTO, error) {
	perms, err := uc.permRepo.List(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list permissions: %w", err)
	}

	dtos := make([]PermissionDTO, len(perms))
	for i, p := range perms {
		dtos[i] = permToDTO(p)
	}
	return dtos, nil
}

func permToDTO(p *project.ProjectPermission) PermissionDTO {
	return PermissionDTO{
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
}
