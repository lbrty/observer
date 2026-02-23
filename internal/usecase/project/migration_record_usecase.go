package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/migration"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// MigrationRecordUseCase handles migration record operations (append-only).
type MigrationRecordUseCase struct {
	repo repository.MigrationRecordRepository
}

// NewMigrationRecordUseCase creates a MigrationRecordUseCase.
func NewMigrationRecordUseCase(repo repository.MigrationRecordRepository) *MigrationRecordUseCase {
	return &MigrationRecordUseCase{repo: repo}
}

// ListByPerson returns all migration records for a person.
func (uc *MigrationRecordUseCase) ListByPerson(ctx context.Context, personID string) ([]MigrationRecordDTO, error) {
	records, err := uc.repo.ListByPerson(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("list migration records: %w", err)
	}
	dtos := make([]MigrationRecordDTO, len(records))
	for i, r := range records {
		dtos[i] = migrationRecordToDTO(r)
	}
	return dtos, nil
}

// Get returns a migration record by ID.
func (uc *MigrationRecordUseCase) Get(ctx context.Context, id string) (*MigrationRecordDTO, error) {
	r, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get migration record: %w", err)
	}
	dto := migrationRecordToDTO(r)
	return &dto, nil
}

// Create creates a new migration record.
func (uc *MigrationRecordUseCase) Create(ctx context.Context, personID string, input CreateMigrationRecordInput) (*MigrationRecordDTO, error) {
	r := &migration.Record{
		ID:                 ulid.NewString(),
		PersonID:           personID,
		FromPlaceID:        input.FromPlaceID,
		DestinationPlaceID: input.DestinationPlaceID,
		Notes:              input.Notes,
	}

	if input.MovementReason != nil {
		mr := migration.MovementReason(*input.MovementReason)
		r.MovementReason = &mr
	}
	if input.HousingAtDestination != nil {
		h := migration.HousingAtDestination(*input.HousingAtDestination)
		r.HousingAtDestination = &h
	}
	if err := parseDateField(input.MigrationDate, &r.MigrationDate); err != nil {
		return nil, fmt.Errorf("invalid migration_date: %w", err)
	}

	if err := uc.repo.Create(ctx, r); err != nil {
		return nil, fmt.Errorf("create migration record: %w", err)
	}
	dto := migrationRecordToDTO(r)
	return &dto, nil
}
