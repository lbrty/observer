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

func newOfficeHandler(ctrl *gomock.Controller) (*handler.OfficeHandler, *repomock.MockOfficeRepository) {
	repo := repomock.NewMockOfficeRepository(ctrl)
	uc := ucadmin.NewOfficeUseCase(repo)
	return handler.NewOfficeHandler(uc), repo
}

func TestOfficeHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newOfficeHandler(ctrl)

	now := time.Now().UTC()
	placeID := testID().String()
	repo.EXPECT().List(gomock.Any()).Return([]*reference.Office{
		{ID: testID().String(), Name: "Main Office", PlaceID: &placeID, CreatedAt: now, UpdatedAt: now},
		{ID: testID().String(), Name: "Branch Office", CreatedAt: now, UpdatedAt: now},
	}, nil)

	c, w := newTestContext(http.MethodGet, "/admin/offices", nil)
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	offices := resp["offices"].([]any)
	assert.Len(t, offices, 2)
}

func TestOfficeHandler_List_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newOfficeHandler(ctrl)

	repo.EXPECT().List(gomock.Any()).Return(nil, fmt.Errorf("db error"))

	c, w := newTestContext(http.MethodGet, "/admin/offices", nil)
	h.List(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestOfficeHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newOfficeHandler(ctrl)

	id := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, reference.ErrOfficeNotFound)

	c, w := newTestContextWithParams(http.MethodGet, "/admin/offices/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestOfficeHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newOfficeHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	placeID := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(&reference.Office{
		ID: id, Name: "Main Office", PlaceID: &placeID, CreatedAt: now, UpdatedAt: now,
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/admin/offices/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, "Main Office", resp["name"])
	assert.Equal(t, placeID, resp["place_id"])
}

func TestOfficeHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _ := newOfficeHandler(ctrl)

	c, w := newTestContext(http.MethodPost, "/admin/offices", map[string]any{})
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOfficeHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newOfficeHandler(ctrl)

	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	placeID := testID().String()
	c, w := newTestContext(http.MethodPost, "/admin/offices", map[string]any{
		"name":     "New Office",
		"place_id": placeID,
	})
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, "New Office", resp["name"])
}

func TestOfficeHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newOfficeHandler(ctrl)

	id := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, reference.ErrOfficeNotFound)

	c, w := newTestContextWithParams(http.MethodPatch, "/admin/offices/"+id, map[string]any{
		"name": "Updated",
	}, gin.Params{{Key: "id", Value: id}})
	h.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestOfficeHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newOfficeHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	existing := &reference.Office{
		ID: id, Name: "Main Office", CreatedAt: now, UpdatedAt: now,
	}

	repo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPatch, "/admin/offices/"+id, map[string]any{
		"name": "Updated Office",
	}, gin.Params{{Key: "id", Value: id}})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
}

func TestOfficeHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newOfficeHandler(ctrl)

	id := testID().String()
	repo.EXPECT().Delete(gomock.Any(), id).Return(reference.ErrOfficeNotFound)

	c, w := newTestContextWithParams(http.MethodDelete, "/admin/offices/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestOfficeHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newOfficeHandler(ctrl)

	id := testID().String()
	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	c, w := newTestContextWithParams(http.MethodDelete, "/admin/offices/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "office deleted", resp["message"])
}
