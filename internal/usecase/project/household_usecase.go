package project

import (
	"context"
	"fmt"
	"time"

	"github.com/lbrty/observer/internal/domain/household"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
	"github.com/lbrty/observer/internal/usecase"
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
)

// HouseholdUseCase handles household operations within a project.
type HouseholdUseCase struct {
	repo       repository.HouseholdRepository
	memberRepo repository.HouseholdMemberRepository
	auditUC    *ucaudit.AuditUseCase
}

// NewHouseholdUseCase creates a HouseholdUseCase.
func NewHouseholdUseCase(
	repo repository.HouseholdRepository,
	memberRepo repository.HouseholdMemberRepository,
	auditUC *ucaudit.AuditUseCase,
) *HouseholdUseCase {
	return &HouseholdUseCase{repo: repo, memberRepo: memberRepo, auditUC: auditUC}
}

// List returns paginated households with members.
func (uc *HouseholdUseCase) List(
	ctx context.Context, projectID string, input ListHouseholdsInput,
) (*ListHouseholdsOutput, error) {
	page, perPage := usecase.ClampPagination(input.Page, input.PerPage)

	filter := household.HouseholdListFilter{
		ProjectID: projectID,
		Page:      page,
		PerPage:   perPage,
	}

	if input.Search != "" {
		filter.Search = &input.Search
	}

	if input.CreatedFrom != "" {
		if err := parseDateField(&input.CreatedFrom, &filter.CreatedFrom); err != nil {
			return nil, fmt.Errorf("invalid created_from: %w", err)
		}
	}

	if input.CreatedTo != "" {
		var t *time.Time
		if err := parseDateField(&input.CreatedTo, &t); err != nil {
			return nil, fmt.Errorf("invalid created_to: %w", err)
		}
		end := t.Add(24*time.Hour - time.Nanosecond)
		filter.CreatedTo = &end
	}

	households, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list households: %w", err)
	}

	dtos := make([]HouseholdDTO, len(households))
	for i, h := range households {
		dtos[i] = householdToDTO(h)
	}

	return &ListHouseholdsOutput{
		Households: dtos,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
	}, nil
}

// Get returns a household by ID with its members.
func (uc *HouseholdUseCase) Get(ctx context.Context, projectID, id string) (*HouseholdDTO, error) {
	h, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get household: %w", err)
	}

	if h.ProjectID != projectID {
		return nil, household.ErrHouseholdNotFound
	}

	dto := householdToDTO(h)
	members, err := uc.memberRepo.List(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("list household members: %w", err)
	}

	dto.Members = make([]HouseholdMemberDTO, len(members))
	for i, m := range members {
		dto.Members[i] = memberToDTO(m)
	}

	return &dto, nil
}

// Create creates a new household.
func (uc *HouseholdUseCase) Create(
	ctx context.Context,
	projectID string,
	input CreateHouseholdInput,
) (*HouseholdDTO, error) {
	h := &household.Household{
		ID:              ulid.NewString(),
		ProjectID:       projectID,
		ReferenceNumber: input.ReferenceNumber,
		HeadPersonID:    input.HeadPersonID,
	}

	if err := uc.repo.Create(ctx, h); err != nil {
		return nil, fmt.Errorf("create household: %w", err)
	}

	details := map[string]any{}
	if h.ReferenceNumber != nil {
		details["reference_number"] = *h.ReferenceNumber
	}
	uc.auditUC.Record(
		ctx,
		&projectID,
		"household.create",
		"household",
		&h.ID,
		fmt.Sprintf("Created household %s", h.ID),
		details,
	)

	dto := householdToDTO(h)
	dto.Members = []HouseholdMemberDTO{}
	return &dto, nil
}

// Update applies a partial update to a household.
func (uc *HouseholdUseCase) Update(
	ctx context.Context,
	projectID, id string,
	input UpdateHouseholdInput,
) (*HouseholdDTO, error) {
	h, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get household for update: %w", err)
	}

	if h.ProjectID != projectID {
		return nil, household.ErrHouseholdNotFound
	}

	if input.ReferenceNumber != nil {
		h.ReferenceNumber = input.ReferenceNumber
	}

	if input.HeadPersonID != nil {
		h.HeadPersonID = input.HeadPersonID
	}

	if err := uc.repo.Update(ctx, h); err != nil {
		return nil, fmt.Errorf("update household: %w", err)
	}

	details := map[string]any{}
	if h.ReferenceNumber != nil {
		details["reference_number"] = *h.ReferenceNumber
	}
	uc.auditUC.Record(
		ctx,
		&projectID,
		"household.update",
		"household",
		&id,
		fmt.Sprintf("Updated household %s", id),
		details,
	)

	dto := householdToDTO(h)
	return &dto, nil
}

// Delete removes a household.
func (uc *HouseholdUseCase) Delete(ctx context.Context, projectID, id string) error {
	h, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get household for delete: %w", err)
	}

	if h.ProjectID != projectID {
		return household.ErrHouseholdNotFound
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete household: %w", err)
	}

	details := map[string]any{}
	if h.ReferenceNumber != nil {
		details["reference_number"] = *h.ReferenceNumber
	}
	uc.auditUC.Record(
		ctx,
		&projectID,
		"household.delete",
		"household",
		&id,
		fmt.Sprintf("Deleted household %s", id),
		details,
	)

	return nil
}

// AddMember adds a member to a household.
func (uc *HouseholdUseCase) AddMember(
	ctx context.Context, projectID, householdID string, input AddMemberInput,
) (*HouseholdMemberDTO, error) {
	h, err := uc.repo.GetByID(ctx, householdID)
	if err != nil {
		return nil, fmt.Errorf("get household for add member: %w", err)
	}

	if h.ProjectID != projectID {
		return nil, household.ErrHouseholdNotFound
	}

	m := &household.Member{
		HouseholdID:  householdID,
		PersonID:     input.PersonID,
		Relationship: household.Relationship(input.Relationship),
	}

	if err := uc.memberRepo.Add(ctx, m); err != nil {
		return nil, fmt.Errorf("add household member: %w", err)
	}

	dto := memberToDTO(m)
	return &dto, nil
}

// RemoveMember removes a member from a household.
func (uc *HouseholdUseCase) RemoveMember(ctx context.Context, projectID, householdID, personID string) error {
	h, err := uc.repo.GetByID(ctx, householdID)
	if err != nil {
		return fmt.Errorf("get household for remove member: %w", err)
	}

	if h.ProjectID != projectID {
		return household.ErrHouseholdNotFound
	}

	if err := uc.memberRepo.Remove(ctx, householdID, personID); err != nil {
		return fmt.Errorf("remove household member: %w", err)
	}

	return nil
}
