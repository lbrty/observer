package project

import (
	"time"

	"github.com/lbrty/observer/internal/domain/support"
)

// SupportRecordDTO is the project-scoped support record representation.
type SupportRecordDTO struct {
	ID               string    `json:"id"`
	PersonID         string    `json:"person_id"`
	ProjectID        string    `json:"project_id"`
	ConsultantID     *string   `json:"consultant_id,omitempty"`
	RecordedBy       *string   `json:"recorded_by,omitempty"`
	OfficeID         *string   `json:"office_id,omitempty"`
	ReferredToOffice *string   `json:"referred_to_office,omitempty"`
	Type             string    `json:"type"`
	Sphere           *string   `json:"sphere,omitempty"`
	ReferralStatus   *string   `json:"referral_status,omitempty"`
	ProvidedAt       *string   `json:"provided_at,omitempty"`
	Notes            *string   `json:"notes,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CreateSupportRecordInput holds data for creating a support record.
type CreateSupportRecordInput struct {
	PersonID         string  `json:"person_id" binding:"required"`
	ConsultantID     *string `json:"consultant_id"`
	OfficeID         *string `json:"office_id"`
	ReferredToOffice *string `json:"referred_to_office"`
	Type             string  `json:"type" binding:"required"`
	Sphere           *string `json:"sphere"`
	ReferralStatus   *string `json:"referral_status"`
	ProvidedAt       *string `json:"provided_at"`
	Notes            *string `json:"notes"`
}

// UpdateSupportRecordInput holds data for updating a support record.
type UpdateSupportRecordInput struct {
	ConsultantID     *string `json:"consultant_id"`
	OfficeID         *string `json:"office_id"`
	ReferredToOffice *string `json:"referred_to_office"`
	Type             *string `json:"type"`
	Sphere           *string `json:"sphere"`
	ReferralStatus   *string `json:"referral_status"`
	ProvidedAt       *string `json:"provided_at"`
	Notes            *string `json:"notes"`
}

// ListSupportRecordsInput holds filter parameters.
type ListSupportRecordsInput struct {
	PersonID       *string `form:"person_id"`
	ConsultantID   *string `form:"consultant_id"`
	OfficeID       *string `form:"office_id"`
	Type           *string `form:"type"`
	Sphere         *string `form:"sphere"`
	ReferralStatus *string `form:"referral_status"`
	DateFrom       *string `form:"date_from"`
	DateTo         *string `form:"date_to"`
	Page           int     `form:"page"`
	PerPage        int     `form:"per_page"`
}

// ListSupportRecordsOutput holds paginated results.
type ListSupportRecordsOutput struct {
	Records []SupportRecordDTO `json:"records"`
	Total   int                `json:"total"`
	Page    int                `json:"page"`
	PerPage int                `json:"per_page"`
}

func supportRecordToDTO(r *support.Record) SupportRecordDTO {
	dto := SupportRecordDTO{
		ID:               r.ID,
		PersonID:         r.PersonID,
		ProjectID:        r.ProjectID,
		ConsultantID:     r.ConsultantID,
		RecordedBy:       r.RecordedBy,
		OfficeID:         r.OfficeID,
		ReferredToOffice: r.ReferredToOffice,
		Type:             string(r.Type),
		Notes:            r.Notes,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
	if r.Sphere != nil {
		s := string(*r.Sphere)
		dto.Sphere = &s
	}
	if r.ReferralStatus != nil {
		s := string(*r.ReferralStatus)
		dto.ReferralStatus = &s
	}
	if r.ProvidedAt != nil {
		s := r.ProvidedAt.Format("2006-01-02")
		dto.ProvidedAt = &s
	}
	return dto
}
