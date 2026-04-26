package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/pet"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
	"github.com/lbrty/observer/internal/usecase"
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
)

// PetUseCase handles pet operations within a project.
type PetUseCase struct {
	repo    repository.PetRepository
	tagRepo repository.PetTagRepository
	auditUC *ucaudit.AuditUseCase
}

// NewPetUseCase creates a PetUseCase.
func NewPetUseCase(
	repo repository.PetRepository,
	tagRepo repository.PetTagRepository,
	auditUC *ucaudit.AuditUseCase,
) *PetUseCase {
	return &PetUseCase{repo: repo, tagRepo: tagRepo, auditUC: auditUC}
}

// List returns paginated pets.
func (uc *PetUseCase) List(ctx context.Context, projectID string, input ListPetsInput) (*ListPetsOutput, error) {
	page, perPage := usecase.ClampPagination(input.Page, input.PerPage)

	filter := pet.PetListFilter{
		ProjectID: projectID,
		TagIDs:    input.TagIDs,
		Page:      page,
		PerPage:   perPage,
	}
	if input.Status != "" {
		s := pet.PetStatus(input.Status)
		filter.Status = &s
	}

	pets, total, err := uc.repo.List(ctx, filter)
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
func (uc *PetUseCase) Get(ctx context.Context, projectID, id string) (*PetDTO, error) {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get pet: %w", err)
	}

	if p.ProjectID != projectID {
		return nil, pet.ErrPetNotFound
	}

	dto := petToDTO(p)
	tagMap, err := uc.tagRepo.ListBulk(ctx, []string{id})
	if err != nil {
		return nil, fmt.Errorf("list pet tags: %w", err)
	}

	if tags, ok := tagMap[id]; ok {
		dto.TagIDs = tags
	} else {
		dto.TagIDs = []string{}
	}

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

	uc.auditUC.Record(
		ctx,
		&projectID,
		"pet.create",
		"pet",
		&p.ID,
		fmt.Sprintf("Created pet %s", p.ID),
	)

	dto := petToDTO(p)
	return &dto, nil
}

// Update applies a partial update to a pet.
func (uc *PetUseCase) Update(ctx context.Context, projectID, id string, input UpdatePetInput) (*PetDTO, error) {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get pet for update: %w", err)
	}

	if p.ProjectID != projectID {
		return nil, pet.ErrPetNotFound
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

	uc.auditUC.Record(
		ctx,
		&projectID,
		"pet.update",
		"pet",
		&id,
		fmt.Sprintf("Updated pet %s", id),
	)

	dto := petToDTO(p)
	return &dto, nil
}

// Delete removes a pet.
func (uc *PetUseCase) Delete(ctx context.Context, projectID, id string) error {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get pet for delete: %w", err)
	}

	if p.ProjectID != projectID {
		return pet.ErrPetNotFound
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete pet: %w", err)
	}

	uc.auditUC.Record(
		ctx,
		&projectID,
		"pet.delete",
		"pet",
		&id,
		fmt.Sprintf("Deleted pet %s", id),
	)

	return nil
}
