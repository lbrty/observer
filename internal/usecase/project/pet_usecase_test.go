package project_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/pet"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func TestPetUseCase_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetRepository(ctrl)
	mockTagRepo := mock_repo.NewMockPetTagRepository(ctrl)
	uc := ucproject.NewPetUseCase(mockRepo, mockTagRepo)

	now := time.Now().UTC()
	ownerID := "p1"
	mockRepo.EXPECT().List(gomock.Any(), "proj1", "", []string(nil), gomock.Any(), gomock.Any()).Return([]*pet.Pet{
		{ID: "pet1", ProjectID: "proj1", OwnerID: &ownerID, Name: "Buddy", Status: pet.PetStatusRegistered, CreatedAt: now, UpdatedAt: now},
		{ID: "pet2", ProjectID: "proj1", Name: "Max", Status: pet.PetStatusUnknown, CreatedAt: now, UpdatedAt: now},
	}, 2, nil)
	mockTagRepo.EXPECT().ListBulk(gomock.Any(), []string{"pet1", "pet2"}).Return(map[string][]string{
		"pet1": {"tag1", "tag2"},
	}, nil)

	out, err := uc.List(context.Background(), "proj1", ucproject.ListPetsInput{Page: 1, PerPage: 20})
	require.NoError(t, err)
	assert.Len(t, out.Pets, 2)
	assert.Equal(t, 2, out.Total)
	assert.Equal(t, "Buddy", out.Pets[0].Name)
	assert.Equal(t, []string{"tag1", "tag2"}, out.Pets[0].TagIDs)
	assert.Equal(t, []string{}, out.Pets[1].TagIDs)
}

func TestPetUseCase_List_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetRepository(ctrl)
	mockTagRepo := mock_repo.NewMockPetTagRepository(ctrl)
	uc := ucproject.NewPetUseCase(mockRepo, mockTagRepo)

	repoErr := errors.New("db error")
	mockRepo.EXPECT().List(gomock.Any(), "proj1", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, 0, repoErr)

	_, err := uc.List(context.Background(), "proj1", ucproject.ListPetsInput{})
	assert.ErrorIs(t, err, repoErr)
}

func TestPetUseCase_List_TagRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetRepository(ctrl)
	mockTagRepo := mock_repo.NewMockPetTagRepository(ctrl)
	uc := ucproject.NewPetUseCase(mockRepo, mockTagRepo)

	now := time.Now().UTC()
	mockRepo.EXPECT().List(gomock.Any(), "proj1", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*pet.Pet{
		{ID: "pet1", ProjectID: "proj1", Name: "Buddy", Status: pet.PetStatusRegistered, CreatedAt: now, UpdatedAt: now},
	}, 1, nil)

	tagErr := errors.New("tag query failed")
	mockTagRepo.EXPECT().ListBulk(gomock.Any(), []string{"pet1"}).Return(nil, tagErr)

	_, err := uc.List(context.Background(), "proj1", ucproject.ListPetsInput{Page: 1, PerPage: 20})
	assert.ErrorIs(t, err, tagErr)
}

func TestPetUseCase_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetRepository(ctrl)
	uc := ucproject.NewPetUseCase(mockRepo, nil)

	now := time.Now().UTC()
	ownerID := "p1"
	mockRepo.EXPECT().GetByID(gomock.Any(), "pet1").Return(&pet.Pet{
		ID:        "pet1",
		ProjectID: "proj1",
		OwnerID:   &ownerID,
		Name:      "Buddy",
		Status:    pet.PetStatusRegistered,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	dto, err := uc.Get(context.Background(), "pet1")
	require.NoError(t, err)
	assert.Equal(t, "pet1", dto.ID)
	assert.Equal(t, "Buddy", dto.Name)
	assert.Equal(t, "registered", dto.Status)
	require.NotNil(t, dto.OwnerID)
	assert.Equal(t, "p1", *dto.OwnerID)
}

func TestPetUseCase_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetRepository(ctrl)
	uc := ucproject.NewPetUseCase(mockRepo, nil)

	mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, pet.ErrPetNotFound)

	_, err := uc.Get(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, pet.ErrPetNotFound)
}

func TestPetUseCase_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetRepository(ctrl)
	uc := ucproject.NewPetUseCase(mockRepo, nil)

	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, p *pet.Pet) error {
		assert.NotEmpty(t, p.ID)
		assert.Equal(t, "proj1", p.ProjectID)
		assert.Equal(t, "Buddy", p.Name)
		assert.Equal(t, pet.PetStatusRegistered, p.Status)
		return nil
	})

	dto, err := uc.Create(context.Background(), "proj1", ucproject.CreatePetInput{
		Name:   "Buddy",
		Status: ptr("registered"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Buddy", dto.Name)
	assert.Equal(t, "registered", dto.Status)
	assert.Equal(t, "proj1", dto.ProjectID)
}

func TestPetUseCase_Create_DefaultStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetRepository(ctrl)
	uc := ucproject.NewPetUseCase(mockRepo, nil)

	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, p *pet.Pet) error {
		assert.Equal(t, pet.PetStatusUnknown, p.Status)
		return nil
	})

	dto, err := uc.Create(context.Background(), "proj1", ucproject.CreatePetInput{
		Name: "Stray",
	})
	require.NoError(t, err)
	assert.Equal(t, "unknown", dto.Status)
}

func TestPetUseCase_Create_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetRepository(ctrl)
	uc := ucproject.NewPetUseCase(mockRepo, nil)

	repoErr := errors.New("insert failed")
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repoErr)

	_, err := uc.Create(context.Background(), "proj1", ucproject.CreatePetInput{Name: "Buddy"})
	assert.ErrorIs(t, err, repoErr)
}

func TestPetUseCase_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetRepository(ctrl)
	uc := ucproject.NewPetUseCase(mockRepo, nil)

	now := time.Now().UTC()
	existing := &pet.Pet{
		ID:        "pet1",
		ProjectID: "proj1",
		Name:      "Buddy",
		Status:    pet.PetStatusUnknown,
		CreatedAt: now,
		UpdatedAt: now,
	}
	mockRepo.EXPECT().GetByID(gomock.Any(), "pet1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, p *pet.Pet) error {
		assert.Equal(t, "Rex", p.Name)
		assert.Equal(t, pet.PetStatusAdopted, p.Status)
		return nil
	})

	dto, err := uc.Update(context.Background(), "pet1", ucproject.UpdatePetInput{
		Name:   ptr("Rex"),
		Status: ptr("adopted"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Rex", dto.Name)
	assert.Equal(t, "adopted", dto.Status)
}

func TestPetUseCase_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetRepository(ctrl)
	uc := ucproject.NewPetUseCase(mockRepo, nil)

	mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, pet.ErrPetNotFound)

	_, err := uc.Update(context.Background(), "nonexistent", ucproject.UpdatePetInput{})
	assert.ErrorIs(t, err, pet.ErrPetNotFound)
}

func TestPetUseCase_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetRepository(ctrl)
	uc := ucproject.NewPetUseCase(mockRepo, nil)

	mockRepo.EXPECT().Delete(gomock.Any(), "pet1").Return(nil)

	err := uc.Delete(context.Background(), "pet1")
	require.NoError(t, err)
}

func TestPetUseCase_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPetRepository(ctrl)
	uc := ucproject.NewPetUseCase(mockRepo, nil)

	mockRepo.EXPECT().Delete(gomock.Any(), "nonexistent").Return(pet.ErrPetNotFound)

	err := uc.Delete(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, pet.ErrPetNotFound)
}
