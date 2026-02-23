package project

import (
	"time"

	"github.com/lbrty/observer/internal/domain/household"
)

// HouseholdDTO is the project-scoped household representation.
type HouseholdDTO struct {
	ID              string               `json:"id"`
	ProjectID       string               `json:"project_id"`
	ReferenceNumber *string              `json:"reference_number,omitempty"`
	HeadPersonID    *string              `json:"head_person_id,omitempty"`
	Members         []HouseholdMemberDTO `json:"members,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

// HouseholdMemberDTO represents a household member.
type HouseholdMemberDTO struct {
	PersonID     string `json:"person_id"`
	Relationship string `json:"relationship"`
}

// CreateHouseholdInput holds data for creating a household.
type CreateHouseholdInput struct {
	ReferenceNumber *string `json:"reference_number"`
	HeadPersonID    *string `json:"head_person_id"`
}

// UpdateHouseholdInput holds data for updating a household.
type UpdateHouseholdInput struct {
	ReferenceNumber *string `json:"reference_number"`
	HeadPersonID    *string `json:"head_person_id"`
}

// AddMemberInput holds data for adding a member.
type AddMemberInput struct {
	PersonID     string `json:"person_id" binding:"required"`
	Relationship string `json:"relationship" binding:"required"`
}

// ListHouseholdsInput holds filter parameters.
type ListHouseholdsInput struct {
	Page    int `form:"page"`
	PerPage int `form:"per_page"`
}

// ListHouseholdsOutput holds paginated results.
type ListHouseholdsOutput struct {
	Households []HouseholdDTO `json:"households"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
}

func householdToDTO(h *household.Household) HouseholdDTO {
	return HouseholdDTO{
		ID:              h.ID,
		ProjectID:       h.ProjectID,
		ReferenceNumber: h.ReferenceNumber,
		HeadPersonID:    h.HeadPersonID,
		CreatedAt:       h.CreatedAt,
		UpdatedAt:       h.UpdatedAt,
	}
}

func memberToDTO(m *household.Member) HouseholdMemberDTO {
	return HouseholdMemberDTO{
		PersonID:     m.PersonID,
		Relationship: string(m.Relationship),
	}
}
