package project

import (
	"time"

	"github.com/lbrty/observer/internal/domain/person"
)

// PersonDTO is the project-scoped person representation.
type PersonDTO struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"project_id"`
	ConsultantID   *string   `json:"consultant_id,omitempty"`
	OfficeID       *string   `json:"office_id,omitempty"`
	CurrentPlaceID *string   `json:"current_place_id,omitempty"`
	OriginPlaceID  *string   `json:"origin_place_id,omitempty"`
	ExternalID     *string   `json:"external_id,omitempty"`
	FirstName      string    `json:"first_name"`
	LastName       *string   `json:"last_name,omitempty"`
	Patronymic     *string   `json:"patronymic,omitempty"`
	Email          *string   `json:"email,omitempty"`
	BirthDate      *string   `json:"birth_date,omitempty"`
	Sex            string    `json:"sex"`
	AgeGroup       *string   `json:"age_group,omitempty"`
	PrimaryPhone   *string   `json:"primary_phone,omitempty"`
	PhoneNumbers   []string  `json:"phone_numbers"`
	CaseStatus     string    `json:"case_status"`
	ConsentGiven   bool      `json:"consent_given"`
	ConsentDate    *string   `json:"consent_date,omitempty"`
	RegisteredAt   *string   `json:"registered_at,omitempty"`
	TagIDs         []string  `json:"tag_ids"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreatePersonInput holds data for creating a person.
type CreatePersonInput struct {
	ConsultantID   *string  `json:"consultant_id"`
	OfficeID       *string  `json:"office_id"`
	CurrentPlaceID *string  `json:"current_place_id"`
	OriginPlaceID  *string  `json:"origin_place_id"`
	ExternalID     *string  `json:"external_id"`
	FirstName      string   `json:"first_name" binding:"required"`
	LastName       *string  `json:"last_name"`
	Patronymic     *string  `json:"patronymic"`
	Email          *string  `json:"email"`
	BirthDate      *string  `json:"birth_date"`
	Sex            *string  `json:"sex"`
	AgeGroup       *string  `json:"age_group"`
	PrimaryPhone   *string  `json:"primary_phone"`
	PhoneNumbers   []string `json:"phone_numbers"`
	CaseStatus     *string  `json:"case_status"`
	ConsentGiven   *bool    `json:"consent_given"`
	ConsentDate    *string  `json:"consent_date"`
	RegisteredAt   *string  `json:"registered_at"`
}

// UpdatePersonInput holds data for updating a person.
type UpdatePersonInput struct {
	ConsultantID   *string  `json:"consultant_id"`
	OfficeID       *string  `json:"office_id"`
	CurrentPlaceID *string  `json:"current_place_id"`
	OriginPlaceID  *string  `json:"origin_place_id"`
	ExternalID     *string  `json:"external_id"`
	FirstName      *string  `json:"first_name"`
	LastName       *string  `json:"last_name"`
	Patronymic     *string  `json:"patronymic"`
	Email          *string  `json:"email"`
	BirthDate      *string  `json:"birth_date"`
	Sex            *string  `json:"sex"`
	AgeGroup       *string  `json:"age_group"`
	PrimaryPhone   *string  `json:"primary_phone"`
	PhoneNumbers   []string `json:"phone_numbers"`
	CaseStatus     *string  `json:"case_status"`
	ConsentGiven   *bool    `json:"consent_given"`
	ConsentDate    *string  `json:"consent_date"`
	RegisteredAt   *string  `json:"registered_at"`
}

// ListPeopleInput holds filter parameters for listing people.
type ListPeopleInput struct {
	ConsultantID *string  `form:"consultant_id"`
	OfficeID     *string  `form:"office_id"`
	CaseStatus   *string  `form:"case_status"`
	Sex          *string  `form:"sex"`
	AgeGroup     *string  `form:"age_group"`
	CategoryID   *string  `form:"category_id"`
	RegionID     *string  `form:"region_id"`
	HasPets      *bool    `form:"has_pets"`
	Search       *string  `form:"search"`
	TagIDs       []string `form:"tag_ids"`
	Page         int      `form:"page"`
	PerPage      int      `form:"per_page"`
}

// ListPeopleOutput holds paginated results.
type ListPeopleOutput struct {
	People  []PersonDTO `json:"people"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
}

// ReplaceIDsInput holds a set of IDs for bulk replace operations (categories/tags).
type ReplaceIDsInput struct {
	IDs []string `json:"ids" binding:"required"`
}

func personToDTO(p *person.Person, canViewContact, canViewPersonal bool) PersonDTO {
	dto := PersonDTO{
		ID:           p.ID,
		ProjectID:    p.ProjectID,
		ConsultantID: p.ConsultantID,
		OfficeID:     p.OfficeID,
		ExternalID:   p.ExternalID,
		FirstName:    p.FirstName,
		Sex:          string(p.Sex),
		CaseStatus:   string(p.CaseStatus),
		ConsentGiven: p.ConsentGiven,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}

	if canViewPersonal {
		dto.LastName = p.LastName
		dto.Patronymic = p.Patronymic
		dto.CurrentPlaceID = p.CurrentPlaceID
		dto.OriginPlaceID = p.OriginPlaceID
		if p.BirthDate != nil {
			s := p.BirthDate.Format("2006-01-02")
			dto.BirthDate = &s
		}
		if p.AgeGroup != nil {
			s := string(*p.AgeGroup)
			dto.AgeGroup = &s
		}
		if p.ConsentDate != nil {
			s := p.ConsentDate.Format("2006-01-02")
			dto.ConsentDate = &s
		}
		if p.RegisteredAt != nil {
			s := p.RegisteredAt.Format("2006-01-02")
			dto.RegisteredAt = &s
		}
	}

	if canViewContact {
		dto.Email = p.Email
		dto.PrimaryPhone = p.PrimaryPhone
		dto.PhoneNumbers = parsePhoneNumbers(p.PhoneNumbers)
	} else {
		dto.PhoneNumbers = []string{}
	}

	return dto
}
