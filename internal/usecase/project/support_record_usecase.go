package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/support"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
	"github.com/lbrty/observer/internal/usecase"
)

// SupportRecordUseCase handles support record operations within a project.
type SupportRecordUseCase struct {
	repo repository.SupportRecordRepository
}

// NewSupportRecordUseCase creates a SupportRecordUseCase.
func NewSupportRecordUseCase(repo repository.SupportRecordRepository) *SupportRecordUseCase {
	return &SupportRecordUseCase{repo: repo}
}

// List returns paginated support records.
func (uc *SupportRecordUseCase) List(ctx context.Context, projectID string, input ListSupportRecordsInput) (*ListSupportRecordsOutput, error) {
	filter := support.RecordListFilter{
		ProjectID:    projectID,
		PersonID:     input.PersonID,
		ConsultantID: input.ConsultantID,
		OfficeID:     input.OfficeID,
		Page:         input.Page,
		PerPage:      input.PerPage,
	}
	if input.Type != nil {
		t := support.SupportType(*input.Type)
		filter.Type = &t
	}
	if input.Sphere != nil {
		s := support.SupportSphere(*input.Sphere)
		filter.Sphere = &s
	}

	records, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list support records: %w", err)
	}

	dtos := make([]SupportRecordDTO, len(records))
	for i, r := range records {
		dtos[i] = supportRecordToDTO(r)
	}

	page, perPage := usecase.ClampPagination(input.Page, input.PerPage)

	return &ListSupportRecordsOutput{
		Records: dtos,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}

// Get returns a support record by ID.
func (uc *SupportRecordUseCase) Get(ctx context.Context, id string) (*SupportRecordDTO, error) {
	r, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get support record: %w", err)
	}
	dto := supportRecordToDTO(r)
	return &dto, nil
}

// Create creates a new support record. recordedBy is auto-set from auth context.
func (uc *SupportRecordUseCase) Create(ctx context.Context, projectID string, recordedBy string, input CreateSupportRecordInput) (*SupportRecordDTO, error) {
	r := &support.Record{
		ID:               ulid.NewString(),
		PersonID:         input.PersonID,
		ProjectID:        projectID,
		ConsultantID:     input.ConsultantID,
		RecordedBy:       &recordedBy,
		OfficeID:         input.OfficeID,
		ReferredToOffice: input.ReferredToOffice,
		Type:             support.SupportType(input.Type),
		Notes:            input.Notes,
	}

	if input.Sphere != nil {
		s := support.SupportSphere(*input.Sphere)
		r.Sphere = &s
	}
	if input.ReferralStatus != nil {
		s := support.ReferralStatus(*input.ReferralStatus)
		r.ReferralStatus = &s
	}
	if err := parseDateField(input.ProvidedAt, &r.ProvidedAt); err != nil {
		return nil, fmt.Errorf("invalid provided_at: %w", err)
	}

	if err := uc.repo.Create(ctx, r); err != nil {
		return nil, fmt.Errorf("create support record: %w", err)
	}
	dto := supportRecordToDTO(r)
	return &dto, nil
}

// Update applies a partial update to a support record.
func (uc *SupportRecordUseCase) Update(ctx context.Context, id string, input UpdateSupportRecordInput) (*SupportRecordDTO, error) {
	r, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get support record for update: %w", err)
	}

	if input.ConsultantID != nil {
		r.ConsultantID = input.ConsultantID
	}
	if input.OfficeID != nil {
		r.OfficeID = input.OfficeID
	}
	if input.ReferredToOffice != nil {
		r.ReferredToOffice = input.ReferredToOffice
	}
	if input.Type != nil {
		r.Type = support.SupportType(*input.Type)
	}
	if input.Sphere != nil {
		s := support.SupportSphere(*input.Sphere)
		r.Sphere = &s
	}
	if input.ReferralStatus != nil {
		s := support.ReferralStatus(*input.ReferralStatus)
		r.ReferralStatus = &s
	}
	if input.Notes != nil {
		r.Notes = input.Notes
	}
	if err := parseDateField(input.ProvidedAt, &r.ProvidedAt); err != nil {
		return nil, fmt.Errorf("invalid provided_at: %w", err)
	}

	if err := uc.repo.Update(ctx, r); err != nil {
		return nil, fmt.Errorf("update support record: %w", err)
	}
	dto := supportRecordToDTO(r)
	return &dto, nil
}

// Delete removes a support record.
func (uc *SupportRecordUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete support record: %w", err)
	}
	return nil
}
