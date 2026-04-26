package admin_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/reference"
	"github.com/lbrty/observer/internal/handler/admin"
	"github.com/lbrty/observer/internal/handler/handlertest"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

func newCountryHandler(ctrl *gomock.Controller) (*admin.CountryHandler, *repomock.MockCountryRepository) {
	repo := repomock.NewMockCountryRepository(ctrl)
	uc := ucadmin.NewCountryUseCase(repo)
	return admin.NewCountryHandler(uc), repo
}

func TestCountryHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCountryHandler(ctrl)

	now := time.Now().UTC()
	repo.EXPECT().List(gomock.Any()).Return([]*reference.Country{
		{ID: handlertest.TestID().String(), Name: "Kyrgyzstan", Code: "KG", CreatedAt: now, UpdatedAt: now},
		{ID: handlertest.TestID().String(), Name: "Ukraine", Code: "UA", CreatedAt: now, UpdatedAt: now},
	}, nil)

	c, w := handlertest.NewTestContext(http.MethodGet, "/admin/countries", nil)
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	countries := resp["countries"].([]any)
	assert.Len(t, countries, 2)
}

func TestCountryHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCountryHandler(ctrl)

	id := handlertest.TestID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, reference.ErrCountryNotFound)

	c, w := handlertest.NewTestContextWithParams(http.MethodGet, "/admin/countries/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCountryHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCountryHandler(ctrl)

	now := time.Now().UTC()
	id := handlertest.TestID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(&reference.Country{
		ID: id, Name: "Kyrgyzstan", Code: "KG", CreatedAt: now, UpdatedAt: now,
	}, nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodGet, "/admin/countries/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, "Kyrgyzstan", resp["name"])
	assert.Equal(t, "KG", resp["code"])
}

func TestCountryHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _ := newCountryHandler(ctrl)

	c, w := handlertest.NewTestContext(http.MethodPost, "/admin/countries", map[string]any{})
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCountryHandler_Create_CodeExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCountryHandler(ctrl)

	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(reference.ErrCountryCodeExists)

	c, w := handlertest.NewTestContext(http.MethodPost, "/admin/countries", map[string]any{
		"name": "Kyrgyzstan",
		"code": "KG",
	})
	h.Create(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestCountryHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCountryHandler(ctrl)

	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := handlertest.NewTestContext(http.MethodPost, "/admin/countries", map[string]any{
		"name": "Kyrgyzstan",
		"code": "KG",
	})
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, "Kyrgyzstan", resp["name"])
	assert.Equal(t, "KG", resp["code"])
}

func TestCountryHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCountryHandler(ctrl)

	id := handlertest.TestID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, reference.ErrCountryNotFound)

	c, w := handlertest.NewTestContextWithParams(http.MethodPatch, "/admin/countries/"+id, map[string]any{
		"name": "Updated",
	}, gin.Params{{Key: "id", Value: id}})
	h.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCountryHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCountryHandler(ctrl)

	now := time.Now().UTC()
	id := handlertest.TestID().String()
	existing := &reference.Country{
		ID: id, Name: "Kyrgyzstan", Code: "KG", CreatedAt: now, UpdatedAt: now,
	}

	repo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodPatch, "/admin/countries/"+id, map[string]any{
		"name": "Updated Name",
	}, gin.Params{{Key: "id", Value: id}})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
}

func TestCountryHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCountryHandler(ctrl)

	id := handlertest.TestID().String()
	repo.EXPECT().Delete(gomock.Any(), id).Return(reference.ErrCountryNotFound)

	c, w := handlertest.NewTestContextWithParams(http.MethodDelete, "/admin/countries/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCountryHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newCountryHandler(ctrl)

	id := handlertest.TestID().String()
	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodDelete, "/admin/countries/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.Equal(t, "country deleted", resp["message"])
}
