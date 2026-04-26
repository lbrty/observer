package project_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	domaintag "github.com/lbrty/observer/internal/domain/tag"
	"github.com/lbrty/observer/internal/handler/handlertest"
	"github.com/lbrty/observer/internal/handler/project"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func newTagHandler(ctrl *gomock.Controller) (*project.TagHandler, *repomock.MockTagRepository) {
	repo := repomock.NewMockTagRepository(ctrl)
	uc := ucproject.NewTagUseCase(repo, nil)
	return project.NewTagHandler(uc), repo
}

func TestTagHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newTagHandler(ctrl)

	projectID := handlertest.TestID().String()
	now := time.Now().UTC()
	repo.EXPECT().List(gomock.Any(), projectID).Return([]*domaintag.Tag{
		{ID: handlertest.TestID().String(), ProjectID: projectID, Name: "urgent", Color: "#ff0000", CreatedAt: now},
		{ID: handlertest.TestID().String(), ProjectID: projectID, Name: "pending", Color: "#ffaa00", CreatedAt: now},
	}, nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/tags", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	tags := resp["tags"].([]any)
	assert.Len(t, tags, 2)
}

func TestTagHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _ := newTagHandler(ctrl)

	projectID := handlertest.TestID().String()
	c, w := handlertest.NewTestContextWithParams(http.MethodPost, "/projects/"+projectID+"/tags", map[string]any{}, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTagHandler_Create_NameExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newTagHandler(ctrl)

	projectID := handlertest.TestID().String()
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domaintag.ErrTagNameExists)

	c, w := handlertest.NewTestContextWithParams(http.MethodPost, "/projects/"+projectID+"/tags", map[string]any{
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

	projectID := handlertest.TestID().String()
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodPost, "/projects/"+projectID+"/tags", map[string]any{
		"name":  "urgent",
		"color": "#ff0000",
	}, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, "urgent", resp["name"])
	assert.Equal(t, "#ff0000", resp["color"])
}

func TestTagHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newTagHandler(ctrl)

	id := handlertest.TestID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, domaintag.ErrTagNotFound)

	c, w := handlertest.NewTestContextWithParams(http.MethodPut, "/projects/x/tags/"+id, map[string]any{
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
	id := handlertest.TestID().String()
	projectID := handlertest.TestID().String()
	existing := &domaintag.Tag{
		ID: id, ProjectID: projectID, Name: "urgent", Color: "#ff0000", CreatedAt: now,
	}

	repo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodPut, "/projects/"+projectID+"/tags/"+id, map[string]any{
		"name": "critical",
	}, gin.Params{
		{Key: "project_id", Value: projectID},
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, "critical", resp["name"])
}

func TestTagHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newTagHandler(ctrl)

	id := handlertest.TestID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, domaintag.ErrTagNotFound)

	c, w := handlertest.NewTestContextWithParams(http.MethodDelete, "/projects/x/tags/"+id, nil, gin.Params{
		{Key: "project_id", Value: "x"},
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTagHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newTagHandler(ctrl)

	projectID := "proj1"
	id := handlertest.TestID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(&domaintag.Tag{ID: id, ProjectID: projectID}, nil)
	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodDelete, "/projects/"+projectID+"/tags/"+id, nil, gin.Params{
		{Key: "project_id", Value: projectID},
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.Equal(t, "tag deleted", resp["message"])
}
