package project_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/support"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func TestSupportRecordUseCase_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockSupportRecordRepository(ctrl)
	uc := ucproject.NewSupportRecordUseCase(mockRepo)

	now := time.Now().UTC()
	mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*support.Record{
		{ID: "sr1", PersonID: "p1", ProjectID: "proj1", Type: support.SupportTypeLegal, CreatedAt: now, UpdatedAt: now},
		{ID: "sr2", PersonID: "p2", ProjectID: "proj1", Type: support.SupportTypeMedical, CreatedAt: now, UpdatedAt: now},
	}, 2, nil)

	out, err := uc.List(context.Background(), "proj1", ucproject.ListSupportRecordsInput{Page: 1, PerPage: 20})
	require.NoError(t, err)
	assert.Len(t, out.Records, 2)
	assert.Equal(t, 2, out.Total)
	assert.Equal(t, "sr1", out.Records[0].ID)
	assert.Equal(t, "legal", out.Records[0].Type)
}

func TestSupportRecordUseCase_List_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockSupportRecordRepository(ctrl)
	uc := ucproject.NewSupportRecordUseCase(mockRepo)

	repoErr := errors.New("db error")
	mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, 0, repoErr)

	_, err := uc.List(context.Background(), "proj1", ucproject.ListSupportRecordsInput{})
	assert.ErrorIs(t, err, repoErr)
}

func TestSupportRecordUseCase_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockSupportRecordRepository(ctrl)
	uc := ucproject.NewSupportRecordUseCase(mockRepo)

	now := time.Now().UTC()
	sphere := support.SphereHousingAssistance
	mockRepo.EXPECT().GetByID(gomock.Any(), "sr1").Return(&support.Record{
		ID:        "sr1",
		PersonID:  "p1",
		ProjectID: "proj1",
		Type:      support.SupportTypeLegal,
		Sphere:    &sphere,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	dto, err := uc.Get(context.Background(), "sr1")
	require.NoError(t, err)
	assert.Equal(t, "sr1", dto.ID)
	assert.Equal(t, "legal", dto.Type)
	require.NotNil(t, dto.Sphere)
	assert.Equal(t, "housing_assistance", *dto.Sphere)
}

func TestSupportRecordUseCase_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockSupportRecordRepository(ctrl)
	uc := ucproject.NewSupportRecordUseCase(mockRepo)

	mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, support.ErrRecordNotFound)

	_, err := uc.Get(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, support.ErrRecordNotFound)
}

func TestSupportRecordUseCase_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockSupportRecordRepository(ctrl)
	uc := ucproject.NewSupportRecordUseCase(mockRepo)

	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, r *support.Record) error {
		assert.NotEmpty(t, r.ID)
		assert.Equal(t, "proj1", r.ProjectID)
		assert.Equal(t, "p1", r.PersonID)
		assert.Equal(t, support.SupportTypeLegal, r.Type)
		require.NotNil(t, r.RecordedBy)
		assert.Equal(t, "user1", *r.RecordedBy)
		return nil
	})

	dto, err := uc.Create(context.Background(), "proj1", "user1", ucproject.CreateSupportRecordInput{
		PersonID: "p1",
		Type:     "legal",
	})
	require.NoError(t, err)
	assert.Equal(t, "p1", dto.PersonID)
	assert.Equal(t, "proj1", dto.ProjectID)
	assert.Equal(t, "legal", dto.Type)
}

func TestSupportRecordUseCase_Create_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockSupportRecordRepository(ctrl)
	uc := ucproject.NewSupportRecordUseCase(mockRepo)

	repoErr := errors.New("insert failed")
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repoErr)

	_, err := uc.Create(context.Background(), "proj1", "user1", ucproject.CreateSupportRecordInput{
		PersonID: "p1",
		Type:     "legal",
	})
	assert.ErrorIs(t, err, repoErr)
}

func TestSupportRecordUseCase_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockSupportRecordRepository(ctrl)
	uc := ucproject.NewSupportRecordUseCase(mockRepo)

	now := time.Now().UTC()
	existing := &support.Record{
		ID:        "sr1",
		PersonID:  "p1",
		ProjectID: "proj1",
		Type:      support.SupportTypeLegal,
		CreatedAt: now,
		UpdatedAt: now,
	}
	mockRepo.EXPECT().GetByID(gomock.Any(), "sr1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, r *support.Record) error {
		assert.Equal(t, support.SupportTypeMedical, r.Type)
		return nil
	})

	dto, err := uc.Update(context.Background(), "sr1", ucproject.UpdateSupportRecordInput{
		Type: ptr("medical"),
	})
	require.NoError(t, err)
	assert.Equal(t, "medical", dto.Type)
}

func TestSupportRecordUseCase_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockSupportRecordRepository(ctrl)
	uc := ucproject.NewSupportRecordUseCase(mockRepo)

	mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, support.ErrRecordNotFound)

	_, err := uc.Update(context.Background(), "nonexistent", ucproject.UpdateSupportRecordInput{})
	assert.ErrorIs(t, err, support.ErrRecordNotFound)
}

func TestSupportRecordUseCase_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockSupportRecordRepository(ctrl)
	uc := ucproject.NewSupportRecordUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "sr1").Return(nil)

	err := uc.Delete(context.Background(), "sr1")
	require.NoError(t, err)
}

func TestSupportRecordUseCase_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockSupportRecordRepository(ctrl)
	uc := ucproject.NewSupportRecordUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "nonexistent").Return(support.ErrRecordNotFound)

	err := uc.Delete(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, support.ErrRecordNotFound)
}
