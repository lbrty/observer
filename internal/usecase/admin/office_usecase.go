package admin

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/reference"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// OfficeUseCase handles CRUD operations for offices.
type OfficeUseCase struct {
	repo repository.OfficeRepository
}

// NewOfficeUseCase creates an OfficeUseCase.
func NewOfficeUseCase(repo repository.OfficeRepository) *OfficeUseCase {
	return &OfficeUseCase{repo: repo}
}

// List returns all offices.
func (uc *OfficeUseCase) List(ctx context.Context) ([]OfficeDTO, error) {
	offices, err := uc.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list offices: %w", err)
	}
	dtos := make([]OfficeDTO, len(offices))
	for i, o := range offices {
		dtos[i] = officeToDTO(o)
	}
	return dtos, nil
}

// Get returns an office by ID.
func (uc *OfficeUseCase) Get(ctx context.Context, id string) (*OfficeDTO, error) {
	o, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get office: %w", err)
	}
	dto := officeToDTO(o)
	return &dto, nil
}

// Create creates a new office.
func (uc *OfficeUseCase) Create(ctx context.Context, input CreateOfficeInput) (*OfficeDTO, error) {
	o := &reference.Office{
		ID:      ulid.NewString(),
		Name:    input.Name,
		PlaceID: input.PlaceID,
	}
	if err := uc.repo.Create(ctx, o); err != nil {
		return nil, fmt.Errorf("create office: %w", err)
	}
	dto := officeToDTO(o)
	return &dto, nil
}

// Update applies a partial update to an office.
func (uc *OfficeUseCase) Update(ctx context.Context, id string, input UpdateOfficeInput) (*OfficeDTO, error) {
	o, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get office for update: %w", err)
	}
	if input.Name != nil {
		o.Name = *input.Name
	}
	if input.PlaceID != nil {
		o.PlaceID = input.PlaceID
	}
	if err := uc.repo.Update(ctx, o); err != nil {
		return nil, fmt.Errorf("update office: %w", err)
	}
	dto := officeToDTO(o)
	return &dto, nil
}

// Delete removes an office by ID.
func (uc *OfficeUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete office: %w", err)
	}
	return nil
}

func officeToDTO(o *reference.Office) OfficeDTO {
	return OfficeDTO{
		ID:        o.ID,
		Name:      o.Name,
		PlaceID:   o.PlaceID,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}
