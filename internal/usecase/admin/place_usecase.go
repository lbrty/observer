package admin

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/reference"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// PlaceUseCase handles CRUD operations for places.
type PlaceUseCase struct {
	repo repository.PlaceRepository
}

// NewPlaceUseCase creates a PlaceUseCase.
func NewPlaceUseCase(repo repository.PlaceRepository) *PlaceUseCase {
	return &PlaceUseCase{repo: repo}
}

// ListAll returns all places.
func (uc *PlaceUseCase) ListAll(ctx context.Context) ([]PlaceDTO, error) {
	places, err := uc.repo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list all places: %w", err)
	}
	dtos := make([]PlaceDTO, len(places))
	for i, p := range places {
		dtos[i] = placeToDTO(p)
	}
	return dtos, nil
}

// List returns all places for a state.
func (uc *PlaceUseCase) List(ctx context.Context, stateID string) ([]PlaceDTO, error) {
	places, err := uc.repo.List(ctx, stateID)
	if err != nil {
		return nil, fmt.Errorf("list places: %w", err)
	}
	dtos := make([]PlaceDTO, len(places))
	for i, p := range places {
		dtos[i] = placeToDTO(p)
	}
	return dtos, nil
}

// Get returns a place by ID.
func (uc *PlaceUseCase) Get(ctx context.Context, id string) (*PlaceDTO, error) {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get place: %w", err)
	}
	dto := placeToDTO(p)
	return &dto, nil
}

// Create creates a new place.
func (uc *PlaceUseCase) Create(ctx context.Context, stateID string, input CreatePlaceInput) (*PlaceDTO, error) {
	p := &reference.Place{
		ID:      ulid.NewString(),
		StateID: stateID,
		Name:    input.Name,
		Lat:     input.Lat,
		Lon:     input.Lon,
	}
	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create place: %w", err)
	}
	dto := placeToDTO(p)
	return &dto, nil
}

// Update applies a partial update to a place.
func (uc *PlaceUseCase) Update(ctx context.Context, id string, input UpdatePlaceInput) (*PlaceDTO, error) {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get place for update: %w", err)
	}
	if input.Name != nil {
		p.Name = *input.Name
	}
	if input.Lat != nil {
		p.Lat = input.Lat
	}
	if input.Lon != nil {
		p.Lon = input.Lon
	}
	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("update place: %w", err)
	}
	dto := placeToDTO(p)
	return &dto, nil
}

// Delete removes a place by ID.
func (uc *PlaceUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete place: %w", err)
	}
	return nil
}

func placeToDTO(p *reference.Place) PlaceDTO {
	return PlaceDTO{
		ID:        p.ID,
		StateID:   p.StateID,
		Name:      p.Name,
		Lat:       p.Lat,
		Lon:       p.Lon,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
