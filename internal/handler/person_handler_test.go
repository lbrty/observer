package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/person"
	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/middleware"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

type personTestDeps struct {
	personRepo      *repomock.MockPersonRepository
	personTagRepo   *repomock.MockPersonTagRepository
	personCatRepo   *repomock.MockPersonCategoryRepository
	handler         *handler.PersonHandler
}

func newPersonTestDeps(ctrl *gomock.Controller) *personTestDeps {
	personRepo := repomock.NewMockPersonRepository(ctrl)
	tagRepo := repomock.NewMockPersonTagRepository(ctrl)
	catRepo := repomock.NewMockPersonCategoryRepository(ctrl)

	h := handler.NewPersonHandler(
		ucproject.NewPersonUseCase(personRepo, tagRepo),
		ucproject.NewPersonCategoryUseCase(catRepo),
		ucproject.NewPersonTagUseCase(tagRepo),
	)

	return &personTestDeps{
		personRepo:    personRepo,
		personTagRepo: tagRepo,
		personCatRepo: catRepo,
		handler:       h,
	}
}

func setSensitivityContext(c *gin.Context) {
	c.Set(string(middleware.CtxCanViewContact), true)
	c.Set(string(middleware.CtxCanViewPersonal), true)
}

func samplePerson(projectID string) *person.Person {
	now := time.Now().UTC()
	return &person.Person{
		ID:           testID().String(),
		ProjectID:    projectID,
		FirstName:    "Aman",
		Sex:          person.SexMale,
		CaseStatus:   person.CaseStatusNew,
		PhoneNumbers: json.RawMessage("[]"),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func TestPersonHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newPersonTestDeps(ctrl)

	projectID := testID().String()
	p := samplePerson(projectID)

	deps.personRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*person.Person{p}, 1, nil)
	deps.personTagRepo.EXPECT().ListBulk(gomock.Any(), []string{p.ID}).Return(map[string][]string{
		p.ID: {"tag1", "tag2"},
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/people", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	setSensitivityContext(c)
	deps.handler.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	people := resp["people"].([]any)
	assert.Len(t, people, 1)
	assert.Equal(t, float64(1), resp["total"])
}

func TestPersonHandler_List_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newPersonTestDeps(ctrl)

	projectID := testID().String()

	deps.personRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, 0, fmt.Errorf("db error"))

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/people", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	setSensitivityContext(c)
	deps.handler.List(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPersonHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newPersonTestDeps(ctrl)

	personID := testID().String()

	deps.personRepo.EXPECT().GetByID(gomock.Any(), personID).Return(nil, person.ErrPersonNotFound)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/people/"+personID, nil, gin.Params{
		{Key: "person_id", Value: personID},
	})
	setSensitivityContext(c)
	deps.handler.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPersonHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newPersonTestDeps(ctrl)

	projectID := testID().String()
	p := samplePerson(projectID)

	deps.personRepo.EXPECT().GetByID(gomock.Any(), p.ID).Return(p, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/people/"+p.ID, nil, gin.Params{
		{Key: "person_id", Value: p.ID},
	})
	setSensitivityContext(c)
	deps.handler.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, p.ID, resp["id"])
	assert.Equal(t, "Aman", resp["first_name"])
}

func TestPersonHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newPersonTestDeps(ctrl)

	projectID := testID().String()

	c, w := newTestContextWithParams(http.MethodPost, "/projects/"+projectID+"/people", map[string]any{}, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	deps.handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPersonHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newPersonTestDeps(ctrl)

	projectID := testID().String()

	deps.personRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPost, "/projects/"+projectID+"/people", map[string]any{
		"first_name": "Aman",
	}, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	deps.handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, "Aman", resp["first_name"])
}

func TestPersonHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newPersonTestDeps(ctrl)

	personID := testID().String()

	deps.personRepo.EXPECT().GetByID(gomock.Any(), personID).Return(nil, person.ErrPersonNotFound)

	name := "Updated"
	c, w := newTestContextWithParams(http.MethodPatch, "/projects/x/people/"+personID, map[string]any{
		"first_name": name,
	}, gin.Params{
		{Key: "person_id", Value: personID},
	})
	deps.handler.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPersonHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newPersonTestDeps(ctrl)

	projectID := testID().String()
	p := samplePerson(projectID)

	deps.personRepo.EXPECT().GetByID(gomock.Any(), p.ID).Return(p, nil)
	deps.personRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPatch, "/projects/"+projectID+"/people/"+p.ID, map[string]any{
		"first_name": "Updated",
	}, gin.Params{
		{Key: "person_id", Value: p.ID},
	})
	deps.handler.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, p.ID, resp["id"])
	assert.Equal(t, "Updated", resp["first_name"])
}

func TestPersonHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newPersonTestDeps(ctrl)

	personID := testID().String()

	deps.personRepo.EXPECT().Delete(gomock.Any(), personID).Return(person.ErrPersonNotFound)

	c, w := newTestContextWithParams(http.MethodDelete, "/projects/x/people/"+personID, nil, gin.Params{
		{Key: "person_id", Value: personID},
	})
	deps.handler.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPersonHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newPersonTestDeps(ctrl)

	personID := testID().String()

	deps.personRepo.EXPECT().Delete(gomock.Any(), personID).Return(nil)

	c, w := newTestContextWithParams(http.MethodDelete, "/projects/x/people/"+personID, nil, gin.Params{
		{Key: "person_id", Value: personID},
	})
	deps.handler.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "person deleted", resp["message"])
}

func TestPersonHandler_ListCategories_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newPersonTestDeps(ctrl)

	personID := testID().String()
	catIDs := []string{testID().String(), testID().String()}

	deps.personCatRepo.EXPECT().List(gomock.Any(), personID).Return(catIDs, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/people/"+personID+"/categories", nil, gin.Params{
		{Key: "person_id", Value: personID},
	})
	deps.handler.ListCategories(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	ids := resp["category_ids"].([]any)
	assert.Len(t, ids, 2)
}

func TestPersonHandler_ReplaceCategories_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newPersonTestDeps(ctrl)

	personID := testID().String()
	catIDs := []string{testID().String(), testID().String()}

	deps.personCatRepo.EXPECT().ReplaceAll(gomock.Any(), personID, catIDs).Return(nil)

	c, w := newTestContextWithParams(http.MethodPut, "/projects/x/people/"+personID+"/categories", map[string]any{
		"ids": catIDs,
	}, gin.Params{
		{Key: "person_id", Value: personID},
	})
	deps.handler.ReplaceCategories(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	ids := resp["category_ids"].([]any)
	assert.Len(t, ids, 2)
}

func TestPersonHandler_ListTags_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newPersonTestDeps(ctrl)

	personID := testID().String()
	tagIDs := []string{testID().String(), testID().String()}

	deps.personTagRepo.EXPECT().List(gomock.Any(), personID).Return(tagIDs, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/people/"+personID+"/tags", nil, gin.Params{
		{Key: "person_id", Value: personID},
	})
	deps.handler.ListTags(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	ids := resp["tag_ids"].([]any)
	assert.Len(t, ids, 2)
}

func TestPersonHandler_ReplaceTags_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	deps := newPersonTestDeps(ctrl)

	personID := testID().String()
	tagIDs := []string{testID().String(), testID().String()}

	deps.personTagRepo.EXPECT().ReplaceAll(gomock.Any(), personID, tagIDs).Return(nil)

	c, w := newTestContextWithParams(http.MethodPut, "/projects/x/people/"+personID+"/tags", map[string]any{
		"ids": tagIDs,
	}, gin.Params{
		{Key: "person_id", Value: personID},
	})
	deps.handler.ReplaceTags(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	ids := resp["tag_ids"].([]any)
	assert.Len(t, ids, 2)
}
