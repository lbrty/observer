package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/tag"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// TagUseCase handles tag operations within a project.
type TagUseCase struct {
	repo repository.TagRepository
}

// NewTagUseCase creates a TagUseCase.
func NewTagUseCase(repo repository.TagRepository) *TagUseCase {
	return &TagUseCase{repo: repo}
}

// List returns all tags for a project.
func (uc *TagUseCase) List(ctx context.Context, projectID string) ([]TagDTO, error) {
	tags, err := uc.repo.List(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	dtos := make([]TagDTO, len(tags))
	for i, t := range tags {
		dtos[i] = tagToDTO(t)
	}
	return dtos, nil
}

// Create creates a new tag in the project.
func (uc *TagUseCase) Create(ctx context.Context, projectID string, input CreateTagInput) (*TagDTO, error) {
	t := &tag.Tag{
		ID:        ulid.NewString(),
		ProjectID: projectID,
		Name:      input.Name,
		Color:     input.Color,
	}
	if err := uc.repo.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("create tag: %w", err)
	}
	dto := tagToDTO(t)
	return &dto, nil
}

// Update applies a partial update to a tag.
func (uc *TagUseCase) Update(ctx context.Context, id string, input UpdateTagInput) (*TagDTO, error) {
	t, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get tag for update: %w", err)
	}
	if input.Name != nil {
		t.Name = *input.Name
	}
	if input.Color != nil {
		t.Color = *input.Color
	}
	if err := uc.repo.Update(ctx, t); err != nil {
		return nil, fmt.Errorf("update tag: %w", err)
	}
	dto := tagToDTO(t)
	return &dto, nil
}

// Delete removes a tag.
func (uc *TagUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}
	return nil
}
