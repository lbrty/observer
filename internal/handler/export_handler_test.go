package handler_test

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/household"
	"github.com/lbrty/observer/internal/domain/person"
	"github.com/lbrty/observer/internal/domain/pet"
	"github.com/lbrty/observer/internal/domain/support"
	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/middleware"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

type exportTestDeps struct {
	personRepo    *repomock.MockPersonRepository
	personTagRepo *repomock.MockPersonTagRepository
	supportRepo   *repomock.MockSupportRecordRepository
	petRepo       *repomock.MockPetRepository
	petTagRepo    *repomock.MockPetTagRepository
	householdRepo *repomock.MockHouseholdRepository
	memberRepo    *repomock.MockHouseholdMemberRepository
	auditRepo     *repomock.MockAuditLogRepository
	handler       *handler.ExportHandler
}

func newExportTestDeps(ctrl *gomock.Controller) *exportTestDeps {
	personRepo := repomock.NewMockPersonRepository(ctrl)
	personTagRepo := repomock.NewMockPersonTagRepository(ctrl)
	supportRepo := repomock.NewMockSupportRecordRepository(ctrl)
	petRepo := repomock.NewMockPetRepository(ctrl)
	petTagRepo := repomock.NewMockPetTagRepository(ctrl)
	householdRepo := repomock.NewMockHouseholdRepository(ctrl)
	memberRepo := repomock.NewMockHouseholdMemberRepository(ctrl)
	auditRepo := repomock.NewMockAuditLogRepository(ctrl)
	auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	personUC := ucproject.NewPersonUseCase(personRepo, personTagRepo, auditUC)
	supportUC := ucproject.NewSupportRecordUseCase(supportRepo, personRepo, auditUC)
	petUC := ucproject.NewPetUseCase(petRepo, petTagRepo, auditUC)
	householdUC := ucproject.NewHouseholdUseCase(householdRepo, memberRepo, auditUC)

	h := handler.NewExportHandler(personUC, supportUC, petUC, householdUC, auditUC)

	return &exportTestDeps{
		personRepo:    personRepo,
		personTagRepo: personTagRepo,
		supportRepo:   supportRepo,
		petRepo:       petRepo,
		petTagRepo:    petTagRepo,
		householdRepo: householdRepo,
		memberRepo:    memberRepo,
		auditRepo:     auditRepo,
		handler:       h,
	}
}

func parseCSV(body string) ([][]string, error) {
	return csv.NewReader(strings.NewReader(body)).ReadAll()
}

// ---- People ----

func TestExportHandler_ExportPeople_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newExportTestDeps(ctrl)

	projectID := testID().String()
	now := time.Now().UTC()

	deps.personRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*person.Person{
		{
			ID:         "person-1",
			ProjectID:  projectID,
			FirstName:  "Aida",
			Sex:        person.SexFemale,
			CaseStatus: person.CaseStatusActive,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}, 1, nil)
	deps.personTagRepo.EXPECT().ListBulk(gomock.Any(), []string{"person-1"}).Return(map[string][]string{}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/export/people", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	c.Set(string(middleware.CtxCanExport), true)
	c.Set(string(middleware.CtxCanViewContact), true)
	c.Set(string(middleware.CtxCanViewPersonal), true)
	deps.handler.ExportPeople(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), fmt.Sprintf("people-%s.csv", projectID))

	rows, err := parseCSV(w.Body.String())
	require.NoError(t, err)
	require.Len(t, rows, 2)
	assert.Equal(t, []string{"id", "first_name", "last_name", "patronymic", "email", "sex", "age_group", "case_status", "primary_phone", "registered_at", "created_at"}, rows[0])
	assert.Equal(t, "person-1", rows[1][0])
	assert.Equal(t, "Aida", rows[1][1])
	assert.Equal(t, "female", rows[1][5])
}

// ---- Support records ----

func TestExportHandler_ExportSupportRecords_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newExportTestDeps(ctrl)

	projectID := testID().String()
	personID := testID().String()
	now := time.Now().UTC()
	sphere := support.SphereFamilyLaw

	deps.supportRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*support.Record{
		{
			ID:        "sr-1",
			PersonID:  personID,
			ProjectID: projectID,
			Type:      support.SupportTypeLegal,
			Sphere:    &sphere,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, 1, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/export/support-records", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	c.Set(string(middleware.CtxCanExport), true)
	deps.handler.ExportSupportRecords(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), fmt.Sprintf("support-records-%s.csv", projectID))

	rows, err := parseCSV(w.Body.String())
	require.NoError(t, err)
	require.Len(t, rows, 2)
	assert.Equal(t, "sr-1", rows[1][0])
	assert.Equal(t, personID, rows[1][1])
}

// ---- Pets ----

func TestExportHandler_ExportPets_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newExportTestDeps(ctrl)

	projectID := testID().String()
	now := time.Now().UTC()

	deps.petRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*pet.Pet{
		{
			ID:        "pet-1",
			ProjectID: projectID,
			Name:      "Bars",
			Status:    pet.PetStatusRegistered,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, 1, nil)
	deps.petTagRepo.EXPECT().ListBulk(gomock.Any(), []string{"pet-1"}).Return(map[string][]string{}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/export/pets", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	c.Set(string(middleware.CtxCanExport), true)
	deps.handler.ExportPets(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), fmt.Sprintf("pets-%s.csv", projectID))

	rows, err := parseCSV(w.Body.String())
	require.NoError(t, err)
	require.Len(t, rows, 2)
	assert.Equal(t, "pet-1", rows[1][0])
	assert.Equal(t, "Bars", rows[1][1])
}

// ---- Households ----

func TestExportHandler_ExportHouseholds_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newExportTestDeps(ctrl)

	projectID := testID().String()
	now := time.Now().UTC()
	ref := "HH-001"
	headID := testID().String()

	deps.householdRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*household.Household{
		{
			ID:              "hh-1",
			ProjectID:       projectID,
			ReferenceNumber: &ref,
			HeadPersonID:    &headID,
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}, 1, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/export/households", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	c.Set(string(middleware.CtxCanExport), true)
	deps.handler.ExportHouseholds(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), fmt.Sprintf("households-%s.csv", projectID))

	rows, err := parseCSV(w.Body.String())
	require.NoError(t, err)
	require.Len(t, rows, 2)
	assert.Equal(t, "hh-1", rows[1][0])
	assert.Equal(t, "HH-001", rows[1][1])
}
