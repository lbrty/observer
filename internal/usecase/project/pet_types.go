package project

import (
	"time"

	"github.com/lbrty/observer/internal/domain/pet"
)

// PetDTO is the project-scoped pet representation.
type PetDTO struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"project_id"`
	OwnerID        *string   `json:"owner_id,omitempty"`
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	RegistrationID *string   `json:"registration_id,omitempty"`
	Notes          *string   `json:"notes,omitempty"`
	TagIDs         []string  `json:"tag_ids"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreatePetInput holds data for creating a pet.
type CreatePetInput struct {
	OwnerID        *string `json:"owner_id"`
	Name           string  `json:"name" binding:"required"`
	Status         *string `json:"status"`
	RegistrationID *string `json:"registration_id"`
	Notes          *string `json:"notes"`
}

// UpdatePetInput holds data for updating a pet.
type UpdatePetInput struct {
	OwnerID        *string `json:"owner_id"`
	Name           *string `json:"name"`
	Status         *string `json:"status"`
	RegistrationID *string `json:"registration_id"`
	Notes          *string `json:"notes"`
}

// ListPetsInput holds filter parameters.
type ListPetsInput struct {
	Page    int      `form:"page"`
	PerPage int      `form:"per_page"`
	Status  string   `form:"status"`
	TagIDs  []string `form:"tag_ids"`
}

// ListPetsOutput holds paginated results.
type ListPetsOutput struct {
	Pets    []PetDTO `json:"pets"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
}

func petToDTO(p *pet.Pet) PetDTO {
	return PetDTO{
		ID:             p.ID,
		ProjectID:      p.ProjectID,
		OwnerID:        p.OwnerID,
		Name:           p.Name,
		Status:         string(p.Status),
		RegistrationID: p.RegistrationID,
		Notes:          p.Notes,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}
