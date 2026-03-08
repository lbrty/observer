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

func TestOfficeUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockOfficeRepository(ctrl)
	uc := ucadmin.NewOfficeUseCase(mockRepo)

	mockRepo.EXPECT().List(gomock.Any()).Return([]*reference.Office{
		{ID: "o1", Name: "Kyiv HQ", PlaceID: ptr("p1")},
		{ID: "o2", Name: "Odesa Branch"},
	}, nil)

	out, err := uc.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, "Kyiv HQ", out[0].Name)
	assert.Equal(t, ptr("p1"), out[0].PlaceID)
	assert.Nil(t, out[1].PlaceID)
}

func TestOfficeUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockOfficeRepository(ctrl)
	uc := ucadmin.NewOfficeUseCase(mockRepo)

	mockRepo.EXPECT().GetByID(gomock.Any(), "o1").Return(&reference.Office{
		ID: "o1", Name: "Kyiv HQ", PlaceID: ptr("p1"),
	}, nil)

	out, err := uc.Get(context.Background(), "o1")
	require.NoError(t, err)
	assert.Equal(t, "o1", out.ID)
	assert.Equal(t, "Kyiv HQ", out.Name)
}

func TestOfficeUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockOfficeRepository(ctrl)
	uc := ucadmin.NewOfficeUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, o *reference.Office) error {
			assert.NotEmpty(t, o.ID)
			assert.Equal(t, "Lviv Office", o.Name)
			assert.Equal(t, ptr("p2"), o.PlaceID)
			return nil
		})

	out, err := uc.Create(context.Background(), ucadmin.CreateOfficeInput{
		Name:    "Lviv Office",
		PlaceID: ptr("p2"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Lviv Office", out.Name)
}

func TestOfficeUseCase_Update_PartialUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockOfficeRepository(ctrl)
	uc := ucadmin.NewOfficeUseCase(mockRepo)

	existing := &reference.Office{ID: "o1", Name: "Kyiv HQ", PlaceID: ptr("p1")}
	mockRepo.EXPECT().GetByID(gomock.Any(), "o1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, o *reference.Office) error {
		assert.Equal(t, "Kyiv Main Office", o.Name)
		assert.Equal(t, ptr("p1"), o.PlaceID) // unchanged
		return nil
	})

	out, err := uc.Update(context.Background(), "o1", ucadmin.UpdateOfficeInput{
		Name: ptr("Kyiv Main Office"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Kyiv Main Office", out.Name)
}

func TestOfficeUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockOfficeRepository(ctrl)
	uc := ucadmin.NewOfficeUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "o1").Return(nil)

	err := uc.Delete(context.Background(), "o1")
	require.NoError(t, err)
}

func TestOfficeUseCase_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockOfficeRepository(ctrl)
	uc := ucadmin.NewOfficeUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "nonexistent").Return(reference.ErrOfficeNotFound)

	err := uc.Delete(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, reference.ErrOfficeNotFound)
}
