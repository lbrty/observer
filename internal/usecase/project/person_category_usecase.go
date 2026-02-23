package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/repository"
)

// PersonCategoryUseCase manages person-category associations.
type PersonCategoryUseCase struct {
	repo repository.PersonCategoryRepository
}

// NewPersonCategoryUseCase creates a PersonCategoryUseCase.
func NewPersonCategoryUseCase(repo repository.PersonCategoryRepository) *PersonCategoryUseCase {
	return &PersonCategoryUseCase{repo: repo}
}

// List returns category IDs for a person.
func (uc *PersonCategoryUseCase) List(ctx context.Context, personID string) ([]string, error) {
	ids, err := uc.repo.List(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("list person categories: %w", err)
	}
	return ids, nil
}

// Replace replaces all categories for a person.
func (uc *PersonCategoryUseCase) Replace(ctx context.Context, personID string, categoryIDs []string) error {
	if err := uc.repo.ReplaceAll(ctx, personID, categoryIDs); err != nil {
		return fmt.Errorf("replace person categories: %w", err)
	}
	return nil
}
