package project_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/person"
	"github.com/lbrty/observer/internal/domain/support"
	"github.com/lbrty/observer/internal/handler/handlertest"
	"github.com/lbrty/observer/internal/handler/project"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func newSupportRecordHandler(ctrl *gomock.Controller) (*project.SupportRecordHandler, *repomock.MockSupportRecordRepository, *repomock.MockPersonRepository) {
	repo := repomock.NewMockSupportRecordRepository(ctrl)
	personRepo := repomock.NewMockPersonRepository(ctrl)
	uc := ucproject.NewSupportRecordUseCase(repo, personRepo, nil)
	return project.NewSupportRecordHandler(uc), repo, personRepo
}

func TestSupportRecordHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newSupportRecordHandler(ctrl)

	projectID := handlertest.TestID().String()
	now := time.Now().UTC()
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*support.Record{
		{ID: handlertest.TestID().String(), PersonID: handlertest.TestID().String(), ProjectID: projectID, Type: support.SupportTypeHumanitarian, CreatedAt: now, UpdatedAt: now},
		{ID: handlertest.TestID().String(), PersonID: handlertest.TestID().String(), ProjectID: projectID, Type: support.SupportTypeLegal, CreatedAt: now, UpdatedAt: now},
	}, 2, nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/support-records", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	records := resp["records"].([]any)
	assert.Len(t, records, 2)
	assert.Equal(t, float64(2), resp["total"])
}

func TestSupportRecordHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newSupportRecordHandler(ctrl)

	id := handlertest.TestID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, support.ErrRecordNotFound)

	c, w := handlertest.NewTestContextWithParams(http.MethodGet, "/projects/x/support-records/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSupportRecordHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newSupportRecordHandler(ctrl)

	now := time.Now().UTC()
	id := handlertest.TestID().String()
	projectID := handlertest.TestID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(&support.Record{
		ID:        id,
		PersonID:  handlertest.TestID().String(),
		ProjectID: projectID,
		Type:      support.SupportTypeSocial,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/support-records/"+id, nil, gin.Params{
		{Key: "project_id", Value: projectID},
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, "social", resp["type"])
}

func TestSupportRecordHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, _ := newSupportRecordHandler(ctrl)

	projectID := handlertest.TestID().String()
	c, w := handlertest.NewTestContextWithParams(http.MethodPost, "/projects/"+projectID+"/support-records", map[string]any{}, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	handlertest.SetAuthContext(c, handlertest.TestID())
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSupportRecordHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, personRepo := newSupportRecordHandler(ctrl)

	projectID := handlertest.TestID().String()
	personID := handlertest.TestID().String()
	userID := handlertest.TestID()

	personRepo.EXPECT().GetByID(gomock.Any(), personID).Return(&person.Person{ID: personID, ProjectID: projectID}, nil)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodPost, "/projects/"+projectID+"/support-records", map[string]any{
		"person_id": personID,
		"type":      "humanitarian",
	}, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	handlertest.SetAuthContext(c, userID)
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, personID, resp["person_id"])
	assert.Equal(t, projectID, resp["project_id"])
	assert.Equal(t, "humanitarian", resp["type"])
	assert.Equal(t, userID.String(), resp["recorded_by"])
}

func TestSupportRecordHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newSupportRecordHandler(ctrl)

	id := handlertest.TestID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, support.ErrRecordNotFound)

	c, w := handlertest.NewTestContextWithParams(http.MethodPatch, "/projects/x/support-records/"+id, map[string]any{
		"notes": "updated",
	}, gin.Params{
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSupportRecordHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newSupportRecordHandler(ctrl)

	now := time.Now().UTC()
	id := handlertest.TestID().String()
	projectID := handlertest.TestID().String()
	existing := &support.Record{
		ID:        id,
		PersonID:  handlertest.TestID().String(),
		ProjectID: projectID,
		Type:      support.SupportTypeHumanitarian,
		CreatedAt: now,
		UpdatedAt: now,
	}

	repo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodPatch, "/projects/"+projectID+"/support-records/"+id, map[string]any{
		"type": "legal",
	}, gin.Params{
		{Key: "project_id", Value: projectID},
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, "legal", resp["type"])
}

func TestSupportRecordHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newSupportRecordHandler(ctrl)

	id := handlertest.TestID().String()
	projectID := handlertest.TestID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, support.ErrRecordNotFound)

	c, w := handlertest.NewTestContextWithParams(http.MethodDelete, "/projects/"+projectID+"/support-records/"+id, nil, gin.Params{
		{Key: "project_id", Value: projectID},
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSupportRecordHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newSupportRecordHandler(ctrl)

	id := handlertest.TestID().String()
	projectID := handlertest.TestID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(&support.Record{ID: id, ProjectID: projectID}, nil)
	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodDelete, "/projects/"+projectID+"/support-records/"+id, nil, gin.Params{
		{Key: "project_id", Value: projectID},
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.Equal(t, "support record deleted", resp["message"])
}
