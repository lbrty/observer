package admin

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
	"github.com/lbrty/observer/internal/usecase"
)

// ProjectUseCase handles CRUD operations for projects.
type ProjectUseCase struct {
	repo repository.ProjectRepository
}

// NewProjectUseCase creates a ProjectUseCase.
func NewProjectUseCase(repo repository.ProjectRepository) *ProjectUseCase {
	return &ProjectUseCase{repo: repo}
}

// List returns paginated projects with optional filters.
func (uc *ProjectUseCase) List(ctx context.Context, input ListProjectsInput) (*ListProjectsOutput, error) {
	filter := project.ProjectListFilter{
		OwnerID: input.OwnerID,
		Page:    input.Page,
		PerPage: input.PerPage,
	}
	if input.Status != nil {
		s := project.ProjectStatus(*input.Status)
		filter.Status = &s
	}

	projects, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}

	dtos := make([]ProjectDTO, len(projects))
	for i, p := range projects {
		dtos[i] = projectToDTO(p)
	}

	page, perPage := usecase.ClampPagination(input.Page, input.PerPage)

	return &ListProjectsOutput{
		Projects: dtos,
		Total:    total,
		Page:     page,
		PerPage:  perPage,
	}, nil
}

// Get returns a project by ID.
func (uc *ProjectUseCase) Get(ctx context.Context, id string) (*ProjectDTO, error) {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	dto := projectToDTO(p)
	return &dto, nil
}

// Create creates a new project.
func (uc *ProjectUseCase) Create(ctx context.Context, ownerID string, input CreateProjectInput) (*ProjectDTO, error) {
	p := &project.Project{
		ID:          ulid.NewString(),
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     ownerID,
		Status:      project.ProjectStatusActive,
	}
	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	dto := projectToDTO(p)
	return &dto, nil
}

// Update applies a partial update to a project (no hard delete — archive instead).
func (uc *ProjectUseCase) Update(ctx context.Context, id string, input UpdateProjectInput) (*ProjectDTO, error) {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get project for update: %w", err)
	}
	if input.Name != nil {
		p.Name = *input.Name
	}
	if input.Description != nil {
		p.Description = input.Description
	}
	if input.Status != nil {
		p.Status = project.ProjectStatus(*input.Status)
	}
	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}
	dto := projectToDTO(p)
	return &dto, nil
}
