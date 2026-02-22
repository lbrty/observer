package admin_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/project"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

func TestAssignPermissionUseCase_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPermissionRepository(ctrl)
	uc := ucadmin.NewAssignPermissionUseCase(mockRepo)

	ctx := context.Background()

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, p *project.ProjectPermission) error {
			assert.NotEmpty(t, p.ID)
			assert.Equal(t, "proj-1", p.ProjectID)
			assert.Equal(t, "user-1", p.UserID)
			assert.Equal(t, project.ProjectRoleManager, p.Role)
			assert.True(t, p.CanViewContact)
			assert.False(t, p.CanViewPersonal)
			return nil
		})

	out, err := uc.Execute(ctx, "proj-1", ucadmin.AssignPermissionInput{
		UserID:         "user-1",
		Role:           "manager",
		CanViewContact: true,
	})
	require.NoError(t, err)
	assert.Equal(t, "proj-1", out.ProjectID)
	assert.Equal(t, "user-1", out.UserID)
	assert.Equal(t, "manager", out.Role)
}

func TestAssignPermissionUseCase_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPermissionRepository(ctrl)
	uc := ucadmin.NewAssignPermissionUseCase(mockRepo)

	_, err := uc.Execute(context.Background(), "proj-1", ucadmin.AssignPermissionInput{
		UserID: "user-1",
		Role:   "superadmin",
	})
	assert.ErrorIs(t, err, project.ErrInvalidProjectRole)
}

func TestAssignPermissionUseCase_DuplicateDetection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockPermissionRepository(ctrl)
	uc := ucadmin.NewAssignPermissionUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(project.ErrPermissionExists)

	_, err := uc.Execute(context.Background(), "proj-1", ucadmin.AssignPermissionInput{
		UserID: "user-1",
		Role:   "viewer",
	})
	assert.ErrorIs(t, err, project.ErrPermissionExists)
}
