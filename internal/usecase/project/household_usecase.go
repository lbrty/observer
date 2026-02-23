package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/household"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// HouseholdUseCase handles household operations within a project.
type HouseholdUseCase struct {
	repo       repository.HouseholdRepository
	memberRepo repository.HouseholdMemberRepository
}

// NewHouseholdUseCase creates a HouseholdUseCase.
func NewHouseholdUseCase(repo repository.HouseholdRepository, memberRepo repository.HouseholdMemberRepository) *HouseholdUseCase {
	return &HouseholdUseCase{repo: repo, memberRepo: memberRepo}
}

// List returns paginated households with members.
func (uc *HouseholdUseCase) List(ctx context.Context, projectID string, input ListHouseholdsInput) (*ListHouseholdsOutput, error) {
	page := input.Page
	if page < 1 {
		page = 1
	}
	perPage := input.PerPage
	if perPage < 1 {
		perPage = 20
	}

	households, total, err := uc.repo.List(ctx, projectID, page, perPage)
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
func (uc *HouseholdUseCase) Get(ctx context.Context, id string) (*HouseholdDTO, error) {
	h, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get household: %w", err)
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
func (uc *HouseholdUseCase) Create(ctx context.Context, projectID string, input CreateHouseholdInput) (*HouseholdDTO, error) {
	h := &household.Household{
		ID:              ulid.NewString(),
		ProjectID:       projectID,
		ReferenceNumber: input.ReferenceNumber,
		HeadPersonID:    input.HeadPersonID,
	}
	if err := uc.repo.Create(ctx, h); err != nil {
		return nil, fmt.Errorf("create household: %w", err)
	}
	dto := householdToDTO(h)
	dto.Members = []HouseholdMemberDTO{}
	return &dto, nil
}

// Update applies a partial update to a household.
func (uc *HouseholdUseCase) Update(ctx context.Context, id string, input UpdateHouseholdInput) (*HouseholdDTO, error) {
	h, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get household for update: %w", err)
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
	dto := householdToDTO(h)
	return &dto, nil
}

// Delete removes a household.
func (uc *HouseholdUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete household: %w", err)
	}
	return nil
}

// AddMember adds a member to a household.
func (uc *HouseholdUseCase) AddMember(ctx context.Context, householdID string, input AddMemberInput) (*HouseholdMemberDTO, error) {
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
func (uc *HouseholdUseCase) RemoveMember(ctx context.Context, householdID, personID string) error {
	if err := uc.memberRepo.Remove(ctx, householdID, personID); err != nil {
		return fmt.Errorf("remove household member: %w", err)
	}
	return nil
}
