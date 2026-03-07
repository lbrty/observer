package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/pet"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
	"github.com/lbrty/observer/internal/usecase"
)

// PetUseCase handles pet operations within a project.
type PetUseCase struct {
	repo    repository.PetRepository
	tagRepo repository.PetTagRepository
}

// NewPetUseCase creates a PetUseCase.
func NewPetUseCase(repo repository.PetRepository, tagRepo repository.PetTagRepository) *PetUseCase {
	return &PetUseCase{repo: repo, tagRepo: tagRepo}
}

// List returns paginated pets.
func (uc *PetUseCase) List(ctx context.Context, projectID string, input ListPetsInput) (*ListPetsOutput, error) {
	page, perPage := usecase.ClampPagination(input.Page, input.PerPage)

	pets, total, err := uc.repo.List(ctx, projectID, input.Status, input.TagIDs, page, perPage)
	if err != nil {
		return nil, fmt.Errorf("list pets: %w", err)
	}

	ids := make([]string, len(pets))
	dtos := make([]PetDTO, len(pets))
	for i, p := range pets {
		ids[i] = p.ID
		dtos[i] = petToDTO(p)
	}

	tagMap, err := uc.tagRepo.ListBulk(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("list pet tags: %w", err)
	}
	for i := range dtos {
		if tags, ok := tagMap[dtos[i].ID]; ok {
			dtos[i].TagIDs = tags
		} else {
			dtos[i].TagIDs = []string{}
		}
	}

	return &ListPetsOutput{
		Pets:    dtos,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}

// Get returns a pet by ID.
func (uc *PetUseCase) Get(ctx context.Context, id string) (*PetDTO, error) {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get pet: %w", err)
	}
	dto := petToDTO(p)
	return &dto, nil
}

// Create creates a new pet.
func (uc *PetUseCase) Create(ctx context.Context, projectID string, input CreatePetInput) (*PetDTO, error) {
	p := &pet.Pet{
		ID:             ulid.NewString(),
		ProjectID:      projectID,
		OwnerID:        input.OwnerID,
		Name:           input.Name,
		Status:         pet.PetStatusUnknown,
		RegistrationID: input.RegistrationID,
		Notes:          input.Notes,
	}
	if input.Status != nil {
		p.Status = pet.PetStatus(*input.Status)
	}
	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create pet: %w", err)
	}
	dto := petToDTO(p)
	return &dto, nil
}

// Update applies a partial update to a pet.
func (uc *PetUseCase) Update(ctx context.Context, id string, input UpdatePetInput) (*PetDTO, error) {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get pet for update: %w", err)
	}
	if input.OwnerID != nil {
		p.OwnerID = input.OwnerID
	}
	if input.Name != nil {
		p.Name = *input.Name
	}
	if input.Status != nil {
		p.Status = pet.PetStatus(*input.Status)
	}
	if input.RegistrationID != nil {
		p.RegistrationID = input.RegistrationID
	}
	if input.Notes != nil {
		p.Notes = input.Notes
	}
	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("update pet: %w", err)
	}
	dto := petToDTO(p)
	return &dto, nil
}

// Delete removes a pet.
func (uc *PetUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete pet: %w", err)
	}
	return nil
}
