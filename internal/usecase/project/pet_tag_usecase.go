package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/pet"
	"github.com/lbrty/observer/internal/repository"
)

// PetTagUseCase manages pet-tag associations.
type PetTagUseCase struct {
	repo    repository.PetTagRepository
	petRepo repository.PetRepository
}

// NewPetTagUseCase creates a PetTagUseCase.
func NewPetTagUseCase(repo repository.PetTagRepository, petRepo repository.PetRepository) *PetTagUseCase {
	return &PetTagUseCase{repo: repo, petRepo: petRepo}
}

// List returns tag IDs for a pet.
func (uc *PetTagUseCase) List(ctx context.Context, petID string) ([]string, error) {
	ids, err := uc.repo.List(ctx, petID)
	if err != nil {
		return nil, fmt.Errorf("list pet tags: %w", err)
	}
	return ids, nil
}

// Replace replaces all tags for a pet, verifying project membership first.
func (uc *PetTagUseCase) Replace(ctx context.Context, projectID, petID string, tagIDs []string) error {
	p, err := uc.petRepo.GetByID(ctx, petID)
	if err != nil {
		return fmt.Errorf("replace pet tags: %w", err)
	}
	if p.ProjectID != projectID {
		return fmt.Errorf("replace pet tags: %w", pet.ErrPetNotFound)
	}
	if err := uc.repo.ReplaceAll(ctx, petID, tagIDs); err != nil {
		return fmt.Errorf("replace pet tags: %w", err)
	}
	return nil
}
