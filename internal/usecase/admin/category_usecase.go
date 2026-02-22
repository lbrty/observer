package admin

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/reference"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// CategoryUseCase handles CRUD operations for categories.
type CategoryUseCase struct {
	repo repository.CategoryRepository
}

// NewCategoryUseCase creates a CategoryUseCase.
func NewCategoryUseCase(repo repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{repo: repo}
}

// List returns all categories.
func (uc *CategoryUseCase) List(ctx context.Context) ([]CategoryDTO, error) {
	categories, err := uc.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	dtos := make([]CategoryDTO, len(categories))
	for i, c := range categories {
		dtos[i] = categoryToDTO(c)
	}
	return dtos, nil
}

// Get returns a category by ID.
func (uc *CategoryUseCase) Get(ctx context.Context, id string) (*CategoryDTO, error) {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get category: %w", err)
	}
	dto := categoryToDTO(c)
	return &dto, nil
}

// Create creates a new category.
func (uc *CategoryUseCase) Create(ctx context.Context, input CreateCategoryInput) (*CategoryDTO, error) {
	c := &reference.Category{
		ID:          ulid.NewString(),
		Name:        input.Name,
		Description: input.Description,
	}
	if err := uc.repo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}
	dto := categoryToDTO(c)
	return &dto, nil
}

// Update applies a partial update to a category.
func (uc *CategoryUseCase) Update(ctx context.Context, id string, input UpdateCategoryInput) (*CategoryDTO, error) {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get category for update: %w", err)
	}
	if input.Name != nil {
		c.Name = *input.Name
	}
	if input.Description != nil {
		c.Description = input.Description
	}
	if err := uc.repo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}
	dto := categoryToDTO(c)
	return &dto, nil
}

// Delete removes a category by ID.
func (uc *CategoryUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return nil
}

func categoryToDTO(c *reference.Category) CategoryDTO {
	return CategoryDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
