package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/repository"
)

// PetTagUseCase manages pet-tag associations.
type PetTagUseCase struct {
	repo repository.PetTagRepository
}

// NewPetTagUseCase creates a PetTagUseCase.
func NewPetTagUseCase(repo repository.PetTagRepository) *PetTagUseCase {
	return &PetTagUseCase{repo: repo}
}

// List returns tag IDs for a pet.
func (uc *PetTagUseCase) List(ctx context.Context, petID string) ([]string, error) {
	ids, err := uc.repo.List(ctx, petID)
	if err != nil {
		return nil, fmt.Errorf("list pet tags: %w", err)
	}
	return ids, nil
}

// Replace replaces all tags for a pet.
func (uc *PetTagUseCase) Replace(ctx context.Context, petID string, tagIDs []string) error {
	if err := uc.repo.ReplaceAll(ctx, petID, tagIDs); err != nil {
		return fmt.Errorf("replace pet tags: %w", err)
	}
	return nil
}
