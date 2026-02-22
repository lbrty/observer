package admin

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/repository"
)

// UpdatePermissionUseCase updates an existing project permission.
type UpdatePermissionUseCase struct {
	permRepo repository.PermissionRepository
}

// NewUpdatePermissionUseCase creates an UpdatePermissionUseCase.
func NewUpdatePermissionUseCase(permRepo repository.PermissionRepository) *UpdatePermissionUseCase {
	return &UpdatePermissionUseCase{permRepo: permRepo}
}

// Execute applies a partial update to a project permission.
func (uc *UpdatePermissionUseCase) Execute(ctx context.Context, id string, input UpdatePermissionInput) (*PermissionDTO, error) {
	perm, err := uc.permRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get permission for update: %w", err)
	}

	if input.Role != nil {
		role, err := project.ValidateProjectRole(*input.Role)
		if err != nil {
			return nil, err
		}
		perm.Role = role
	}
	if input.CanViewContact != nil {
		perm.CanViewContact = *input.CanViewContact
	}
	if input.CanViewPersonal != nil {
		perm.CanViewPersonal = *input.CanViewPersonal
	}
	if input.CanViewDocuments != nil {
		perm.CanViewDocuments = *input.CanViewDocuments
	}

	if err := uc.permRepo.Update(ctx, perm); err != nil {
		return nil, fmt.Errorf("update permission: %w", err)
	}

	dto := permToDTO(perm)
	return &dto, nil
}
