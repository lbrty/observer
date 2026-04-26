package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/migration"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
)

// MigrationRecordUseCase handles migration record operations.
type MigrationRecordUseCase struct {
	repo    repository.MigrationRecordRepository
	auditUC *ucaudit.AuditUseCase
}

// NewMigrationRecordUseCase creates a MigrationRecordUseCase.
func NewMigrationRecordUseCase(
	repo repository.MigrationRecordRepository,
	auditUC *ucaudit.AuditUseCase,
) *MigrationRecordUseCase {
	return &MigrationRecordUseCase{repo: repo, auditUC: auditUC}
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
func (uc *MigrationRecordUseCase) Get(ctx context.Context, personID, id string) (*MigrationRecordDTO, error) {
	r, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get migration record: %w", err)
	}

	if r.PersonID != personID {
		return nil, migration.ErrRecordNotFound
	}

	dto := migrationRecordToDTO(r)
	return &dto, nil
}

// Create creates a new migration record.
func (uc *MigrationRecordUseCase) Create(ctx context.Context, projectID, personID string, input CreateMigrationRecordInput) (*MigrationRecordDTO, error) {
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

	uc.auditUC.Record(
		ctx,
		&projectID,
		"migration_record.create",
		"migration_record",
		&r.ID,
		fmt.Sprintf("Created migration record %s", r.ID),
	)

	dto := migrationRecordToDTO(r)
	return &dto, nil
}

// Delete removes a migration record after verifying it belongs to the given person.
func (uc *MigrationRecordUseCase) Delete(ctx context.Context, projectID, personID, id string) error {
	rec, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get migration record for delete: %w", err)
	}

	if rec.PersonID != personID {
		return migration.ErrRecordNotFound
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete migration record: %w", err)
	}

	uc.auditUC.Record(
		ctx,
		&projectID,
		"migration_record.delete",
		"migration_record",
		&id,
		fmt.Sprintf("Deleted migration record %s", id),
	)

	return nil
}

// Update updates a migration record.
func (uc *MigrationRecordUseCase) Update(ctx context.Context, projectID, personID, id string, input UpdateMigrationRecordInput) (*MigrationRecordDTO, error) {
	r, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get migration record for update: %w", err)
	}

	if r.PersonID != personID {
		return nil, migration.ErrRecordNotFound
	}

	if input.FromPlaceID != nil {
		r.FromPlaceID = input.FromPlaceID
	}
	if input.DestinationPlaceID != nil {
		r.DestinationPlaceID = input.DestinationPlaceID
	}
	if input.MovementReason != nil {
		mr := migration.MovementReason(*input.MovementReason)
		r.MovementReason = &mr
	}
	if input.HousingAtDestination != nil {
		h := migration.HousingAtDestination(*input.HousingAtDestination)
		r.HousingAtDestination = &h
	}
	if input.Notes != nil {
		r.Notes = input.Notes
	}

	if err := parseDateField(input.MigrationDate, &r.MigrationDate); err != nil {
		return nil, fmt.Errorf("invalid migration_date: %w", err)
	}

	if err := uc.repo.Update(ctx, r); err != nil {
		return nil, fmt.Errorf("update migration record: %w", err)
	}

	uc.auditUC.Record(
		ctx,
		&projectID,
		"migration_record.update",
		"migration_record",
		&id,
		fmt.Sprintf("Updated migration record %s", id),
	)

	dto := migrationRecordToDTO(r)
	return &dto, nil
}
