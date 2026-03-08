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

func newCategoryHandler(ctrl *gomock.Controller) (*handler.CategoryHandler, *repomock.MockCategoryRepository) {
	repo := repomock.NewMockCategoryRepository(ctrl)
	uc := ucadmin.NewCategoryUseCase(repo)
	return handler.NewCategoryHandler(uc), repo
}

func TestCategoryHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCategoryHandler(ctrl)

	now := time.Now().UTC()
	desc := "Internally displaced people"
	repo.EXPECT().List(gomock.Any()).Return([]*reference.Category{
		{ID: testID().String(), Name: "IDP", Description: &desc, CreatedAt: now, UpdatedAt: now},
		{ID: testID().String(), Name: "Refugee", CreatedAt: now, UpdatedAt: now},
	}, nil)

	c, w := newTestContext(http.MethodGet, "/admin/categories", nil)
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	categories := resp["categories"].([]any)
	assert.Len(t, categories, 2)
}

func TestCategoryHandler_List_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCategoryHandler(ctrl)

	repo.EXPECT().List(gomock.Any()).Return(nil, fmt.Errorf("db error"))

	c, w := newTestContext(http.MethodGet, "/admin/categories", nil)
	h.List(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCategoryHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCategoryHandler(ctrl)

	id := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, reference.ErrCategoryNotFound)

	c, w := newTestContextWithParams(http.MethodGet, "/admin/categories/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCategoryHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCategoryHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	desc := "Internally displaced people"
	repo.EXPECT().GetByID(gomock.Any(), id).Return(&reference.Category{
		ID: id, Name: "IDP", Description: &desc, CreatedAt: now, UpdatedAt: now,
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/admin/categories/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, "IDP", resp["name"])
	assert.Equal(t, "Internally displaced people", resp["description"])
}

func TestCategoryHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _ := newCategoryHandler(ctrl)

	c, w := newTestContext(http.MethodPost, "/admin/categories", map[string]any{})
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCategoryHandler_Create_NameExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCategoryHandler(ctrl)

	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(reference.ErrCategoryNameExists)

	c, w := newTestContext(http.MethodPost, "/admin/categories", map[string]any{
		"name": "IDP",
	})
	h.Create(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestCategoryHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCategoryHandler(ctrl)

	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	desc := "Internally displaced people"
	c, w := newTestContext(http.MethodPost, "/admin/categories", map[string]any{
		"name":        "IDP",
		"description": desc,
	})
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, "IDP", resp["name"])
}

func TestCategoryHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCategoryHandler(ctrl)

	id := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, reference.ErrCategoryNotFound)

	c, w := newTestContextWithParams(http.MethodPatch, "/admin/categories/"+id, map[string]any{
		"name": "Updated",
	}, gin.Params{{Key: "id", Value: id}})
	h.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCategoryHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCategoryHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	existing := &reference.Category{
		ID: id, Name: "IDP", CreatedAt: now, UpdatedAt: now,
	}

	repo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPatch, "/admin/categories/"+id, map[string]any{
		"name": "Updated Category",
	}, gin.Params{{Key: "id", Value: id}})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
}

func TestCategoryHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCategoryHandler(ctrl)

	id := testID().String()
	repo.EXPECT().Delete(gomock.Any(), id).Return(reference.ErrCategoryNotFound)

	c, w := newTestContextWithParams(http.MethodDelete, "/admin/categories/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCategoryHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCategoryHandler(ctrl)

	id := testID().String()
	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	c, w := newTestContextWithParams(http.MethodDelete, "/admin/categories/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "category deleted", resp["message"])
}
