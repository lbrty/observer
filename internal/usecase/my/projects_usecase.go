package my

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/repository"
)

// MyProjectsUseCase returns projects accessible to the current user.
type MyProjectsUseCase struct {
	permRepo    repository.PermissionRepository
	projectRepo repository.ProjectRepository
}

// NewMyProjectsUseCase creates a MyProjectsUseCase.
func NewMyProjectsUseCase(permRepo repository.PermissionRepository, projectRepo repository.ProjectRepository) *MyProjectsUseCase {
	return &MyProjectsUseCase{permRepo: permRepo, projectRepo: projectRepo}
}

// Execute returns the list of projects accessible to the given user.
func (uc *MyProjectsUseCase) Execute(ctx context.Context, userID string, role user.Role) (*MyProjectsOutput, error) {
	if role == user.RoleAdmin {
		return uc.adminProjects(ctx)
	}
	return uc.userProjects(ctx, userID)
}

func (uc *MyProjectsUseCase) adminProjects(ctx context.Context) (*MyProjectsOutput, error) {
	active := project.ProjectStatusActive
	projects, _, err := uc.projectRepo.List(ctx, project.ProjectListFilter{
		Status:  &active,
		Page:    1,
		PerPage: 1000,
	})
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}

	dtos := make([]MyProjectDTO, 0, len(projects))
	for _, p := range projects {
		dtos = append(dtos, MyProjectDTO{
			ID:               p.ID,
			Name:             p.Name,
			Description:      p.Description,
			Status:           string(p.Status),
			Role:             "owner",
			CanViewContact:   true,
			CanViewPersonal:  true,
			CanViewDocuments: true,
			CreatedAt:        p.CreatedAt,
			UpdatedAt:        p.UpdatedAt,
		})
	}
	return &MyProjectsOutput{Projects: dtos}, nil
}

func (uc *MyProjectsUseCase) userProjects(ctx context.Context, userID string) (*MyProjectsOutput, error) {
	perms, err := uc.permRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list user permissions: %w", err)
	}

	dtos := make([]MyProjectDTO, 0, len(perms))
	for _, perm := range perms {
		p, err := uc.projectRepo.GetByID(ctx, perm.ProjectID)
		if err != nil {
			continue
		}
		if p.Status != project.ProjectStatusActive {
			continue
		}
		dtos = append(dtos, MyProjectDTO{
			ID:               p.ID,
			Name:             p.Name,
			Description:      p.Description,
			Status:           string(p.Status),
			Role:             string(perm.Role),
			CanViewContact:   perm.CanViewContact,
			CanViewPersonal:  perm.CanViewPersonal,
			CanViewDocuments: perm.CanViewDocuments,
			CreatedAt:        p.CreatedAt,
			UpdatedAt:        p.UpdatedAt,
		})
	}
	return &MyProjectsOutput{Projects: dtos}, nil
}
