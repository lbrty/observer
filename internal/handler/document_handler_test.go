package handler_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/document"
	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/middleware"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	storagemock "github.com/lbrty/observer/internal/storage/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func newDocumentHandler(ctrl *gomock.Controller) (*handler.DocumentHandler, *repomock.MockDocumentRepository, *storagemock.MockFileStorage) {
	docRepo := repomock.NewMockDocumentRepository(ctrl)
	fs := storagemock.NewMockFileStorage(ctrl)
	uc := ucproject.NewDocumentUseCase(docRepo, fs)
	return handler.NewDocumentHandler(uc), docRepo, fs
}

func setCanViewDocuments(c *gin.Context, canView bool) {
	c.Set(string(middleware.CtxCanViewDocuments), canView)
}

func TestDocumentHandler_List_NoPermission(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, _ := newDocumentHandler(ctrl)

	personID := testID().String()
	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/people/"+personID+"/documents", nil, gin.Params{
		{Key: "person_id", Value: personID},
	})
	setCanViewDocuments(c, false)
	h.List(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDocumentHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, docRepo, _ := newDocumentHandler(ctrl)

	personID := testID().String()
	now := time.Now().UTC()

	docRepo.EXPECT().List(gomock.Any(), personID).Return([]*document.Document{
		{ID: testID().String(), PersonID: personID, ProjectID: testID().String(), Name: "passport.pdf", Path: "a/b/c", MimeType: "application/pdf", Size: 1024, CreatedAt: now},
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/people/"+personID+"/documents", nil, gin.Params{
		{Key: "person_id", Value: personID},
	})
	setCanViewDocuments(c, true)
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	docs := resp["documents"].([]any)
	assert.Len(t, docs, 1)
}

func TestDocumentHandler_List_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, docRepo, _ := newDocumentHandler(ctrl)

	personID := testID().String()
	docRepo.EXPECT().List(gomock.Any(), personID).Return(nil, fmt.Errorf("db error"))

	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/people/"+personID+"/documents", nil, gin.Params{
		{Key: "person_id", Value: personID},
	})
	setCanViewDocuments(c, true)
	h.List(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDocumentHandler_Get_NoPermission(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, _ := newDocumentHandler(ctrl)

	id := testID().String()
	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/documents/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	setCanViewDocuments(c, false)
	h.Get(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDocumentHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, docRepo, _ := newDocumentHandler(ctrl)

	id := testID().String()
	docRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, document.ErrDocumentNotFound)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/documents/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	setCanViewDocuments(c, true)
	h.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDocumentHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, docRepo, _ := newDocumentHandler(ctrl)

	id := testID().String()
	now := time.Now().UTC()

	docRepo.EXPECT().GetByID(gomock.Any(), id).Return(&document.Document{
		ID: id, PersonID: testID().String(), ProjectID: testID().String(),
		Name: "passport.pdf", Path: "a/b/c", MimeType: "application/pdf", Size: 1024, CreatedAt: now,
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/documents/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	setCanViewDocuments(c, true)
	h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, "passport.pdf", resp["name"])
}

func TestDocumentHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, docRepo, _ := newDocumentHandler(ctrl)

	id := testID().String()
	docRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, document.ErrDocumentNotFound)

	newName := "updated.pdf"
	c, w := newTestContextWithParams(http.MethodPatch, "/projects/x/documents/"+id, map[string]any{
		"name": newName,
	}, gin.Params{
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDocumentHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, docRepo, _ := newDocumentHandler(ctrl)

	id := testID().String()
	now := time.Now().UTC()
	existing := &document.Document{
		ID: id, PersonID: testID().String(), ProjectID: testID().String(),
		Name: "passport.pdf", Path: "a/b/c", MimeType: "application/pdf", Size: 1024, CreatedAt: now,
	}

	docRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	docRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPatch, "/projects/x/documents/"+id, map[string]any{
		"name": "renamed.pdf",
	}, gin.Params{
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, "renamed.pdf", resp["name"])
}

func TestDocumentHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, docRepo, _ := newDocumentHandler(ctrl)

	id := testID().String()
	docRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, document.ErrDocumentNotFound)

	c, w := newTestContextWithParams(http.MethodDelete, "/projects/x/documents/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDocumentHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, docRepo, fs := newDocumentHandler(ctrl)

	id := testID().String()
	now := time.Now().UTC()
	existing := &document.Document{
		ID: id, PersonID: testID().String(), ProjectID: testID().String(),
		Name: "passport.pdf", Path: "a/b/c", MimeType: "application/pdf", Size: 1024, CreatedAt: now,
	}

	docRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	docRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)
	fs.EXPECT().Delete(gomock.Any(), "a/b/c").Return(nil)

	c, w := newTestContextWithParams(http.MethodDelete, "/projects/x/documents/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "document deleted", resp["message"])
}
