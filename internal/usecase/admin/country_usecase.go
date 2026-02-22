package admin

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/reference"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// CountryUseCase handles CRUD operations for countries.
type CountryUseCase struct {
	repo repository.CountryRepository
}

// NewCountryUseCase creates a CountryUseCase.
func NewCountryUseCase(repo repository.CountryRepository) *CountryUseCase {
	return &CountryUseCase{repo: repo}
}

// List returns all countries.
func (uc *CountryUseCase) List(ctx context.Context) ([]CountryDTO, error) {
	countries, err := uc.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list countries: %w", err)
	}
	dtos := make([]CountryDTO, len(countries))
	for i, c := range countries {
		dtos[i] = countryToDTO(c)
	}
	return dtos, nil
}

// Get returns a country by ID.
func (uc *CountryUseCase) Get(ctx context.Context, id string) (*CountryDTO, error) {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get country: %w", err)
	}
	dto := countryToDTO(c)
	return &dto, nil
}

// Create creates a new country.
func (uc *CountryUseCase) Create(ctx context.Context, input CreateCountryInput) (*CountryDTO, error) {
	c := &reference.Country{
		ID:   ulid.NewString(),
		Name: input.Name,
		Code: input.Code,
	}
	if err := uc.repo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("create country: %w", err)
	}
	dto := countryToDTO(c)
	return &dto, nil
}

// Update applies a partial update to a country.
func (uc *CountryUseCase) Update(ctx context.Context, id string, input UpdateCountryInput) (*CountryDTO, error) {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get country for update: %w", err)
	}
	if input.Name != nil {
		c.Name = *input.Name
	}
	if input.Code != nil {
		c.Code = *input.Code
	}
	if err := uc.repo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("update country: %w", err)
	}
	dto := countryToDTO(c)
	return &dto, nil
}

// Delete removes a country by ID.
func (uc *CountryUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete country: %w", err)
	}
	return nil
}

func countryToDTO(c *reference.Country) CountryDTO {
	return CountryDTO{
		ID:        c.ID,
		Name:      c.Name,
		Code:      c.Code,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
