package admin

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// AssignPermissionUseCase assigns a project permission to a user.
type AssignPermissionUseCase struct {
	permRepo repository.PermissionRepository
}

// NewAssignPermissionUseCase creates an AssignPermissionUseCase.
func NewAssignPermissionUseCase(permRepo repository.PermissionRepository) *AssignPermissionUseCase {
	return &AssignPermissionUseCase{permRepo: permRepo}
}

// Execute creates a new project permission.
func (uc *AssignPermissionUseCase) Execute(ctx context.Context, projectID string, input AssignPermissionInput) (*PermissionDTO, error) {
	role, err := project.ValidateProjectRole(input.Role)
	if err != nil {
		return nil, err
	}

	perm := &project.ProjectPermission{
		ID:               ulid.NewString(),
		ProjectID:        projectID,
		UserID:           input.UserID,
		Role:             role,
		CanViewContact:   input.CanViewContact,
		CanViewPersonal:  input.CanViewPersonal,
		CanViewDocuments: input.CanViewDocuments,
	}

	if err := uc.permRepo.Create(ctx, perm); err != nil {
		return nil, fmt.Errorf("assign permission: %w", err)
	}

	dto := permToDTO(perm)
	return &dto, nil
}
