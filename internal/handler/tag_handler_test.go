package handler_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	domaintag "github.com/lbrty/observer/internal/domain/tag"
	"github.com/lbrty/observer/internal/handler"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func newTagHandler(ctrl *gomock.Controller) (*handler.TagHandler, *repomock.MockTagRepository) {
	repo := repomock.NewMockTagRepository(ctrl)
	uc := ucproject.NewTagUseCase(repo)
	return handler.NewTagHandler(uc), repo
}

func TestTagHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newTagHandler(ctrl)

	projectID := testID().String()
	now := time.Now().UTC()
	repo.EXPECT().List(gomock.Any(), projectID).Return([]*domaintag.Tag{
		{ID: testID().String(), ProjectID: projectID, Name: "urgent", Color: "#ff0000", CreatedAt: now},
		{ID: testID().String(), ProjectID: projectID, Name: "pending", Color: "#ffaa00", CreatedAt: now},
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/tags", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	tags := resp["tags"].([]any)
	assert.Len(t, tags, 2)
}

func TestTagHandler_List_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newTagHandler(ctrl)

	projectID := testID().String()
	repo.EXPECT().List(gomock.Any(), projectID).Return(nil, fmt.Errorf("db error"))

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/tags", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.List(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTagHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _ := newTagHandler(ctrl)

	projectID := testID().String()
	c, w := newTestContextWithParams(http.MethodPost, "/projects/"+projectID+"/tags", map[string]any{}, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTagHandler_Create_NameExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newTagHandler(ctrl)

	projectID := testID().String()
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domaintag.ErrTagNameExists)

	c, w := newTestContextWithParams(http.MethodPost, "/projects/"+projectID+"/tags", map[string]any{
		"name":  "urgent",
		"color": "#ff0000",
	}, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.Create(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestTagHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newTagHandler(ctrl)

	projectID := testID().String()
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPost, "/projects/"+projectID+"/tags", map[string]any{
		"name":  "urgent",
		"color": "#ff0000",
	}, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, "urgent", resp["name"])
	assert.Equal(t, "#ff0000", resp["color"])
}

func TestTagHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newTagHandler(ctrl)

	id := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, domaintag.ErrTagNotFound)

	c, w := newTestContextWithParams(http.MethodPut, "/projects/x/tags/"+id, map[string]any{
		"name": "updated",
	}, gin.Params{
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTagHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newTagHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	projectID := testID().String()
	existing := &domaintag.Tag{
		ID: id, ProjectID: projectID, Name: "urgent", Color: "#ff0000", CreatedAt: now,
	}

	repo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPut, "/projects/"+projectID+"/tags/"+id, map[string]any{
		"name": "critical",
	}, gin.Params{
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, "critical", resp["name"])
}

func TestTagHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newTagHandler(ctrl)

	id := testID().String()
	repo.EXPECT().Delete(gomock.Any(), id).Return(domaintag.ErrTagNotFound)

	c, w := newTestContextWithParams(http.MethodDelete, "/projects/x/tags/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTagHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newTagHandler(ctrl)

	id := testID().String()
	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	c, w := newTestContextWithParams(http.MethodDelete, "/projects/x/tags/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "tag deleted", resp["message"])
}
