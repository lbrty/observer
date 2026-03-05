package project

import (
	"time"

	"github.com/lbrty/observer/internal/domain/migration"
)

// MigrationRecordDTO is the project-scoped migration record representation.
type MigrationRecordDTO struct {
	ID                   string     `json:"id"`
	PersonID             string     `json:"person_id"`
	FromPlaceID          *string    `json:"from_place_id,omitempty"`
	DestinationPlaceID   *string    `json:"destination_place_id,omitempty"`
	MigrationDate        *string    `json:"migration_date,omitempty"`
	MovementReason       *string    `json:"movement_reason,omitempty"`
	HousingAtDestination *string    `json:"housing_at_destination,omitempty"`
	Notes                *string    `json:"notes,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty"`
}

// CreateMigrationRecordInput holds data for creating a migration record.
type CreateMigrationRecordInput struct {
	FromPlaceID          *string `json:"from_place_id"`
	DestinationPlaceID   *string `json:"destination_place_id"`
	MigrationDate        *string `json:"migration_date"`
	MovementReason       *string `json:"movement_reason"`
	HousingAtDestination *string `json:"housing_at_destination"`
	Notes                *string `json:"notes"`
}

// UpdateMigrationRecordInput holds data for updating a migration record.
type UpdateMigrationRecordInput struct {
	FromPlaceID          *string `json:"from_place_id"`
	DestinationPlaceID   *string `json:"destination_place_id"`
	MigrationDate        *string `json:"migration_date"`
	MovementReason       *string `json:"movement_reason"`
	HousingAtDestination *string `json:"housing_at_destination"`
	Notes                *string `json:"notes"`
}

func migrationRecordToDTO(r *migration.Record) MigrationRecordDTO {
	dto := MigrationRecordDTO{
		ID:                 r.ID,
		PersonID:           r.PersonID,
		FromPlaceID:        r.FromPlaceID,
		DestinationPlaceID: r.DestinationPlaceID,
		Notes:              r.Notes,
		CreatedAt:          r.CreatedAt,
	}
	if !r.UpdatedAt.IsZero() {
		t := r.UpdatedAt
		dto.UpdatedAt = &t
	}
	if r.MigrationDate != nil {
		s := r.MigrationDate.Format("2006-01-02")
		dto.MigrationDate = &s
	}
	if r.MovementReason != nil {
		s := string(*r.MovementReason)
		dto.MovementReason = &s
	}
	if r.HousingAtDestination != nil {
		s := string(*r.HousingAtDestination)
		dto.HousingAtDestination = &s
	}
	return dto
}
