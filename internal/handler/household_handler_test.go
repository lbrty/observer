package handler_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/household"
	"github.com/lbrty/observer/internal/handler"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func newHouseholdHandler(ctrl *gomock.Controller) (*handler.HouseholdHandler, *repomock.MockHouseholdRepository, *repomock.MockHouseholdMemberRepository) {
	repo := repomock.NewMockHouseholdRepository(ctrl)
	memberRepo := repomock.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(repo, memberRepo)
	return handler.NewHouseholdHandler(uc), repo, memberRepo
}

func TestHouseholdHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newHouseholdHandler(ctrl)

	projectID := testID().String()
	now := time.Now().UTC()
	repo.EXPECT().List(gomock.Any(), projectID, 1, 20).Return([]*household.Household{
		{ID: testID().String(), ProjectID: projectID, CreatedAt: now, UpdatedAt: now},
		{ID: testID().String(), ProjectID: projectID, CreatedAt: now, UpdatedAt: now},
	}, 2, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/households", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	households := resp["households"].([]any)
	assert.Len(t, households, 2)
	assert.Equal(t, float64(2), resp["total"])
}

func TestHouseholdHandler_List_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newHouseholdHandler(ctrl)

	projectID := testID().String()
	repo.EXPECT().List(gomock.Any(), projectID, 1, 20).Return(nil, 0, fmt.Errorf("db error"))

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/households", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.List(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHouseholdHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newHouseholdHandler(ctrl)

	id := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, household.ErrHouseholdNotFound)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/households/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHouseholdHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, memberRepo := newHouseholdHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	projectID := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(&household.Household{
		ID:        id,
		ProjectID: projectID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)
	memberRepo.EXPECT().List(gomock.Any(), id).Return([]*household.Member{}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/households/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, projectID, resp["project_id"])
}

func TestHouseholdHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newHouseholdHandler(ctrl)

	projectID := testID().String()
	refNum := "HH-001"

	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPost, "/projects/"+projectID+"/households", map[string]any{
		"reference_number": refNum,
	}, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, projectID, resp["project_id"])
	assert.Equal(t, refNum, resp["reference_number"])
}

func TestHouseholdHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newHouseholdHandler(ctrl)

	id := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, household.ErrHouseholdNotFound)

	c, w := newTestContextWithParams(http.MethodPatch, "/projects/x/households/"+id, map[string]any{
		"reference_number": "HH-002",
	}, gin.Params{
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHouseholdHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newHouseholdHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	projectID := testID().String()
	existing := &household.Household{
		ID:        id,
		ProjectID: projectID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	repo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	refNum := "HH-002"
	c, w := newTestContextWithParams(http.MethodPatch, "/projects/"+projectID+"/households/"+id, map[string]any{
		"reference_number": refNum,
	}, gin.Params{
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, refNum, resp["reference_number"])
}

func TestHouseholdHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newHouseholdHandler(ctrl)

	id := testID().String()
	repo.EXPECT().Delete(gomock.Any(), id).Return(household.ErrHouseholdNotFound)

	c, w := newTestContextWithParams(http.MethodDelete, "/projects/x/households/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHouseholdHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo, _ := newHouseholdHandler(ctrl)

	id := testID().String()
	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	c, w := newTestContextWithParams(http.MethodDelete, "/projects/x/households/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "household deleted", resp["message"])
}

func TestHouseholdHandler_AddMember_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, _ := newHouseholdHandler(ctrl)

	id := testID().String()
	c, w := newTestContextWithParams(http.MethodPost, "/projects/x/households/"+id+"/members", map[string]any{}, gin.Params{
		{Key: "id", Value: id},
	})
	h.AddMember(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHouseholdHandler_AddMember_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, memberRepo := newHouseholdHandler(ctrl)

	id := testID().String()
	personID := testID().String()

	memberRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPost, "/projects/x/households/"+id+"/members", map[string]any{
		"person_id":    personID,
		"relationship": "spouse",
	}, gin.Params{
		{Key: "id", Value: id},
	})
	h.AddMember(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, personID, resp["person_id"])
	assert.Equal(t, "spouse", resp["relationship"])
}

func TestHouseholdHandler_RemoveMember_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, memberRepo := newHouseholdHandler(ctrl)

	id := testID().String()
	personID := testID().String()
	memberRepo.EXPECT().Remove(gomock.Any(), id, personID).Return(household.ErrMemberNotFound)

	c, w := newTestContextWithParams(http.MethodDelete, "/projects/x/households/"+id+"/members/"+personID, nil, gin.Params{
		{Key: "id", Value: id},
		{Key: "person_id", Value: personID},
	})
	h.RemoveMember(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHouseholdHandler_RemoveMember_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, memberRepo := newHouseholdHandler(ctrl)

	id := testID().String()
	personID := testID().String()
	memberRepo.EXPECT().Remove(gomock.Any(), id, personID).Return(nil)

	c, w := newTestContextWithParams(http.MethodDelete, "/projects/x/households/"+id+"/members/"+personID, nil, gin.Params{
		{Key: "id", Value: id},
		{Key: "person_id", Value: personID},
	})
	h.RemoveMember(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "member removed", resp["message"])
}
