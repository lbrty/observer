package admin_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/reference"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

func TestStateUseCase_ListAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockStateRepository(ctrl)
	uc := ucadmin.NewStateUseCase(mockRepo)

	mockRepo.EXPECT().ListAll(gomock.Any()).Return([]*reference.State{
		{ID: "s1", CountryID: "c1", Name: "Kyivska oblast", Code: ptr("KY")},
		{ID: "s2", CountryID: "c1", Name: "Odeska oblast"},
	}, nil)

	out, err := uc.ListAll(context.Background())
	require.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, "Kyivska oblast", out[0].Name)
	assert.Equal(t, ptr("KY"), out[0].Code)
	assert.Nil(t, out[1].Code)
}

func TestStateUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockStateRepository(ctrl)
	uc := ucadmin.NewStateUseCase(mockRepo)

	mockRepo.EXPECT().List(gomock.Any(), "c1").Return([]*reference.State{
		{ID: "s1", CountryID: "c1", Name: "Kyivska oblast"},
	}, nil)

	out, err := uc.List(context.Background(), "c1")
	require.NoError(t, err)
	assert.Len(t, out, 1)
	assert.Equal(t, "c1", out[0].CountryID)
}

func TestStateUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockStateRepository(ctrl)
	uc := ucadmin.NewStateUseCase(mockRepo)

	mockRepo.EXPECT().GetByID(gomock.Any(), "s1").Return(&reference.State{
		ID: "s1", CountryID: "c1", Name: "Kyivska oblast", ConflictZone: ptr("active"),
	}, nil)

	out, err := uc.Get(context.Background(), "s1")
	require.NoError(t, err)
	assert.Equal(t, "s1", out.ID)
	assert.Equal(t, "Kyivska oblast", out.Name)
	assert.Equal(t, ptr("active"), out.ConflictZone)
}

func TestStateUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockStateRepository(ctrl)
	uc := ucadmin.NewStateUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, s *reference.State) error {
			assert.NotEmpty(t, s.ID)
			assert.Equal(t, "c1", s.CountryID)
			assert.Equal(t, "Lvivska oblast", s.Name)
			assert.Equal(t, ptr("LV"), s.Code)
			return nil
		})

	out, err := uc.Create(context.Background(), "c1", ucadmin.CreateStateInput{
		Name: "Lvivska oblast",
		Code: ptr("LV"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Lvivska oblast", out.Name)
	assert.Equal(t, "c1", out.CountryID)
}

func TestStateUseCase_Update_PartialUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockStateRepository(ctrl)
	uc := ucadmin.NewStateUseCase(mockRepo)

	existing := &reference.State{ID: "s1", CountryID: "c1", Name: "Kyivska", Code: ptr("KY")}
	mockRepo.EXPECT().GetByID(gomock.Any(), "s1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, s *reference.State) error {
		assert.Equal(t, "Kyivska oblast", s.Name)
		assert.Equal(t, ptr("KY"), s.Code) // unchanged
		assert.Equal(t, ptr("frontline"), s.ConflictZone)
		return nil
	})

	out, err := uc.Update(context.Background(), "s1", ucadmin.UpdateStateInput{
		Name:         ptr("Kyivska oblast"),
		ConflictZone: ptr("frontline"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Kyivska oblast", out.Name)
}

func TestStateUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockStateRepository(ctrl)
	uc := ucadmin.NewStateUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "s1").Return(nil)

	err := uc.Delete(context.Background(), "s1")
	require.NoError(t, err)
}

func TestStateUseCase_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockStateRepository(ctrl)
	uc := ucadmin.NewStateUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "nonexistent").Return(reference.ErrStateNotFound)

	err := uc.Delete(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, reference.ErrStateNotFound)
}
