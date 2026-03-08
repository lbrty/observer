package my_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucmy "github.com/lbrty/observer/internal/usecase/my"
)

func TestMyProjectsUseCase_Execute_Admin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permRepo := mock_repo.NewMockPermissionRepository(ctrl)
	projectRepo := mock_repo.NewMockProjectRepository(ctrl)

	now := time.Now()
	projects := []*project.Project{
		{ID: "p1", Name: "Alpha", Status: project.ProjectStatusActive, CreatedAt: now, UpdatedAt: now},
		{ID: "p2", Name: "Beta", Status: project.ProjectStatusActive, CreatedAt: now, UpdatedAt: now},
	}

	projectRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(projects, 2, nil)

	uc := ucmy.NewMyProjectsUseCase(permRepo, projectRepo)
	out, err := uc.Execute(context.Background(), "user-1", user.RoleAdmin)

	require.NoError(t, err)
	require.Len(t, out.Projects, 2)
	assert.Equal(t, "p1", out.Projects[0].ID)
	assert.Equal(t, "owner", out.Projects[0].Role)
	assert.Equal(t, "p2", out.Projects[1].ID)
	assert.Equal(t, "owner", out.Projects[1].Role)
}

func TestMyProjectsUseCase_Execute_Staff(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permRepo := mock_repo.NewMockPermissionRepository(ctrl)
	projectRepo := mock_repo.NewMockProjectRepository(ctrl)

	now := time.Now()
	perms := []*project.ProjectPermission{
		{
			ID:               "perm-1",
			ProjectID:        "p1",
			UserID:           "user-1",
			Role:             project.ProjectRoleConsultant,
			CanViewContact:   true,
			CanViewPersonal:  false,
			CanViewDocuments: true,
		},
	}
	activeProject := &project.Project{
		ID:     "p1",
		Name:   "Alpha",
		Status: project.ProjectStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	permRepo.EXPECT().ListByUserID(gomock.Any(), "user-1").Return(perms, nil)
	projectRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(activeProject, nil)

	uc := ucmy.NewMyProjectsUseCase(permRepo, projectRepo)
	out, err := uc.Execute(context.Background(), "user-1", user.RoleStaff)

	require.NoError(t, err)
	require.Len(t, out.Projects, 1)
	assert.Equal(t, "p1", out.Projects[0].ID)
	assert.Equal(t, "consultant", out.Projects[0].Role)
	assert.True(t, out.Projects[0].CanViewContact)
	assert.False(t, out.Projects[0].CanViewPersonal)
	assert.True(t, out.Projects[0].CanViewDocuments)
}

func TestMyProjectsUseCase_Execute_Staff_SkipsInactive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permRepo := mock_repo.NewMockPermissionRepository(ctrl)
	projectRepo := mock_repo.NewMockProjectRepository(ctrl)

	now := time.Now()
	perms := []*project.ProjectPermission{
		{ID: "perm-1", ProjectID: "p1", UserID: "user-1", Role: project.ProjectRoleViewer},
		{ID: "perm-2", ProjectID: "p2", UserID: "user-1", Role: project.ProjectRoleManager},
	}

	permRepo.EXPECT().ListByUserID(gomock.Any(), "user-1").Return(perms, nil)
	projectRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&project.Project{
		ID: "p1", Name: "Archived", Status: project.ProjectStatusArchived, CreatedAt: now, UpdatedAt: now,
	}, nil)
	projectRepo.EXPECT().GetByID(gomock.Any(), "p2").Return(&project.Project{
		ID: "p2", Name: "Active", Status: project.ProjectStatusActive, CreatedAt: now, UpdatedAt: now,
	}, nil)

	uc := ucmy.NewMyProjectsUseCase(permRepo, projectRepo)
	out, err := uc.Execute(context.Background(), "user-1", user.RoleStaff)

	require.NoError(t, err)
	require.Len(t, out.Projects, 1)
	assert.Equal(t, "p2", out.Projects[0].ID)
	assert.Equal(t, "manager", out.Projects[0].Role)
}

func TestMyProjectsUseCase_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permRepo := mock_repo.NewMockPermissionRepository(ctrl)
	projectRepo := mock_repo.NewMockProjectRepository(ctrl)

	repoErr := errors.New("db connection lost")
	projectRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, 0, repoErr)

	uc := ucmy.NewMyProjectsUseCase(permRepo, projectRepo)
	_, err := uc.Execute(context.Background(), "user-1", user.RoleAdmin)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}
