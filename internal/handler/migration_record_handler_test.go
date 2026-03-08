package handler_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/migration"
	"github.com/lbrty/observer/internal/handler"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func newMigrationRecordHandler(ctrl *gomock.Controller) (*handler.MigrationRecordHandler, *repomock.MockMigrationRecordRepository) {
	repo := repomock.NewMockMigrationRecordRepository(ctrl)
	uc := ucproject.NewMigrationRecordUseCase(repo)
	return handler.NewMigrationRecordHandler(uc), repo
}

func TestMigrationRecordHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newMigrationRecordHandler(ctrl)

	personID := testID().String()
	now := time.Now().UTC()
	repo.EXPECT().ListByPerson(gomock.Any(), personID).Return([]*migration.Record{
		{ID: testID().String(), PersonID: personID, CreatedAt: now},
		{ID: testID().String(), PersonID: personID, CreatedAt: now},
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/people/"+personID+"/migration-records", nil, gin.Params{
		{Key: "person_id", Value: personID},
	})
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	records := resp["records"].([]any)
	assert.Len(t, records, 2)
}

func TestMigrationRecordHandler_List_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newMigrationRecordHandler(ctrl)

	personID := testID().String()
	repo.EXPECT().ListByPerson(gomock.Any(), personID).Return(nil, fmt.Errorf("db error"))

	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/people/"+personID+"/migration-records", nil, gin.Params{
		{Key: "person_id", Value: personID},
	})
	h.List(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMigrationRecordHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newMigrationRecordHandler(ctrl)

	id := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, migration.ErrRecordNotFound)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/people/y/migration-records/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMigrationRecordHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newMigrationRecordHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	personID := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(&migration.Record{
		ID:        id,
		PersonID:  personID,
		CreatedAt: now,
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/x/people/"+personID+"/migration-records/"+id, nil, gin.Params{
		{Key: "id", Value: id},
	})
	h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, personID, resp["person_id"])
}

func TestMigrationRecordHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newMigrationRecordHandler(ctrl)

	personID := testID().String()
	reason := "conflict"

	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPost, "/projects/x/people/"+personID+"/migration-records", map[string]any{
		"movement_reason": reason,
	}, gin.Params{
		{Key: "person_id", Value: personID},
	})
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, personID, resp["person_id"])
	assert.Equal(t, reason, resp["movement_reason"])
}

func TestMigrationRecordHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newMigrationRecordHandler(ctrl)

	id := testID().String()
	repo.EXPECT().GetByID(gomock.Any(), id).Return(nil, migration.ErrRecordNotFound)

	c, w := newTestContextWithParams(http.MethodPatch, "/projects/x/people/y/migration-records/"+id, map[string]any{
		"notes": "updated",
	}, gin.Params{
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMigrationRecordHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newMigrationRecordHandler(ctrl)

	now := time.Now().UTC()
	id := testID().String()
	personID := testID().String()
	existing := &migration.Record{
		ID:        id,
		PersonID:  personID,
		CreatedAt: now,
	}

	repo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	notes := "updated notes"
	c, w := newTestContextWithParams(http.MethodPatch, "/projects/x/people/"+personID+"/migration-records/"+id, map[string]any{
		"notes": notes,
	}, gin.Params{
		{Key: "id", Value: id},
	})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, id, resp["id"])
	assert.Equal(t, notes, resp["notes"])
}
