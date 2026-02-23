package pet

import "time"

// PetStatus represents the current status of a pet.
type PetStatus string

const (
	PetStatusRegistered   PetStatus = "registered"
	PetStatusAdopted      PetStatus = "adopted"
	PetStatusOwnerFound   PetStatus = "owner_found"
	PetStatusNeedsShelter PetStatus = "needs_shelter"
	PetStatusUnknown      PetStatus = "unknown"
)

// Pet represents an animal associated with a person in a project.
type Pet struct {
	ID             string
	ProjectID      string
	OwnerID        *string
	Name           string
	Status         PetStatus
	RegistrationID *string
	Notes          *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
