package project_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/person"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func TestPersonCategoryUseCase_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonCategoryRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonCategoryUseCase(mockRepo, mockPersonRepo)

	mockPersonRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&person.Person{ID: "p1", ProjectID: "proj-1"}, nil)
	mockRepo.EXPECT().List(gomock.Any(), "p1").Return([]string{"cat1", "cat2"}, nil)

	ids, err := uc.List(context.Background(), "proj-1", "p1")
	require.NoError(t, err)
	assert.Equal(t, []string{"cat1", "cat2"}, ids)
}

func TestPersonCategoryUseCase_List_WrongProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonCategoryRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonCategoryUseCase(mockRepo, mockPersonRepo)

	mockPersonRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&person.Person{ID: "p1", ProjectID: "other-proj"}, nil)

	_, err := uc.List(context.Background(), "proj-1", "p1")
	assert.ErrorIs(t, err, person.ErrPersonNotFound)
}

func TestPersonCategoryUseCase_Replace_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonCategoryRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonCategoryUseCase(mockRepo, mockPersonRepo)

	mockPersonRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&person.Person{ID: "p1", ProjectID: "proj-1"}, nil)
	mockRepo.EXPECT().ReplaceAll(gomock.Any(), "p1", []string{"cat1", "cat3"}).Return(nil)

	err := uc.Replace(context.Background(), "proj-1", "p1", []string{"cat1", "cat3"})
	require.NoError(t, err)
}

func TestPersonCategoryUseCase_Replace_WrongProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonCategoryRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonCategoryUseCase(mockRepo, mockPersonRepo)

	mockPersonRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&person.Person{ID: "p1", ProjectID: "other-proj"}, nil)

	err := uc.Replace(context.Background(), "proj-1", "p1", []string{"cat1"})
	assert.ErrorIs(t, err, person.ErrPersonNotFound)
}
