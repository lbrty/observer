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

func TestPersonCategoryUseCase_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonCategoryRepository(ctrl)
	uc := ucproject.NewPersonCategoryUseCase(mockRepo)

	mockRepo.EXPECT().List(gomock.Any(), "p1").Return([]string{"cat1", "cat2"}, nil)

	ids, err := uc.List(context.Background(), "p1")
	require.NoError(t, err)
	assert.Equal(t, []string{"cat1", "cat2"}, ids)
}

func TestPersonCategoryUseCase_List_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonCategoryRepository(ctrl)
	uc := ucproject.NewPersonCategoryUseCase(mockRepo)

	repoErr := errors.New("db error")
	mockRepo.EXPECT().List(gomock.Any(), "p1").Return(nil, repoErr)

	_, err := uc.List(context.Background(), "p1")
	assert.ErrorIs(t, err, repoErr)
}

func TestPersonCategoryUseCase_Replace_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonCategoryRepository(ctrl)
	uc := ucproject.NewPersonCategoryUseCase(mockRepo)

	mockRepo.EXPECT().ReplaceAll(gomock.Any(), "p1", []string{"cat1", "cat3"}).Return(nil)

	err := uc.Replace(context.Background(), "p1", []string{"cat1", "cat3"})
	require.NoError(t, err)
}

func TestPersonCategoryUseCase_Replace_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonCategoryRepository(ctrl)
	uc := ucproject.NewPersonCategoryUseCase(mockRepo)

	repoErr := errors.New("constraint violation")
	mockRepo.EXPECT().ReplaceAll(gomock.Any(), "p1", []string{"bad"}).Return(repoErr)

	err := uc.Replace(context.Background(), "p1", []string{"bad"})
	assert.ErrorIs(t, err, repoErr)
}
