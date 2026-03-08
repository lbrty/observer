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
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
)

func TestPermissionUseCase_Assign_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermRepo := mock_repo.NewMockPermissionRepository(ctrl)
	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	auditRepo := mock_repo.NewMockAuditLogRepository(ctrl)
	auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	uc := ucadmin.NewPermissionUseCase(mockPermRepo, mockUserRepo, auditUC)

	ctx := context.Background()

	mockPermRepo.EXPECT().
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

	out, err := uc.Assign(ctx, "proj-1", ucadmin.AssignPermissionInput{
		UserID:         "user-1",
		Role:           "manager",
		CanViewContact: true,
	})
	require.NoError(t, err)
	assert.Equal(t, "proj-1", out.ProjectID)
	assert.Equal(t, "user-1", out.UserID)
	assert.Equal(t, "manager", out.Role)
}

func TestPermissionUseCase_Assign_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermRepo := mock_repo.NewMockPermissionRepository(ctrl)
	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	auditRepo := mock_repo.NewMockAuditLogRepository(ctrl)
	auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	uc := ucadmin.NewPermissionUseCase(mockPermRepo, mockUserRepo, auditUC)

	_, err := uc.Assign(context.Background(), "proj-1", ucadmin.AssignPermissionInput{
		UserID: "user-1",
		Role:   "superadmin",
	})
	assert.ErrorIs(t, err, project.ErrInvalidProjectRole)
}

func TestPermissionUseCase_Assign_DuplicateDetection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermRepo := mock_repo.NewMockPermissionRepository(ctrl)
	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	auditRepo := mock_repo.NewMockAuditLogRepository(ctrl)
	auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	auditUC := ucaudit.NewAuditUseCase(auditRepo)
	uc := ucadmin.NewPermissionUseCase(mockPermRepo, mockUserRepo, auditUC)

	mockPermRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(project.ErrPermissionExists)

	_, err := uc.Assign(context.Background(), "proj-1", ucadmin.AssignPermissionInput{
		UserID: "user-1",
		Role:   "viewer",
	})
	assert.ErrorIs(t, err, project.ErrPermissionExists)
}
