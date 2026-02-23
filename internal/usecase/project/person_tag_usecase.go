package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/repository"
)

// PersonTagUseCase manages person-tag associations.
type PersonTagUseCase struct {
	repo repository.PersonTagRepository
}

// NewPersonTagUseCase creates a PersonTagUseCase.
func NewPersonTagUseCase(repo repository.PersonTagRepository) *PersonTagUseCase {
	return &PersonTagUseCase{repo: repo}
}

// List returns tag IDs for a person.
func (uc *PersonTagUseCase) List(ctx context.Context, personID string) ([]string, error) {
	ids, err := uc.repo.List(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("list person tags: %w", err)
	}
	return ids, nil
}

// Replace replaces all tags for a person.
func (uc *PersonTagUseCase) Replace(ctx context.Context, personID string, tagIDs []string) error {
	if err := uc.repo.ReplaceAll(ctx, personID, tagIDs); err != nil {
		return fmt.Errorf("replace person tags: %w", err)
	}
	return nil
}
