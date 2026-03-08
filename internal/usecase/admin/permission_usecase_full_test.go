package admin_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

func TestPermissionUseCase_List_Admin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermRepo := mock_repo.NewMockPermissionRepository(ctrl)
	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	uc := ucadmin.NewPermissionUseCase(mockPermRepo, mockUserRepo)

	ctx := context.Background()
	userID := ulid.Make()
	userIDStr := userID.String()

	mockPermRepo.EXPECT().List(ctx, "proj-1").Return([]*project.ProjectPermission{
		{ID: "perm-1", ProjectID: "proj-1", UserID: userIDStr, Role: project.ProjectRoleManager, CanViewContact: true},
	}, nil)
	mockUserRepo.EXPECT().GetByIDs(ctx, []ulid.ULID{userID}).Return([]*user.User{
		{ID: userID, FirstName: "Alice", LastName: "Smith", Email: "alice@example.com", Role: user.RoleStaff},
	}, nil)

	out, err := uc.List(ctx, "proj-1", "caller-admin", user.RoleAdmin)
	require.NoError(t, err)
	assert.Len(t, out, 1)
	assert.Equal(t, "perm-1", out[0].ID)
	assert.Equal(t, "Alice", out[0].UserFirstName)
	assert.Equal(t, "manager", out[0].Role)
}

func TestPermissionUseCase_List_NonMember_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermRepo := mock_repo.NewMockPermissionRepository(ctrl)
	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	uc := ucadmin.NewPermissionUseCase(mockPermRepo, mockUserRepo)

	ctx := context.Background()

	// Caller is a consultant, not a member of proj-1
	mockPermRepo.EXPECT().ListByUserID(ctx, "caller-1").Return([]*project.ProjectPermission{
		{ID: "perm-99", ProjectID: "proj-other", UserID: "caller-1", Role: project.ProjectRoleConsultant},
	}, nil)

	out, err := uc.List(ctx, "proj-1", "caller-1", user.RoleConsultant)
	require.NoError(t, err)
	assert.Empty(t, out)
}

func TestPermissionUseCase_List_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermRepo := mock_repo.NewMockPermissionRepository(ctrl)
	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	uc := ucadmin.NewPermissionUseCase(mockPermRepo, mockUserRepo)

	ctx := context.Background()

	mockPermRepo.EXPECT().List(ctx, "proj-1").Return(nil, fmt.Errorf("db connection lost"))

	_, err := uc.List(ctx, "proj-1", "caller-admin", user.RoleAdmin)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db connection lost")
}

func TestPermissionUseCase_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermRepo := mock_repo.NewMockPermissionRepository(ctrl)
	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	uc := ucadmin.NewPermissionUseCase(mockPermRepo, mockUserRepo)

	ctx := context.Background()

	existing := &project.ProjectPermission{
		ID:              "perm-1",
		ProjectID:       "proj-1",
		UserID:          "user-1",
		Role:            project.ProjectRoleViewer,
		CanViewContact:  false,
		CanViewPersonal: false,
	}
	mockPermRepo.EXPECT().GetByID(ctx, "perm-1").Return(existing, nil)
	mockPermRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, p *project.ProjectPermission) error {
		assert.Equal(t, project.ProjectRoleManager, p.Role)
		assert.True(t, p.CanViewContact)
		assert.False(t, p.CanViewPersonal) // unchanged
		return nil
	})

	out, err := uc.Update(ctx, "perm-1", ucadmin.UpdatePermissionInput{
		Role:           ptr("manager"),
		CanViewContact: ptr(true),
	})
	require.NoError(t, err)
	assert.Equal(t, "manager", out.Role)
	assert.True(t, out.CanViewContact)
}

func TestPermissionUseCase_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermRepo := mock_repo.NewMockPermissionRepository(ctrl)
	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	uc := ucadmin.NewPermissionUseCase(mockPermRepo, mockUserRepo)

	mockPermRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, project.ErrPermissionNotFound)

	_, err := uc.Update(context.Background(), "nonexistent", ucadmin.UpdatePermissionInput{
		Role: ptr("manager"),
	})
	assert.ErrorIs(t, err, project.ErrPermissionNotFound)
}

func TestPermissionUseCase_Update_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermRepo := mock_repo.NewMockPermissionRepository(ctrl)
	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	uc := ucadmin.NewPermissionUseCase(mockPermRepo, mockUserRepo)

	existing := &project.ProjectPermission{
		ID:        "perm-1",
		ProjectID: "proj-1",
		UserID:    "user-1",
		Role:      project.ProjectRoleViewer,
	}
	mockPermRepo.EXPECT().GetByID(gomock.Any(), "perm-1").Return(existing, nil)

	_, err := uc.Update(context.Background(), "perm-1", ucadmin.UpdatePermissionInput{
		Role: ptr("superadmin"),
	})
	assert.ErrorIs(t, err, project.ErrInvalidProjectRole)
}

func TestPermissionUseCase_Revoke_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermRepo := mock_repo.NewMockPermissionRepository(ctrl)
	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	uc := ucadmin.NewPermissionUseCase(mockPermRepo, mockUserRepo)

	mockPermRepo.EXPECT().Delete(gomock.Any(), "perm-1").Return(nil)

	err := uc.Revoke(context.Background(), "perm-1")
	require.NoError(t, err)
}

func TestPermissionUseCase_Revoke_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPermRepo := mock_repo.NewMockPermissionRepository(ctrl)
	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	uc := ucadmin.NewPermissionUseCase(mockPermRepo, mockUserRepo)

	mockPermRepo.EXPECT().Delete(gomock.Any(), "nonexistent").Return(project.ErrPermissionNotFound)

	err := uc.Revoke(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, project.ErrPermissionNotFound)
}
