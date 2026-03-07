package project

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lbrty/observer/internal/domain/person"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
	"github.com/lbrty/observer/internal/usecase"
)

// PersonUseCase handles person operations within a project.
type PersonUseCase struct {
	repo    repository.PersonRepository
	tagRepo repository.PersonTagRepository
}

// NewPersonUseCase creates a PersonUseCase.
func NewPersonUseCase(repo repository.PersonRepository, tagRepo repository.PersonTagRepository) *PersonUseCase {
	return &PersonUseCase{repo: repo, tagRepo: tagRepo}
}

// List returns paginated people with sensitivity-aware redaction.
func (uc *PersonUseCase) List(ctx context.Context, projectID string, input ListPeopleInput, canViewContact, canViewPersonal bool) (*ListPeopleOutput, error) {
	filter := person.PersonListFilter{
		ProjectID:    projectID,
		ConsultantID: input.ConsultantID,
		OfficeID:     input.OfficeID,
		Search:       input.Search,
		TagIDs:       input.TagIDs,
		Page:         input.Page,
		PerPage:      input.PerPage,
	}
	if input.CaseStatus != nil {
		s := person.CaseStatus(*input.CaseStatus)
		filter.CaseStatus = &s
	}

	people, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list people: %w", err)
	}

	ids := make([]string, len(people))
	dtos := make([]PersonDTO, len(people))
	for i, p := range people {
		ids[i] = p.ID
		dtos[i] = personToDTO(p, canViewContact, canViewPersonal)
	}

	tagMap, err := uc.tagRepo.ListBulk(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("list person tags: %w", err)
	}
	for i := range dtos {
		if tags, ok := tagMap[dtos[i].ID]; ok {
			dtos[i].TagIDs = tags
		} else {
			dtos[i].TagIDs = []string{}
		}
	}

	page, perPage := usecase.ClampPagination(input.Page, input.PerPage)

	return &ListPeopleOutput{
		People:  dtos,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}

// Get returns a person by ID with sensitivity-aware redaction.
func (uc *PersonUseCase) Get(ctx context.Context, id string, canViewContact, canViewPersonal bool) (*PersonDTO, error) {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get person: %w", err)
	}
	dto := personToDTO(p, canViewContact, canViewPersonal)
	return &dto, nil
}

// Create creates a new person in the project.
func (uc *PersonUseCase) Create(ctx context.Context, projectID string, input CreatePersonInput) (*PersonDTO, error) {
	p := &person.Person{
		ID:             ulid.NewString(),
		ProjectID:      projectID,
		ConsultantID:   input.ConsultantID,
		OfficeID:       input.OfficeID,
		CurrentPlaceID: input.CurrentPlaceID,
		OriginPlaceID:  input.OriginPlaceID,
		ExternalID:     input.ExternalID,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Patronymic:     input.Patronymic,
		Email:          input.Email,
		Sex:            person.SexUnknown,
		CaseStatus:     person.CaseStatusNew,
		PhoneNumbers:   json.RawMessage("[]"),
	}

	if input.Sex != nil {
		p.Sex = person.Sex(*input.Sex)
	}
	if input.CaseStatus != nil {
		p.CaseStatus = person.CaseStatus(*input.CaseStatus)
	}
	if input.AgeGroup != nil {
		ag := person.AgeGroup(*input.AgeGroup)
		p.AgeGroup = &ag
	}
	if input.PrimaryPhone != nil {
		p.PrimaryPhone = input.PrimaryPhone
	}
	if input.PhoneNumbers != nil {
		b, _ := json.Marshal(input.PhoneNumbers)
		p.PhoneNumbers = b
	}
	if input.ConsentGiven != nil {
		p.ConsentGiven = *input.ConsentGiven
	}

	if err := parseDateField(input.BirthDate, &p.BirthDate); err != nil {
		return nil, fmt.Errorf("invalid birth_date: %w", err)
	}
	if err := parseDateField(input.ConsentDate, &p.ConsentDate); err != nil {
		return nil, fmt.Errorf("invalid consent_date: %w", err)
	}
	if err := parseDateField(input.RegisteredAt, &p.RegisteredAt); err != nil {
		return nil, fmt.Errorf("invalid registered_at: %w", err)
	}

	if err := validatePersonConstraints(p); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create person: %w", err)
	}
	dto := personToDTO(p, true, true)
	return &dto, nil
}

// Update applies a partial update to a person.
func (uc *PersonUseCase) Update(ctx context.Context, id string, input UpdatePersonInput) (*PersonDTO, error) {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get person for update: %w", err)
	}

	if input.ConsultantID != nil {
		p.ConsultantID = input.ConsultantID
	}
	if input.OfficeID != nil {
		p.OfficeID = input.OfficeID
	}
	if input.CurrentPlaceID != nil {
		p.CurrentPlaceID = input.CurrentPlaceID
	}
	if input.OriginPlaceID != nil {
		p.OriginPlaceID = input.OriginPlaceID
	}
	if input.ExternalID != nil {
		p.ExternalID = input.ExternalID
	}
	if input.FirstName != nil {
		p.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		p.LastName = input.LastName
	}
	if input.Patronymic != nil {
		p.Patronymic = input.Patronymic
	}
	if input.Email != nil {
		p.Email = input.Email
	}
	if input.Sex != nil {
		p.Sex = person.Sex(*input.Sex)
	}
	if input.AgeGroup != nil {
		ag := person.AgeGroup(*input.AgeGroup)
		p.AgeGroup = &ag
	}
	if input.PrimaryPhone != nil {
		p.PrimaryPhone = input.PrimaryPhone
	}
	if input.PhoneNumbers != nil {
		b, _ := json.Marshal(input.PhoneNumbers)
		p.PhoneNumbers = b
	}
	if input.CaseStatus != nil {
		p.CaseStatus = person.CaseStatus(*input.CaseStatus)
	}
	if input.ConsentGiven != nil {
		p.ConsentGiven = *input.ConsentGiven
	}

	if err := parseDateField(input.BirthDate, &p.BirthDate); err != nil {
		return nil, fmt.Errorf("invalid birth_date: %w", err)
	}
	if err := parseDateField(input.ConsentDate, &p.ConsentDate); err != nil {
		return nil, fmt.Errorf("invalid consent_date: %w", err)
	}
	if err := parseDateField(input.RegisteredAt, &p.RegisteredAt); err != nil {
		return nil, fmt.Errorf("invalid registered_at: %w", err)
	}

	if err := validatePersonConstraints(p); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("update person: %w", err)
	}
	dto := personToDTO(p, true, true)
	return &dto, nil
}

// Delete removes a person.
func (uc *PersonUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete person: %w", err)
	}
	return nil
}

func validatePersonConstraints(p *person.Person) error {
	if p.BirthDate != nil && p.AgeGroup != nil {
		return person.ErrAgeConstraint
	}
	if !p.ConsentGiven && p.ConsentDate != nil {
		return person.ErrConsentConstraint
	}
	return nil
}

func parseDateField(s *string, target **time.Time) error {
	if s == nil {
		return nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return err
	}
	*target = &t
	return nil
}

func parsePhoneNumbers(raw json.RawMessage) []string {
	if raw == nil {
		return []string{}
	}
	var phones []string
	if err := json.Unmarshal(raw, &phones); err != nil {
		return []string{}
	}
	return phones
}
