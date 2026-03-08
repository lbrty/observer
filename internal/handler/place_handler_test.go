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

func newPlaceHandler(ctrl *gomock.Controller) (*handler.PlaceHandler, *repomock.MockPlaceRepository) {
	repo := repomock.NewMockPlaceRepository(ctrl)
	uc := ucadmin.NewPlaceUseCase(repo)
	return handler.NewPlaceHandler(uc), repo
}

func TestPlaceHandler_List_AllPlaces(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newPlaceHandler(ctrl)

	now := time.Now().UTC()
	repo.EXPECT().ListAll(gomock.Any()).Return([]*reference.Place{
		{ID: testID().String(), StateID: testID().String(), Name: "Bishkek City", CreatedAt: now, UpdatedAt: now},
	}, nil)

	c, w := newTestContext(http.MethodGet, "/admin/places", nil)
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	places := resp["places"].([]any)
	assert.Len(t, places, 1)
}

func TestPlaceHandler_List_ByState(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newPlaceHandler(ctrl)

	now := time.Now().UTC()
	stateID := testID().String()
	repo.EXPECT().List(gomock.Any(), stateID).Return([]*reference.Place{
		{ID: testID().String(), StateID: stateID, Name: "Osh City", CreatedAt: now, UpdatedAt: now},
	}, nil)

	c, w := newTestContext(http.MethodGet, "/admin/places?state_id="+stateID, nil)
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	places := resp["places"].([]any)
	assert.Len(t, places, 1)
}

func TestPlaceHandler_List_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newPlaceHandler(ctrl)

	repo.EXPECT().ListAll(gomock.Any()).Return(nil, fmt.Errorf("db error"))

	c, w := newTestContext(http.MethodGet, "/admin/places", nil)
	h.List(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPlaceHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newPlaceHandler(ctrl)

	id := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, reference.ErrPlaceNotFound)

	c, w := newTestContextWithParams(http.MethodGet, "/admin/places/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPlaceHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newPlaceHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	stateID := testID().String()
	lat := 42.87
	lon := 74.59
	repo.EXPECT().GetByID(gomock.Any(), id).Return(&reference.Place{
		ID: id, StateID: stateID, Name: "Bishkek", Lat: &lat, Lon: &lon, CreatedAt: now, UpdatedAt: now,
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/admin/places/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, "Bishkek", resp["name"])
	assert.Equal(t, 42.87, resp["lat"])
	assert.Equal(t, 74.59, resp["lon"])
}

func TestPlaceHandler_Create_NoStateID(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _ := newPlaceHandler(ctrl)

	c, w := newTestContext(http.MethodPost, "/admin/places", map[string]any{
		"name": "Bishkek",
	})
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPlaceHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _ := newPlaceHandler(ctrl)

	stateID := testID().String()
	c, w := newTestContext(http.MethodPost, "/admin/places?state_id="+stateID, map[string]any{})
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPlaceHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newPlaceHandler(ctrl)

	stateID := testID().String()
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	lat := 42.87
	lon := 74.59
	c, w := newTestContext(http.MethodPost, "/admin/places?state_id="+stateID, map[string]any{
		"name": "Bishkek",
		"lat":  lat,
		"lon":  lon,
	})
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, "Bishkek", resp["name"])
	assert.Equal(t, stateID, resp["state_id"])
}

func TestPlaceHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newPlaceHandler(ctrl)

	id := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, reference.ErrPlaceNotFound)

	c, w := newTestContextWithParams(http.MethodPatch, "/admin/places/"+id, map[string]any{
		"name": "Updated",
	}, gin.Params{{Key: "id", Value: id}})
	h.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPlaceHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newPlaceHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	existing := &reference.Place{
		ID: id, StateID: testID().String(), Name: "Bishkek", CreatedAt: now, UpdatedAt: now,
	}

	repo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPatch, "/admin/places/"+id, map[string]any{
		"name": "Updated Name",
	}, gin.Params{{Key: "id", Value: id}})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
}

func TestPlaceHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newPlaceHandler(ctrl)

	id := testID().String()
	repo.EXPECT().Delete(gomock.Any(), id).Return(reference.ErrPlaceNotFound)

	c, w := newTestContextWithParams(http.MethodDelete, "/admin/places/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPlaceHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newPlaceHandler(ctrl)

	id := testID().String()
	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	c, w := newTestContextWithParams(http.MethodDelete, "/admin/places/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "place deleted", resp["message"])
}
