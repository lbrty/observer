package handler_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/reference"
	"github.com/lbrty/observer/internal/handler"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

func newStateHandler(ctrl *gomock.Controller) (*handler.StateHandler, *repomock.MockStateRepository) {
	repo := repomock.NewMockStateRepository(ctrl)
	uc := ucadmin.NewStateUseCase(repo)
	return handler.NewStateHandler(uc), repo
}

func TestStateHandler_List_AllStates(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newStateHandler(ctrl)

	now := time.Now().UTC()
	repo.EXPECT().ListAll(gomock.Any()).Return([]*reference.State{
		{ID: testID().String(), CountryID: testID().String(), Name: "Bishkek", CreatedAt: now, UpdatedAt: now},
	}, nil)

	c, w := newTestContext(http.MethodGet, "/admin/states", nil)
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	states := resp["states"].([]any)
	assert.Len(t, states, 1)
}

func TestStateHandler_List_ByCountry(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newStateHandler(ctrl)

	now := time.Now().UTC()
	countryID := testID().String()
	repo.EXPECT().List(gomock.Any(), countryID).Return([]*reference.State{
		{ID: testID().String(), CountryID: countryID, Name: "Osh", CreatedAt: now, UpdatedAt: now},
	}, nil)

	c, w := newTestContext(http.MethodGet, "/admin/states?country_id="+countryID, nil)
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	states := resp["states"].([]any)
	assert.Len(t, states, 1)
}

func TestStateHandler_List_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newStateHandler(ctrl)

	repo.EXPECT().ListAll(gomock.Any()).Return(nil, fmt.Errorf("db error"))

	c, w := newTestContext(http.MethodGet, "/admin/states", nil)
	h.List(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestStateHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newStateHandler(ctrl)

	id := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, reference.ErrStateNotFound)

	c, w := newTestContextWithParams(http.MethodGet, "/admin/states/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStateHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newStateHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	countryID := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(&reference.State{
		ID: id, CountryID: countryID, Name: "Bishkek", CreatedAt: now, UpdatedAt: now,
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/admin/states/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, "Bishkek", resp["name"])
}

func TestStateHandler_Create_NoCountryID(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _ := newStateHandler(ctrl)

	c, w := newTestContext(http.MethodPost, "/admin/states", map[string]any{
		"name": "Bishkek",
	})
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStateHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _ := newStateHandler(ctrl)

	countryID := testID().String()
	c, w := newTestContext(http.MethodPost, "/admin/states?country_id="+countryID, map[string]any{})
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStateHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newStateHandler(ctrl)

	countryID := testID().String()
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContext(http.MethodPost, "/admin/states?country_id="+countryID, map[string]any{
		"name": "Bishkek",
	})
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, "Bishkek", resp["name"])
	assert.Equal(t, countryID, resp["country_id"])
}

func TestStateHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newStateHandler(ctrl)

	id := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, reference.ErrStateNotFound)

	c, w := newTestContextWithParams(http.MethodPatch, "/admin/states/"+id, map[string]any{
		"name": "Updated",
	}, gin.Params{{Key: "id", Value: id}})
	h.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStateHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newStateHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	existing := &reference.State{
		ID: id, CountryID: testID().String(), Name: "Bishkek", CreatedAt: now, UpdatedAt: now,
	}

	repo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPatch, "/admin/states/"+id, map[string]any{
		"name": "Updated Name",
	}, gin.Params{{Key: "id", Value: id}})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
}

func TestStateHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newStateHandler(ctrl)

	id := testID().String()
	repo.EXPECT().Delete(gomock.Any(), id).Return(reference.ErrStateNotFound)

	c, w := newTestContextWithParams(http.MethodDelete, "/admin/states/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStateHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newStateHandler(ctrl)

	id := testID().String()
	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	c, w := newTestContextWithParams(http.MethodDelete, "/admin/states/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "state deleted", resp["message"])
}
