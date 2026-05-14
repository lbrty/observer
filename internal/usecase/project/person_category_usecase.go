package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/person"
	"github.com/lbrty/observer/internal/repository"
)

// PersonCategoryUseCase manages person-category associations.
type PersonCategoryUseCase struct {
	repo       repository.PersonCategoryRepository
	personRepo repository.PersonRepository
}

// NewPersonCategoryUseCase creates a PersonCategoryUseCase.
func NewPersonCategoryUseCase(
	repo repository.PersonCategoryRepository,
	personRepo repository.PersonRepository,
) *PersonCategoryUseCase {
	return &PersonCategoryUseCase{repo: repo, personRepo: personRepo}
}

// List returns category IDs for a person.
func (uc *PersonCategoryUseCase) List(ctx context.Context, projectID, personID string) ([]string, error) {
	p, err := uc.personRepo.GetByID(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("list person categories: %w", err)
	}

	if p.ProjectID != projectID {
		return nil, fmt.Errorf("list person categories: %w", person.ErrPersonNotFound)
	}

	ids, err := uc.repo.List(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("list person categories: %w", err)
	}

	return ids, nil
}

// Replace replaces all categories for a person, verifying project membership first.
func (uc *PersonCategoryUseCase) Replace(ctx context.Context, projectID, personID string, categoryIDs []string) error {
	p, err := uc.personRepo.GetByID(ctx, personID)
	if err != nil {
		return fmt.Errorf("replace person categories: %w", err)
	}

	if p.ProjectID != projectID {
		return fmt.Errorf("replace person categories: %w", person.ErrPersonNotFound)
	}

	if err := uc.repo.ReplaceAll(ctx, personID, categoryIDs); err != nil {
		return fmt.Errorf("replace person categories: %w", err)
	}

	return nil
}
