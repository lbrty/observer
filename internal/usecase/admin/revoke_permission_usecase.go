package admin

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/repository"
)

// RevokePermissionUseCase removes a project permission.
type RevokePermissionUseCase struct {
	permRepo repository.PermissionRepository
}

// NewRevokePermissionUseCase creates a RevokePermissionUseCase.
func NewRevokePermissionUseCase(permRepo repository.PermissionRepository) *RevokePermissionUseCase {
	return &RevokePermissionUseCase{permRepo: permRepo}
}

// Execute deletes a project permission by ID.
func (uc *RevokePermissionUseCase) Execute(ctx context.Context, id string) error {
	if err := uc.permRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("revoke permission: %w", err)
	}
	return nil
}
