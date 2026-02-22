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

func TestCountryUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCountryRepository(ctrl)
	uc := ucadmin.NewCountryUseCase(mockRepo)

	mockRepo.EXPECT().List(gomock.Any()).Return([]*reference.Country{
		{ID: "c1", Name: "Ukraine", Code: "UA"},
		{ID: "c2", Name: "Kyrgyzstan", Code: "KG"},
	}, nil)

	out, err := uc.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, "Ukraine", out[0].Name)
	assert.Equal(t, "UA", out[0].Code)
}

func TestCountryUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCountryRepository(ctrl)
	uc := ucadmin.NewCountryUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, c *reference.Country) error {
			assert.NotEmpty(t, c.ID)
			assert.Equal(t, "Moldova", c.Name)
			assert.Equal(t, "MD", c.Code)
			return nil
		})

	out, err := uc.Create(context.Background(), ucadmin.CreateCountryInput{
		Name: "Moldova",
		Code: "MD",
	})
	require.NoError(t, err)
	assert.Equal(t, "Moldova", out.Name)
	assert.Equal(t, "MD", out.Code)
}

func TestCountryUseCase_Create_DuplicateCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCountryRepository(ctrl)
	uc := ucadmin.NewCountryUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(reference.ErrCountryCodeExists)

	_, err := uc.Create(context.Background(), ucadmin.CreateCountryInput{
		Name: "Ukraine",
		Code: "UA",
	})
	assert.ErrorIs(t, err, reference.ErrCountryCodeExists)
}

func TestCountryUseCase_Update_PartialUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCountryRepository(ctrl)
	uc := ucadmin.NewCountryUseCase(mockRepo)

	existing := &reference.Country{ID: "c1", Name: "Ukraina", Code: "UA"}
	mockRepo.EXPECT().GetByID(gomock.Any(), "c1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, c *reference.Country) error {
		assert.Equal(t, "Ukraine", c.Name)
		assert.Equal(t, "UA", c.Code) // unchanged
		return nil
	})

	out, err := uc.Update(context.Background(), "c1", ucadmin.UpdateCountryInput{
		Name: ptr("Ukraine"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Ukraine", out.Name)
}

func TestCountryUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCountryRepository(ctrl)
	uc := ucadmin.NewCountryUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "c1").Return(nil)

	err := uc.Delete(context.Background(), "c1")
	require.NoError(t, err)
}

func TestCountryUseCase_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCountryRepository(ctrl)
	uc := ucadmin.NewCountryUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "nonexistent").Return(reference.ErrCountryNotFound)

	err := uc.Delete(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, reference.ErrCountryNotFound)
}
