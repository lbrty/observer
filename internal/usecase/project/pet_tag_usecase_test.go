package project_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func TestPetTagUseCase_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetTagRepository(ctrl)
	uc := ucproject.NewPetTagUseCase(mockRepo)

	mockRepo.EXPECT().List(gomock.Any(), "pet1").Return([]string{"tag1", "tag2"}, nil)

	ids, err := uc.List(context.Background(), "pet1")
	require.NoError(t, err)
	assert.Equal(t, []string{"tag1", "tag2"}, ids)
}

func TestPetTagUseCase_List_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetTagRepository(ctrl)
	uc := ucproject.NewPetTagUseCase(mockRepo)

	repoErr := errors.New("db error")
	mockRepo.EXPECT().List(gomock.Any(), "pet1").Return(nil, repoErr)

	_, err := uc.List(context.Background(), "pet1")
	assert.ErrorIs(t, err, repoErr)
}

func TestPetTagUseCase_Replace_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetTagRepository(ctrl)
	uc := ucproject.NewPetTagUseCase(mockRepo)

	mockRepo.EXPECT().ReplaceAll(gomock.Any(), "pet1", []string{"tag1", "tag3"}).Return(nil)

	err := uc.Replace(context.Background(), "pet1", []string{"tag1", "tag3"})
	require.NoError(t, err)
}

func TestPetTagUseCase_Replace_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetTagRepository(ctrl)
	uc := ucproject.NewPetTagUseCase(mockRepo)

	repoErr := errors.New("constraint violation")
	mockRepo.EXPECT().ReplaceAll(gomock.Any(), "pet1", []string{"bad"}).Return(repoErr)

	err := uc.Replace(context.Background(), "pet1", []string{"bad"})
	assert.ErrorIs(t, err, repoErr)
}
