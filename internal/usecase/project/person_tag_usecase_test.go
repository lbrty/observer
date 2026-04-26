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

func TestPersonTagUseCase_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonTagRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonTagUseCase(mockRepo, mockPersonRepo)

	mockRepo.EXPECT().List(gomock.Any(), "p1").Return([]string{"tag1", "tag2"}, nil)

	ids, err := uc.List(context.Background(), "p1")
	require.NoError(t, err)
	assert.Equal(t, []string{"tag1", "tag2"}, ids)
}

func TestPersonTagUseCase_Replace_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonTagRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonTagUseCase(mockRepo, mockPersonRepo)

	mockPersonRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&person.Person{ID: "p1", ProjectID: "proj-1"}, nil)
	mockRepo.EXPECT().ReplaceAll(gomock.Any(), "p1", []string{"tag1", "tag3"}).Return(nil)

	err := uc.Replace(context.Background(), "proj-1", "p1", []string{"tag1", "tag3"})
	require.NoError(t, err)
}

func TestPersonTagUseCase_Replace_WrongProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPersonTagRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewPersonTagUseCase(mockRepo, mockPersonRepo)

	mockPersonRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&person.Person{ID: "p1", ProjectID: "other-proj"}, nil)

	err := uc.Replace(context.Background(), "proj-1", "p1", []string{"tag1"})
	assert.ErrorIs(t, err, person.ErrPersonNotFound)
}
