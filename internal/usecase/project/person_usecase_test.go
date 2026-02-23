package project_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/person"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func ptr[T any](v T) *T { return &v }

func TestPersonUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonUseCase(mockRepo)

	mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*person.Person{
		{ID: "p1", ProjectID: "proj1", FirstName: "Aida", Sex: person.SexFemale, CaseStatus: person.CaseStatusActive, PhoneNumbers: json.RawMessage("[]")},
		{ID: "p2", ProjectID: "proj1", FirstName: "Bek", Sex: person.SexMale, CaseStatus: person.CaseStatusNew, PhoneNumbers: json.RawMessage("[]")},
	}, 2, nil)

	out, err := uc.List(context.Background(), "proj1", ucproject.ListPeopleInput{Page: 1, PerPage: 20}, true, true)
	require.NoError(t, err)
	assert.Len(t, out.People, 2)
	assert.Equal(t, 2, out.Total)
	assert.Equal(t, "Aida", out.People[0].FirstName)
}

func TestPersonUseCase_Get_WithRedaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonUseCase(mockRepo)

	email := "aida@example.com"
	phone := "+996555111222"
	lastName := "Akmatova"
	bd := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)

	mockRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&person.Person{
		ID:           "p1",
		ProjectID:    "proj1",
		FirstName:    "Aida",
		LastName:     &lastName,
		Email:        &email,
		PrimaryPhone: &phone,
		BirthDate:    &bd,
		Sex:          person.SexFemale,
		CaseStatus:   person.CaseStatusActive,
		PhoneNumbers: json.RawMessage(`["+996555333444"]`),
	}, nil)

	// Without contact or personal visibility
	out, err := uc.Get(context.Background(), "p1", false, false)
	require.NoError(t, err)
	assert.Equal(t, "Aida", out.FirstName)
	assert.Nil(t, out.Email, "email should be redacted")
	assert.Nil(t, out.PrimaryPhone, "phone should be redacted")
	assert.Empty(t, out.PhoneNumbers, "phone_numbers should be empty")
	assert.Nil(t, out.LastName, "last_name should be redacted")
	assert.Nil(t, out.BirthDate, "birth_date should be redacted")
}

func TestPersonUseCase_Get_FullVisibility(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonUseCase(mockRepo)

	email := "aida@example.com"
	phone := "+996555111222"
	lastName := "Akmatova"

	mockRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&person.Person{
		ID:           "p1",
		ProjectID:    "proj1",
		FirstName:    "Aida",
		LastName:     &lastName,
		Email:        &email,
		PrimaryPhone: &phone,
		Sex:          person.SexFemale,
		CaseStatus:   person.CaseStatusActive,
		PhoneNumbers: json.RawMessage(`["+996555333444"]`),
	}, nil)

	out, err := uc.Get(context.Background(), "p1", true, true)
	require.NoError(t, err)
	assert.Equal(t, &email, out.Email)
	assert.Equal(t, &phone, out.PrimaryPhone)
	assert.Equal(t, &lastName, out.LastName)
	assert.Len(t, out.PhoneNumbers, 1)
}

func TestPersonUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, p *person.Person) error {
			assert.NotEmpty(t, p.ID)
			assert.Equal(t, "proj1", p.ProjectID)
			assert.Equal(t, "Aida", p.FirstName)
			assert.Equal(t, person.SexFemale, p.Sex)
			assert.Equal(t, person.CaseStatusNew, p.CaseStatus)
			return nil
		})

	out, err := uc.Create(context.Background(), "proj1", ucproject.CreatePersonInput{
		FirstName: "Aida",
		Sex:       ptr("female"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Aida", out.FirstName)
	assert.Equal(t, "female", out.Sex)
}

func TestPersonUseCase_Create_AgeConstraint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonUseCase(mockRepo)

	// Both birth_date and age_group should fail validation
	_, err := uc.Create(context.Background(), "proj1", ucproject.CreatePersonInput{
		FirstName: "Aida",
		BirthDate: ptr("1990-01-15"),
		AgeGroup:  ptr("young_adult"),
	})
	assert.ErrorIs(t, err, person.ErrAgeConstraint)
}

func TestPersonUseCase_Create_ConsentConstraint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonUseCase(mockRepo)

	// consent_date without consent_given=true should fail
	_, err := uc.Create(context.Background(), "proj1", ucproject.CreatePersonInput{
		FirstName:   "Aida",
		ConsentDate: ptr("2024-01-15"),
	})
	assert.ErrorIs(t, err, person.ErrConsentConstraint)
}

func TestPersonUseCase_Create_ExternalIDConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(person.ErrExternalIDExists)

	_, err := uc.Create(context.Background(), "proj1", ucproject.CreatePersonInput{
		FirstName:  "Aida",
		ExternalID: ptr("EXT-001"),
	})
	assert.ErrorIs(t, err, person.ErrExternalIDExists)
}

func TestPersonUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "p1").Return(nil)

	err := uc.Delete(context.Background(), "p1")
	require.NoError(t, err)
}

func TestPersonUseCase_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "nonexistent").Return(person.ErrPersonNotFound)

	err := uc.Delete(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, person.ErrPersonNotFound)
}
