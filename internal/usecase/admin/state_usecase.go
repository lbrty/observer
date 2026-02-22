package admin

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/reference"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// StateUseCase handles CRUD operations for states.
type StateUseCase struct {
	repo repository.StateRepository
}

// NewStateUseCase creates a StateUseCase.
func NewStateUseCase(repo repository.StateRepository) *StateUseCase {
	return &StateUseCase{repo: repo}
}

// List returns all states for a country.
func (uc *StateUseCase) List(ctx context.Context, countryID string) ([]StateDTO, error) {
	states, err := uc.repo.List(ctx, countryID)
	if err != nil {
		return nil, fmt.Errorf("list states: %w", err)
	}
	dtos := make([]StateDTO, len(states))
	for i, s := range states {
		dtos[i] = stateToDTO(s)
	}
	return dtos, nil
}

// Get returns a state by ID.
func (uc *StateUseCase) Get(ctx context.Context, id string) (*StateDTO, error) {
	s, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get state: %w", err)
	}
	dto := stateToDTO(s)
	return &dto, nil
}

// Create creates a new state.
func (uc *StateUseCase) Create(ctx context.Context, countryID string, input CreateStateInput) (*StateDTO, error) {
	s := &reference.State{
		ID:           ulid.NewString(),
		CountryID:    countryID,
		Name:         input.Name,
		Code:         input.Code,
		ConflictZone: input.ConflictZone,
	}
	if err := uc.repo.Create(ctx, s); err != nil {
		return nil, fmt.Errorf("create state: %w", err)
	}
	dto := stateToDTO(s)
	return &dto, nil
}

// Update applies a partial update to a state.
func (uc *StateUseCase) Update(ctx context.Context, id string, input UpdateStateInput) (*StateDTO, error) {
	s, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get state for update: %w", err)
	}
	if input.Name != nil {
		s.Name = *input.Name
	}
	if input.Code != nil {
		s.Code = input.Code
	}
	if input.ConflictZone != nil {
		s.ConflictZone = input.ConflictZone
	}
	if err := uc.repo.Update(ctx, s); err != nil {
		return nil, fmt.Errorf("update state: %w", err)
	}
	dto := stateToDTO(s)
	return &dto, nil
}

// Delete removes a state by ID.
func (uc *StateUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete state: %w", err)
	}
	return nil
}

func stateToDTO(s *reference.State) StateDTO {
	return StateDTO{
		ID:           s.ID,
		CountryID:    s.CountryID,
		Name:         s.Name,
		Code:         s.Code,
		ConflictZone: s.ConflictZone,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}
