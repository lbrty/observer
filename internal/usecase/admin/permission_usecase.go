package admin

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
	iulid "github.com/lbrty/observer/internal/ulid"
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
)

// PermissionUseCase handles project permission management.
type PermissionUseCase struct {
	permRepo repository.PermissionRepository
	userRepo repository.UserRepository
	auditUC  *ucaudit.AuditUseCase
}

// NewPermissionUseCase creates a PermissionUseCase.
func NewPermissionUseCase(
	permRepo repository.PermissionRepository,
	userRepo repository.UserRepository,
	auditUC *ucaudit.AuditUseCase,
) *PermissionUseCase {
	return &PermissionUseCase{permRepo: permRepo, userRepo: userRepo, auditUC: auditUC}
}

// List returns permissions for a project with user details.
// Admin and Staff see all permissions; other roles see permissions
// only if they are a member of the project, otherwise get an empty list.
func (uc *PermissionUseCase) List(ctx context.Context, projectID string, callerID string, callerRole user.Role) ([]PermissionMemberDTO, error) {
	if callerRole != user.RoleAdmin && callerRole != user.RoleStaff {
		userPerms, err := uc.permRepo.ListByUserID(ctx, callerID)
		if err != nil {
			return nil, fmt.Errorf("check user permissions: %w", err)
		}
		allowed := false
		for _, p := range userPerms {
			if p.ProjectID == projectID {
				allowed = true
				break
			}
		}
		if !allowed {
			return []PermissionMemberDTO{}, nil
		}
	}

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

// Assign creates a new project permission.
func (uc *PermissionUseCase) Assign(ctx context.Context, projectID string, input AssignPermissionInput) (*PermissionDTO, error) {
	role, err := project.ValidateProjectRole(input.Role)
	if err != nil {
		return nil, err
	}

	uid, err := ulid.Parse(input.UserID)
	if err != nil {
		return nil, user.ErrUserNotFound
	}

	if _, err := uc.userRepo.GetByID(ctx, uid); err != nil {
		return nil, fmt.Errorf("verify user for assign: %w", err)
	}

	perm := &project.ProjectPermission{
		ID:               iulid.NewString(),
		ProjectID:        projectID,
		UserID:           input.UserID,
		Role:             role,
		CanViewContact:   input.CanViewContact,
		CanViewPersonal:  input.CanViewPersonal,
		CanViewDocuments: input.CanViewDocuments,
		CanExport:        input.CanExport,
	}

	if err := uc.permRepo.Create(ctx, perm); err != nil {
		return nil, fmt.Errorf("assign permission: %w", err)
	}

	uc.auditUC.Record(
		ctx,
		&projectID,
		"permission.grant",
		"permission",
		&perm.ID,
		fmt.Sprintf("Granted %s to user %s", input.Role, input.UserID),
		nil,
	)

	dto := permToDTO(perm)
	return &dto, nil
}

// Update applies a partial update to a project permission.
func (uc *PermissionUseCase) Update(ctx context.Context, projectID, id string, input UpdatePermissionInput) (*PermissionDTO, error) {
	perm, err := uc.permRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get permission for update: %w", err)
	}

	if perm.ProjectID != projectID {
		return nil, project.ErrPermissionNotFound
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
	if input.CanExport != nil {
		perm.CanExport = *input.CanExport
	}

	if err := uc.permRepo.Update(ctx, perm); err != nil {
		return nil, fmt.Errorf("update permission: %w", err)
	}

	uc.auditUC.Record(
		ctx,
		&projectID,
		"admin.permission.update",
		"permission",
		&id,
		fmt.Sprintf("Updated permission %s", id),
		nil,
	)

	dto := permToDTO(perm)
	return &dto, nil
}

// Revoke deletes a project permission by ID.
func (uc *PermissionUseCase) Revoke(ctx context.Context, projectID, id string) error {
	perm, err := uc.permRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get permission for revoke: %w", err)
	}

	if perm.ProjectID != projectID {
		return project.ErrPermissionNotFound
	}

	if err := uc.permRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("revoke permission: %w", err)
	}

	uc.auditUC.Record(
		ctx,
		&projectID,
		"permission.revoke",
		"permission",
		&id,
		fmt.Sprintf("Revoked permission %s", id),
		nil,
	)

	return nil
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
		CanExport:        p.CanExport,
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
