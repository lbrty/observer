package migration

import "time"

// MovementReason describes why a person moved.
type MovementReason string

const (
	ReasonConflict          MovementReason = "conflict"
	ReasonSecurity          MovementReason = "security"
	ReasonServiceAccess     MovementReason = "service_access"
	ReasonReturn            MovementReason = "return"
	ReasonRelocationProgram MovementReason = "relocation_program"
	ReasonEconomic          MovementReason = "economic"
	ReasonOther             MovementReason = "other"
)

// HousingAtDestination describes housing situation after migration.
type HousingAtDestination string

const (
	HousingOwnProperty    HousingAtDestination = "own_property"
	HousingRenting        HousingAtDestination = "renting"
	HousingWithRelatives  HousingAtDestination = "with_relatives"
	HousingCollectiveSite HousingAtDestination = "collective_site"
	HousingHotel          HousingAtDestination = "hotel"
	HousingOther          HousingAtDestination = "other"
	HousingUnknown        HousingAtDestination = "unknown"
)

// Record represents a single migration event for a person.
type Record struct {
	ID                   string
	PersonID             string
	FromPlaceID          *string
	DestinationPlaceID   *string
	MigrationDate        *time.Time
	MovementReason       *MovementReason
	HousingAtDestination *HousingAtDestination
	Notes                *string
	CreatedAt            time.Time
}
