package project_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/note"
	"github.com/lbrty/observer/internal/domain/person"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func TestNoteUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonNoteRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	auditRepo := mock_repo.NewMockAuditLogRepository(ctrl)
	auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	uc := ucproject.NewNoteUseCase(mockRepo, mockPersonRepo, auditUC)

	author := "user1"
	existing := &note.Note{
		ID:        "n1",
		PersonID:  "p1",
		AuthorID:  &author,
		Body:      "original body",
		CreatedAt: time.Now().UTC(),
	}

	mockRepo.EXPECT().GetByID(gomock.Any(), "n1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, n *note.Note) error {
		assert.Equal(t, "updated body", n.Body)
		return nil
	})

	out, err := uc.Update(context.Background(), "proj1", "p1", "n1", ucproject.UpdateNoteInput{Body: "updated body"})
	require.NoError(t, err)
	assert.Equal(t, "n1", out.ID)
	assert.Equal(t, "updated body", out.Body)
	assert.Equal(t, "user1", *out.AuthorID)
}

func TestNoteUseCase_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonNoteRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	auditRepo := mock_repo.NewMockAuditLogRepository(ctrl)
	auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	uc := ucproject.NewNoteUseCase(mockRepo, mockPersonRepo, auditUC)

	mockRepo.EXPECT().GetByID(gomock.Any(), "n1").Return(nil, errors.New("not found"))

	_, err := uc.Update(context.Background(), "proj1", "p1", "n1", ucproject.UpdateNoteInput{Body: "new"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "get note for update")
}

func TestNoteUseCase_Update_CrossPersonIDOR(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonNoteRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewNoteUseCase(mockRepo, mockPersonRepo, nil)

	mockRepo.EXPECT().GetByID(gomock.Any(), "n1").Return(&note.Note{
		ID: "n1", PersonID: "other-person",
	}, nil)

	_, err := uc.Update(context.Background(), "proj1", "p1", "n1", ucproject.UpdateNoteInput{Body: "x"})
	assert.ErrorIs(t, err, note.ErrNoteNotFound)
}

func TestNoteUseCase_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonNoteRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	auditRepo := mock_repo.NewMockAuditLogRepository(ctrl)
	auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	uc := ucproject.NewNoteUseCase(mockRepo, mockPersonRepo, auditUC)

	mockRepo.EXPECT().GetByID(gomock.Any(), "n1").Return(&note.Note{ID: "n1", PersonID: "p1"}, nil)
	mockRepo.EXPECT().Delete(gomock.Any(), "n1").Return(nil)

	err := uc.Delete(context.Background(), "proj1", "p1", "n1")
	require.NoError(t, err)
}

func TestNoteUseCase_Delete_CrossPersonIDOR(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonNoteRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewNoteUseCase(mockRepo, mockPersonRepo, nil)

	mockRepo.EXPECT().GetByID(gomock.Any(), "n1").Return(&note.Note{
		ID: "n1", PersonID: "other-person",
	}, nil)

	err := uc.Delete(context.Background(), "proj1", "p1", "n1")
	assert.ErrorIs(t, err, note.ErrNoteNotFound)
}

func TestNoteUseCase_List_WrongProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonNoteRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewNoteUseCase(mockRepo, mockPersonRepo, nil)

	mockPersonRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&person.Person{ID: "p1", ProjectID: "other-proj"}, nil)

	_, err := uc.List(context.Background(), "proj1", "p1")
	assert.ErrorIs(t, err, person.ErrPersonNotFound)
}

func TestNoteUseCase_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonNoteRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	auditRepo := mock_repo.NewMockAuditLogRepository(ctrl)
	auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	uc := ucproject.NewNoteUseCase(mockRepo, mockPersonRepo, auditUC)

	mockPersonRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&person.Person{ID: "p1", ProjectID: "proj1"}, nil)
	mockRepo.EXPECT().List(gomock.Any(), "p1").Return([]*note.Note{
		{ID: "n1", PersonID: "p1", Body: "hello"},
	}, nil)

	out, err := uc.List(context.Background(), "proj1", "p1")
	require.NoError(t, err)
	assert.Len(t, out, 1)
	assert.Equal(t, "n1", out[0].ID)
}

func TestNoteUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonNoteRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	auditRepo := mock_repo.NewMockAuditLogRepository(ctrl)
	auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	uc := ucproject.NewNoteUseCase(mockRepo, mockPersonRepo, auditUC)

	mockPersonRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&person.Person{ID: "p1", ProjectID: "proj1"}, nil)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	out, err := uc.Create(context.Background(), "proj1", "p1", "author1", ucproject.CreateNoteInput{
		Body: "hello",
	})
	require.NoError(t, err)
	assert.Equal(t, "p1", out.PersonID)
	assert.Equal(t, "hello", out.Body)
	assert.NotEmpty(t, out.ID)
}

func TestNoteUseCase_Create_WrongProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonNoteRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewNoteUseCase(mockRepo, mockPersonRepo, nil)

	mockPersonRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&person.Person{ID: "p1", ProjectID: "other-proj"}, nil)

	_, err := uc.Create(context.Background(), "proj1", "p1", "author1", ucproject.CreateNoteInput{Body: "hello"})
	assert.ErrorIs(t, err, person.ErrPersonNotFound)
}
