package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/person"
	"github.com/lbrty/observer/internal/repository"
)

// PersonTagUseCase manages person-tag associations.
type PersonTagUseCase struct {
	repo       repository.PersonTagRepository
	personRepo repository.PersonRepository
}

// NewPersonTagUseCase creates a PersonTagUseCase.
func NewPersonTagUseCase(repo repository.PersonTagRepository, personRepo repository.PersonRepository) *PersonTagUseCase {
	return &PersonTagUseCase{repo: repo, personRepo: personRepo}
}

// List returns tag IDs for a person.
func (uc *PersonTagUseCase) List(ctx context.Context, personID string) ([]string, error) {
	ids, err := uc.repo.List(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("list person tags: %w", err)
	}

	return ids, nil
}

// Replace replaces all tags for a person, verifying project membership first.
func (uc *PersonTagUseCase) Replace(ctx context.Context, projectID, personID string, tagIDs []string) error {
	p, err := uc.personRepo.GetByID(ctx, personID)
	if err != nil {
		return fmt.Errorf("replace person tags: %w", err)
	}

	if p.ProjectID != projectID {
		return fmt.Errorf("replace person tags: %w", person.ErrPersonNotFound)
	}

	if err := uc.repo.ReplaceAll(ctx, personID, tagIDs); err != nil {
		return fmt.Errorf("replace person tags: %w", err)
	}

	return nil
}
