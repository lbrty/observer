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

func TestPlaceUseCase_ListAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPlaceRepository(ctrl)
	uc := ucadmin.NewPlaceUseCase(mockRepo)

	mockRepo.EXPECT().ListAll(gomock.Any()).Return([]*reference.Place{
		{ID: "p1", StateID: "s1", Name: "Kyiv", Lat: ptr(50.45), Lon: ptr(30.52)},
		{ID: "p2", StateID: "s1", Name: "Brovary"},
	}, nil)

	out, err := uc.ListAll(context.Background())
	require.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, "Kyiv", out[0].Name)
	assert.Equal(t, ptr(50.45), out[0].Lat)
	assert.Nil(t, out[1].Lat)
}

func TestPlaceUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPlaceRepository(ctrl)
	uc := ucadmin.NewPlaceUseCase(mockRepo)

	mockRepo.EXPECT().List(gomock.Any(), "s1").Return([]*reference.Place{
		{ID: "p1", StateID: "s1", Name: "Kyiv"},
	}, nil)

	out, err := uc.List(context.Background(), "s1")
	require.NoError(t, err)
	assert.Len(t, out, 1)
	assert.Equal(t, "s1", out[0].StateID)
}

func TestPlaceUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPlaceRepository(ctrl)
	uc := ucadmin.NewPlaceUseCase(mockRepo)

	mockRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&reference.Place{
		ID: "p1", StateID: "s1", Name: "Kyiv", Lat: ptr(50.45), Lon: ptr(30.52),
	}, nil)

	out, err := uc.Get(context.Background(), "p1")
	require.NoError(t, err)
	assert.Equal(t, "p1", out.ID)
	assert.Equal(t, "Kyiv", out.Name)
	assert.Equal(t, ptr(50.45), out.Lat)
}

func TestPlaceUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPlaceRepository(ctrl)
	uc := ucadmin.NewPlaceUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, p *reference.Place) error {
			assert.NotEmpty(t, p.ID)
			assert.Equal(t, "s1", p.StateID)
			assert.Equal(t, "Odesa", p.Name)
			assert.Equal(t, ptr(46.47), p.Lat)
			assert.Equal(t, ptr(30.73), p.Lon)
			return nil
		})

	out, err := uc.Create(context.Background(), "s1", ucadmin.CreatePlaceInput{
		Name: "Odesa",
		Lat:  ptr(46.47),
		Lon:  ptr(30.73),
	})
	require.NoError(t, err)
	assert.Equal(t, "Odesa", out.Name)
	assert.Equal(t, "s1", out.StateID)
}

func TestPlaceUseCase_Update_PartialUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPlaceRepository(ctrl)
	uc := ucadmin.NewPlaceUseCase(mockRepo)

	existing := &reference.Place{ID: "p1", StateID: "s1", Name: "Kyv", Lat: ptr(50.45), Lon: ptr(30.52)}
	mockRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, p *reference.Place) error {
		assert.Equal(t, "Kyiv", p.Name)
		assert.Equal(t, ptr(50.45), p.Lat) // unchanged
		return nil
	})

	out, err := uc.Update(context.Background(), "p1", ucadmin.UpdatePlaceInput{
		Name: ptr("Kyiv"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Kyiv", out.Name)
}

func TestPlaceUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPlaceRepository(ctrl)
	uc := ucadmin.NewPlaceUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "p1").Return(nil)

	err := uc.Delete(context.Background(), "p1")
	require.NoError(t, err)
}

func TestPlaceUseCase_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPlaceRepository(ctrl)
	uc := ucadmin.NewPlaceUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "nonexistent").Return(reference.ErrPlaceNotFound)

	err := uc.Delete(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, reference.ErrPlaceNotFound)
}
