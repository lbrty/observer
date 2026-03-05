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
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func TestNoteUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonNoteRepository(ctrl)
	uc := ucproject.NewNoteUseCase(mockRepo)

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

	out, err := uc.Update(context.Background(), "n1", ucproject.UpdateNoteInput{Body: "updated body"})
	require.NoError(t, err)
	assert.Equal(t, "n1", out.ID)
	assert.Equal(t, "updated body", out.Body)
	assert.Equal(t, "user1", *out.AuthorID)
}

func TestNoteUseCase_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonNoteRepository(ctrl)
	uc := ucproject.NewNoteUseCase(mockRepo)

	mockRepo.EXPECT().GetByID(gomock.Any(), "n1").Return(nil, errors.New("not found"))

	_, err := uc.Update(context.Background(), "n1", ucproject.UpdateNoteInput{Body: "new"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "get note for update")
}

func TestNoteUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonNoteRepository(ctrl)
	uc := ucproject.NewNoteUseCase(mockRepo)

	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	out, err := uc.Create(context.Background(), "p1", "author1", ucproject.CreateNoteInput{Body: "hello"})
	require.NoError(t, err)
	assert.Equal(t, "p1", out.PersonID)
	assert.Equal(t, "hello", out.Body)
	assert.NotEmpty(t, out.ID)
}
