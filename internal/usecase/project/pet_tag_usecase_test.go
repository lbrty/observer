package project_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/pet"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func TestPetTagUseCase_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetTagRepository(ctrl)
	mockPetRepo := mock_repo.NewMockPetRepository(ctrl)
	uc := ucproject.NewPetTagUseCase(mockRepo, mockPetRepo)

	mockRepo.EXPECT().List(gomock.Any(), "pet1").Return([]string{"tag1", "tag2"}, nil)

	ids, err := uc.List(context.Background(), "pet1")
	require.NoError(t, err)
	assert.Equal(t, []string{"tag1", "tag2"}, ids)
}

func TestPetTagUseCase_Replace_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetTagRepository(ctrl)
	mockPetRepo := mock_repo.NewMockPetRepository(ctrl)
	uc := ucproject.NewPetTagUseCase(mockRepo, mockPetRepo)

	mockPetRepo.EXPECT().GetByID(gomock.Any(), "pet1").Return(&pet.Pet{ID: "pet1", ProjectID: "proj-1"}, nil)
	mockRepo.EXPECT().ReplaceAll(gomock.Any(), "pet1", []string{"tag1", "tag3"}).Return(nil)

	err := uc.Replace(context.Background(), "proj-1", "pet1", []string{"tag1", "tag3"})
	require.NoError(t, err)
}

func TestPetTagUseCase_Replace_WrongProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetTagRepository(ctrl)
	mockPetRepo := mock_repo.NewMockPetRepository(ctrl)
	uc := ucproject.NewPetTagUseCase(mockRepo, mockPetRepo)

	mockPetRepo.EXPECT().GetByID(gomock.Any(), "pet1").Return(&pet.Pet{ID: "pet1", ProjectID: "other-proj"}, nil)

	err := uc.Replace(context.Background(), "proj-1", "pet1", []string{"tag1"})
	assert.ErrorIs(t, err, pet.ErrPetNotFound)
}
