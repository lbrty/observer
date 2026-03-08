package project_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/tag"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func TestTagUseCase_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockTagRepository(ctrl)
	uc := ucproject.NewTagUseCase(mockRepo)

	now := time.Now().UTC()
	mockRepo.EXPECT().List(gomock.Any(), "proj1").Return([]*tag.Tag{
		{ID: "t1", ProjectID: "proj1", Name: "urgent", Color: "#ff0000", CreatedAt: now},
		{ID: "t2", ProjectID: "proj1", Name: "follow-up", Color: "#00ff00", CreatedAt: now},
	}, nil)

	dtos, err := uc.List(context.Background(), "proj1")
	require.NoError(t, err)
	assert.Len(t, dtos, 2)
	assert.Equal(t, "t1", dtos[0].ID)
	assert.Equal(t, "urgent", dtos[0].Name)
	assert.Equal(t, "#ff0000", dtos[0].Color)
	assert.Equal(t, "follow-up", dtos[1].Name)
}

func TestTagUseCase_List_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockTagRepository(ctrl)
	uc := ucproject.NewTagUseCase(mockRepo)

	repoErr := errors.New("db connection lost")
	mockRepo.EXPECT().List(gomock.Any(), "proj1").Return(nil, repoErr)

	_, err := uc.List(context.Background(), "proj1")
	assert.ErrorIs(t, err, repoErr)
}

func TestTagUseCase_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockTagRepository(ctrl)
	uc := ucproject.NewTagUseCase(mockRepo)

	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, tg *tag.Tag) error {
		assert.NotEmpty(t, tg.ID)
		assert.Equal(t, "proj1", tg.ProjectID)
		assert.Equal(t, "urgent", tg.Name)
		assert.Equal(t, "#ff0000", tg.Color)
		return nil
	})

	dto, err := uc.Create(context.Background(), "proj1", ucproject.CreateTagInput{
		Name:  "urgent",
		Color: "#ff0000",
	})
	require.NoError(t, err)
	assert.Equal(t, "urgent", dto.Name)
	assert.Equal(t, "#ff0000", dto.Color)
	assert.Equal(t, "proj1", dto.ProjectID)
}

func TestTagUseCase_Create_DuplicateName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockTagRepository(ctrl)
	uc := ucproject.NewTagUseCase(mockRepo)

	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tag.ErrTagNameExists)

	_, err := uc.Create(context.Background(), "proj1", ucproject.CreateTagInput{Name: "urgent"})
	assert.ErrorIs(t, err, tag.ErrTagNameExists)
}

func TestTagUseCase_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockTagRepository(ctrl)
	uc := ucproject.NewTagUseCase(mockRepo)

	existing := &tag.Tag{ID: "t1", ProjectID: "proj1", Name: "urgent", Color: "#ff0000"}
	mockRepo.EXPECT().GetByID(gomock.Any(), "t1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, tg *tag.Tag) error {
		assert.Equal(t, "critical", tg.Name)
		assert.Equal(t, "#cc0000", tg.Color)
		return nil
	})

	dto, err := uc.Update(context.Background(), "t1", ucproject.UpdateTagInput{
		Name:  ptr("critical"),
		Color: ptr("#cc0000"),
	})
	require.NoError(t, err)
	assert.Equal(t, "critical", dto.Name)
	assert.Equal(t, "#cc0000", dto.Color)
}

func TestTagUseCase_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockTagRepository(ctrl)
	uc := ucproject.NewTagUseCase(mockRepo)

	mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, tag.ErrTagNotFound)

	_, err := uc.Update(context.Background(), "nonexistent", ucproject.UpdateTagInput{Name: ptr("x")})
	assert.ErrorIs(t, err, tag.ErrTagNotFound)
}

func TestTagUseCase_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockTagRepository(ctrl)
	uc := ucproject.NewTagUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "t1").Return(nil)

	err := uc.Delete(context.Background(), "t1")
	require.NoError(t, err)
}

func TestTagUseCase_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockTagRepository(ctrl)
	uc := ucproject.NewTagUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "nonexistent").Return(tag.ErrTagNotFound)

	err := uc.Delete(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, tag.ErrTagNotFound)
}
