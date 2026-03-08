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

func TestCategoryUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCategoryRepository(ctrl)
	uc := ucadmin.NewCategoryUseCase(mockRepo)

	mockRepo.EXPECT().List(gomock.Any()).Return([]*reference.Category{
		{ID: "cat1", Name: "IDP", Description: ptr("Internally displaced person")},
		{ID: "cat2", Name: "Refugee"},
	}, nil)

	out, err := uc.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, "IDP", out[0].Name)
	assert.Equal(t, ptr("Internally displaced person"), out[0].Description)
	assert.Nil(t, out[1].Description)
}

func TestCategoryUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCategoryRepository(ctrl)
	uc := ucadmin.NewCategoryUseCase(mockRepo)

	mockRepo.EXPECT().GetByID(gomock.Any(), "cat1").Return(&reference.Category{
		ID: "cat1", Name: "IDP", Description: ptr("Internally displaced person"),
	}, nil)

	out, err := uc.Get(context.Background(), "cat1")
	require.NoError(t, err)
	assert.Equal(t, "cat1", out.ID)
	assert.Equal(t, "IDP", out.Name)
}

func TestCategoryUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCategoryRepository(ctrl)
	uc := ucadmin.NewCategoryUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, c *reference.Category) error {
			assert.NotEmpty(t, c.ID)
			assert.Equal(t, "Veteran", c.Name)
			assert.Equal(t, ptr("War veteran"), c.Description)
			return nil
		})

	out, err := uc.Create(context.Background(), ucadmin.CreateCategoryInput{
		Name:        "Veteran",
		Description: ptr("War veteran"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Veteran", out.Name)
}

func TestCategoryUseCase_Create_DuplicateName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCategoryRepository(ctrl)
	uc := ucadmin.NewCategoryUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(reference.ErrCategoryNameExists)

	_, err := uc.Create(context.Background(), ucadmin.CreateCategoryInput{
		Name: "IDP",
	})
	assert.ErrorIs(t, err, reference.ErrCategoryNameExists)
}

func TestCategoryUseCase_Update_PartialUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCategoryRepository(ctrl)
	uc := ucadmin.NewCategoryUseCase(mockRepo)

	existing := &reference.Category{ID: "cat1", Name: "IDP", Description: ptr("old desc")}
	mockRepo.EXPECT().GetByID(gomock.Any(), "cat1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, c *reference.Category) error {
		assert.Equal(t, "Internally Displaced Person", c.Name)
		assert.Equal(t, ptr("old desc"), c.Description) // unchanged
		return nil
	})

	out, err := uc.Update(context.Background(), "cat1", ucadmin.UpdateCategoryInput{
		Name: ptr("Internally Displaced Person"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Internally Displaced Person", out.Name)
}

func TestCategoryUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCategoryRepository(ctrl)
	uc := ucadmin.NewCategoryUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "cat1").Return(nil)

	err := uc.Delete(context.Background(), "cat1")
	require.NoError(t, err)
}

func TestCategoryUseCase_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCategoryRepository(ctrl)
	uc := ucadmin.NewCategoryUseCase(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "nonexistent").Return(reference.ErrCategoryNotFound)

	err := uc.Delete(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, reference.ErrCategoryNotFound)
}
