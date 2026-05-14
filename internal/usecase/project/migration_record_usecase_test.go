package project_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/migration"
	"github.com/lbrty/observer/internal/domain/person"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func TestMigrationRecordUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMigrationRecordRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	auditRepo := mock_repo.NewMockAuditLogRepository(ctrl)
	auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	uc := ucproject.NewMigrationRecordUseCase(mockRepo, mockPersonRepo, auditUC)

	reason := migration.ReasonConflict
	existing := &migration.Record{
		ID:             "mr1",
		PersonID:       "p1",
		MovementReason: &reason,
		CreatedAt:      time.Now().UTC(),
	}

	mockRepo.EXPECT().GetByID(gomock.Any(), "mr1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	newReason := "economic"
	notes := "updated notes"
	out, err := uc.Update(context.Background(), "proj1", "p1", "mr1", ucproject.UpdateMigrationRecordInput{
		MovementReason: &newReason,
		Notes:          &notes,
	})

	require.NoError(t, err)
	assert.Equal(t, "mr1", out.ID)
	assert.Equal(t, "economic", *out.MovementReason)
	assert.Equal(t, "updated notes", *out.Notes)
}

func TestMigrationRecordUseCase_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMigrationRecordRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	auditRepo := mock_repo.NewMockAuditLogRepository(ctrl)
	auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	uc := ucproject.NewMigrationRecordUseCase(mockRepo, mockPersonRepo, auditUC)

	mockRepo.EXPECT().GetByID(gomock.Any(), "mr1").Return(nil, errors.New("not found"))

	_, err := uc.Update(context.Background(), "proj1", "p1", "mr1", ucproject.UpdateMigrationRecordInput{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "get migration record for update")
}

func TestMigrationRecordUseCase_Update_PartialFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMigrationRecordRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	auditRepo := mock_repo.NewMockAuditLogRepository(ctrl)
	auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	uc := ucproject.NewMigrationRecordUseCase(mockRepo, mockPersonRepo, auditUC)

	fromPlace := "place1"
	existing := &migration.Record{
		ID:          "mr1",
		PersonID:    "p1",
		FromPlaceID: &fromPlace,
		CreatedAt:   time.Now().UTC(),
	}

	mockRepo.EXPECT().GetByID(gomock.Any(), "mr1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, r *migration.Record) error {
		assert.Equal(t, "place1", *r.FromPlaceID)
		assert.Equal(t, "place2", *r.DestinationPlaceID)
		return nil
	})

	dest := "place2"
	_, err := uc.Update(context.Background(), "proj1", "p1", "mr1", ucproject.UpdateMigrationRecordInput{
		DestinationPlaceID: &dest,
	})
	require.NoError(t, err)
}

func TestMigrationRecordUseCase_Get_CrossPersonIDOR(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMigrationRecordRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewMigrationRecordUseCase(mockRepo, mockPersonRepo, nil)

	mockRepo.EXPECT().GetByID(gomock.Any(), "mr1").Return(&migration.Record{
		ID: "mr1", PersonID: "other-person",
	}, nil)

	_, err := uc.Get(context.Background(), "p1", "mr1")
	assert.ErrorIs(t, err, migration.ErrRecordNotFound)
}

func TestMigrationRecordUseCase_Update_CrossPersonIDOR(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMigrationRecordRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewMigrationRecordUseCase(mockRepo, mockPersonRepo, nil)

	mockRepo.EXPECT().GetByID(gomock.Any(), "mr1").Return(&migration.Record{
		ID: "mr1", PersonID: "other-person",
	}, nil)

	_, err := uc.Update(context.Background(), "proj1", "p1", "mr1", ucproject.UpdateMigrationRecordInput{})
	assert.ErrorIs(t, err, migration.ErrRecordNotFound)
}

func TestMigrationRecordUseCase_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMigrationRecordRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	auditRepo := mock_repo.NewMockAuditLogRepository(ctrl)
	auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	uc := ucproject.NewMigrationRecordUseCase(mockRepo, mockPersonRepo, auditUC)

	mockRepo.EXPECT().GetByID(gomock.Any(), "mr1").Return(&migration.Record{
		ID:       "mr1",
		PersonID: "p1",
	}, nil)
	mockRepo.EXPECT().Delete(gomock.Any(), "mr1").Return(nil)

	err := uc.Delete(context.Background(), "proj1", "p1", "mr1")
	require.NoError(t, err)
}

func TestMigrationRecordUseCase_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMigrationRecordRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewMigrationRecordUseCase(mockRepo, mockPersonRepo, nil)

	mockRepo.EXPECT().GetByID(gomock.Any(), "mr99").Return(nil, migration.ErrRecordNotFound)

	err := uc.Delete(context.Background(), "proj1", "p1", "mr99")
	assert.ErrorIs(t, err, migration.ErrRecordNotFound)
}

func TestMigrationRecordUseCase_Delete_CrossPersonIDOR(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMigrationRecordRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewMigrationRecordUseCase(mockRepo, mockPersonRepo, nil)

	mockRepo.EXPECT().GetByID(gomock.Any(), "mr1").Return(&migration.Record{
		ID: "mr1", PersonID: "other-person",
	}, nil)

	err := uc.Delete(context.Background(), "proj1", "p1", "mr1")
	assert.ErrorIs(t, err, migration.ErrRecordNotFound)
}

func TestMigrationRecordUseCase_ListByPerson_WrongProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMigrationRecordRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewMigrationRecordUseCase(mockRepo, mockPersonRepo, nil)

	mockPersonRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&person.Person{ID: "p1", ProjectID: "other-proj"}, nil)

	_, err := uc.ListByPerson(context.Background(), "proj1", "p1")
	assert.ErrorIs(t, err, person.ErrPersonNotFound)
}

func TestMigrationRecordUseCase_Create_WrongProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMigrationRecordRepository(ctrl)
	mockPersonRepo := mock_repo.NewMockPersonRepository(ctrl)
	uc := ucproject.NewMigrationRecordUseCase(mockRepo, mockPersonRepo, nil)

	mockPersonRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&person.Person{ID: "p1", ProjectID: "other-proj"}, nil)

	_, err := uc.Create(context.Background(), "proj1", "p1", ucproject.CreateMigrationRecordInput{})
	assert.ErrorIs(t, err, person.ErrPersonNotFound)
}
