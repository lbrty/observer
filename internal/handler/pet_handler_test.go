package handler_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/pet"
	"github.com/lbrty/observer/internal/handler"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func newPetHandler(ctrl *gomock.Controller) (*handler.PetHandler, *repomock.MockPetRepository, *repomock.MockPetTagRepository) {
	petRepo := repomock.NewMockPetRepository(ctrl)
	petTagRepo := repomock.NewMockPetTagRepository(ctrl)
	uc := ucproject.NewPetUseCase(petRepo, petTagRepo)
	tagUC := ucproject.NewPetTagUseCase(petTagRepo)
	return handler.NewPetHandler(uc, tagUC), petRepo, petTagRepo
}

func TestPetHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, petRepo, petTagRepo := newPetHandler(ctrl)

	projectID := testID().String()
	now := time.Now().UTC()
	id1 := testID().String()
	id2 := testID().String()

	petRepo.EXPECT().List(gomock.Any(), projectID, "", []string(nil), 1, 20).Return([]*pet.Pet{
		{ID: id1, ProjectID: projectID, Name: "Rex", Status: pet.PetStatusRegistered, CreatedAt: now, UpdatedAt: now},
		{ID: id2, ProjectID: projectID, Name: "Luna", Status: pet.PetStatusAdopted, CreatedAt: now, UpdatedAt: now},
	}, 2, nil)
	petTagRepo.EXPECT().ListBulk(gomock.Any(), []string{id1, id2}).Return(map[string][]string{
		id1: {"tag1"},
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/pets", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	pets := resp["pets"].([]any)
	assert.Len(t, pets, 2)
	assert.Equal(t, float64(2), resp["total"])
}

func TestPetHandler_List_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, petRepo, _ := newPetHandler(ctrl)

	projectID := testID().String()
	petRepo.EXPECT().List(gomock.Any(), projectID, "", []string(nil), 1, 20).Return(nil, 0, fmt.Errorf("db error"))

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/pets", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.List(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPetHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, petRepo, _ := newPetHandler(ctrl)

	id := testID().String()
	petRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, pet.ErrPetNotFound)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/pets/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPetHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, petRepo, _ := newPetHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	projectID := testID().String()

	petRepo.EXPECT().GetByID(gomock.Any(), id).Return(&pet.Pet{
		ID: id, ProjectID: projectID, Name: "Rex", Status: pet.PetStatusRegistered,
		CreatedAt: now, UpdatedAt: now,
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/pets/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, "Rex", resp["name"])
}

func TestPetHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, _ := newPetHandler(ctrl)

	projectID := testID().String()
	c, w := newTestContextWithParams(http.MethodPost, "/projects/"+projectID+"/pets", map[string]any{}, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPetHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, petRepo, _ := newPetHandler(ctrl)

	projectID := testID().String()
	petRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPost, "/projects/"+projectID+"/pets", map[string]any{
		"name": "Rex",
	}, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, "Rex", resp["name"])
	assert.Equal(t, "unknown", resp["status"])
}

func TestPetHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, petRepo, _ := newPetHandler(ctrl)

	id := testID().String()
	petRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, pet.ErrPetNotFound)

	c, w := newTestContextWithParams(http.MethodPatch, "/projects/x/pets/"+id, map[string]any{
		"name": "updated",
	}, gin.Params{
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPetHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, petRepo, _ := newPetHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	projectID := testID().String()
	existing := &pet.Pet{
		ID: id, ProjectID: projectID, Name: "Rex", Status: pet.PetStatusRegistered,
		CreatedAt: now, UpdatedAt: now,
	}

	petRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	petRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPatch, "/projects/"+projectID+"/pets/"+id, map[string]any{
		"name": "Rex Jr",
	}, gin.Params{
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, "Rex Jr", resp["name"])
}

func TestPetHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, petRepo, _ := newPetHandler(ctrl)

	id := testID().String()
	petRepo.EXPECT().Delete(gomock.Any(), id).Return(pet.ErrPetNotFound)

	c, w := newTestContextWithParams(http.MethodDelete, "/projects/x/pets/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPetHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, petRepo, _ := newPetHandler(ctrl)

	id := testID().String()
	petRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	c, w := newTestContextWithParams(http.MethodDelete, "/projects/x/pets/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "pet deleted", resp["message"])
}

func TestPetHandler_ListTags_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, petTagRepo := newPetHandler(ctrl)

	id := testID().String()
	petTagRepo.EXPECT().List(gomock.Any(), id).Return([]string{"tag1", "tag2"}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/pets/"+id+"/tags", nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.ListTags(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	tagIDs := resp["tag_ids"].([]any)
	assert.Len(t, tagIDs, 2)
}

func TestPetHandler_ReplaceTags_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, _ := newPetHandler(ctrl)

	id := testID().String()
	c, w := newTestContextWithParams(http.MethodPut, "/projects/x/pets/"+id+"/tags", map[string]any{}, gin.Params{
		{Key: "id", Value: id},
	})
	h.ReplaceTags(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPetHandler_ReplaceTags_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, petTagRepo := newPetHandler(ctrl)

	id := testID().String()
	tagIDs := []string{"tag1", "tag2"}
	petTagRepo.EXPECT().ReplaceAll(gomock.Any(), id, tagIDs).Return(nil)

	c, w := newTestContextWithParams(http.MethodPut, "/projects/x/pets/"+id+"/tags", map[string]any{
		"ids": tagIDs,
	}, gin.Params{
		{Key: "id", Value: id},
	})
	h.ReplaceTags(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	resultIDs := resp["tag_ids"].([]any)
	assert.Len(t, resultIDs, 2)
}
