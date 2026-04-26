package project_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/note"
	"github.com/lbrty/observer/internal/handler/handlertest"
	"github.com/lbrty/observer/internal/handler/project"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func newNoteHandler(ctrl *gomock.Controller) (*project.NoteHandler, *repomock.MockPersonNoteRepository) {
	repo := repomock.NewMockPersonNoteRepository(ctrl)
	uc := ucproject.NewNoteUseCase(repo, nil)
	return project.NewNoteHandler(uc), repo
}

func TestNoteHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newNoteHandler(ctrl)

	personID := handlertest.TestID().String()
	now := time.Now().UTC()
	authorID := handlertest.TestID().String()
	repo.EXPECT().List(gomock.Any(), personID).Return([]*note.Note{
		{ID: handlertest.TestID().String(), PersonID: personID, AuthorID: &authorID, Body: "note 1", CreatedAt: now},
		{ID: handlertest.TestID().String(), PersonID: personID, AuthorID: &authorID, Body: "note 2", CreatedAt: now},
	}, nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodGet, "/projects/x/people/"+personID+"/notes", nil, gin.Params{
		{Key: "person_id", Value: personID},
	})
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	notes := resp["notes"].([]any)
	assert.Len(t, notes, 2)
}

func TestNoteHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _ := newNoteHandler(ctrl)

	personID := handlertest.TestID().String()
	c, w := handlertest.NewTestContextWithParams(http.MethodPost, "/projects/x/people/"+personID+"/notes", map[string]any{}, gin.Params{
		{Key: "person_id", Value: personID},
	})
	handlertest.SetAuthContext(c, handlertest.TestID())
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNoteHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newNoteHandler(ctrl)

	personID := handlertest.TestID().String()
	userID := handlertest.TestID()

	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodPost, "/projects/x/people/"+personID+"/notes", map[string]any{
		"body": "important note",
	}, gin.Params{
		{Key: "person_id", Value: personID},
	})
	handlertest.SetAuthContext(c, userID)
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, personID, resp["person_id"])
	assert.Equal(t, userID.String(), resp["author_id"])
	assert.Equal(t, "important note", resp["body"])
}

func TestNoteHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newNoteHandler(ctrl)

	id := handlertest.TestID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, note.ErrNoteNotFound)

	c, w := handlertest.NewTestContextWithParams(http.MethodPatch, "/projects/x/people/y/notes/"+id, map[string]any{
		"body": "updated",
	}, gin.Params{
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestNoteHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newNoteHandler(ctrl)

	now := time.Now().UTC()
	id := handlertest.TestID().String()
	personID := handlertest.TestID().String()
	authorID := handlertest.TestID().String()
	existing := &note.Note{
		ID:        id,
		PersonID:  personID,
		AuthorID:  &authorID,
		Body:      "old body",
		CreatedAt: now,
	}

	repo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodPatch, "/projects/x/people/"+personID+"/notes/"+id, map[string]any{
		"body": "new body",
	}, gin.Params{
		{Key: "person_id", Value: personID},
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, "new body", resp["body"])
}

func TestNoteHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newNoteHandler(ctrl)

	id := handlertest.TestID().String()
	personID := handlertest.TestID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, note.ErrNoteNotFound)

	c, w := handlertest.NewTestContextWithParams(http.MethodDelete, "/projects/x/people/"+personID+"/notes/"+id, nil, gin.Params{
		{Key: "person_id", Value: personID},
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestNoteHandler_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newNoteHandler(ctrl)

	id := handlertest.TestID().String()
	personID := handlertest.TestID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(&note.Note{ID: id, PersonID: personID}, nil)
	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodDelete, "/projects/x/people/"+personID+"/notes/"+id, nil, gin.Params{
		{Key: "person_id", Value: personID},
		{Key: "id", Value: id},
	})
	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.Equal(t, "note deleted", resp["message"])
}
