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

func TestProjectUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockProjectRepository(ctrl)
	uc := ucadmin.NewProjectUseCase(mockRepo)

	mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*project.Project{
		{ID: "p1", Name: "Project A", OwnerID: "u1", Status: project.ProjectStatusActive},
		{ID: "p2", Name: "Project B", OwnerID: "u1", Status: project.ProjectStatusArchived},
	}, 2, nil)

	out, err := uc.List(context.Background(), ucadmin.ListProjectsInput{Page: 1, PerPage: 20})
	require.NoError(t, err)
	assert.Len(t, out.Projects, 2)
	assert.Equal(t, 2, out.Total)
	assert.Equal(t, "Project A", out.Projects[0].Name)
	assert.Equal(t, "active", out.Projects[0].Status)
}

func TestProjectUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockProjectRepository(ctrl)
	uc := ucadmin.NewProjectUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, p *project.Project) error {
			assert.NotEmpty(t, p.ID)
			assert.Equal(t, "New Project", p.Name)
			assert.Equal(t, "owner1", p.OwnerID)
			assert.Equal(t, project.ProjectStatusActive, p.Status)
			return nil
		})

	out, err := uc.Create(context.Background(), "owner1", ucadmin.CreateProjectInput{
		Name: "New Project",
	})
	require.NoError(t, err)
	assert.Equal(t, "New Project", out.Name)
	assert.Equal(t, "owner1", out.OwnerID)
}

func TestProjectUseCase_Create_DuplicateName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockProjectRepository(ctrl)
	uc := ucadmin.NewProjectUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(project.ErrProjectNameExists)

	_, err := uc.Create(context.Background(), "owner1", ucadmin.CreateProjectInput{
		Name: "Existing",
	})
	assert.ErrorIs(t, err, project.ErrProjectNameExists)
}

func TestProjectUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockProjectRepository(ctrl)
	uc := ucadmin.NewProjectUseCase(mockRepo)

	mockRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(&project.Project{
		ID: "p1", Name: "Project A", OwnerID: "u1", Status: project.ProjectStatusActive,
	}, nil)

	out, err := uc.Get(context.Background(), "p1")
	require.NoError(t, err)
	assert.Equal(t, "p1", out.ID)
	assert.Equal(t, "Project A", out.Name)
}

func TestProjectUseCase_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockProjectRepository(ctrl)
	uc := ucadmin.NewProjectUseCase(mockRepo)

	mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, project.ErrProjectNotFound)

	_, err := uc.Get(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, project.ErrProjectNotFound)
}

func TestProjectUseCase_Update_Archive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockProjectRepository(ctrl)
	uc := ucadmin.NewProjectUseCase(mockRepo)

	existing := &project.Project{
		ID: "p1", Name: "Project A", OwnerID: "u1", Status: project.ProjectStatusActive,
	}
	mockRepo.EXPECT().GetByID(gomock.Any(), "p1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, p *project.Project) error {
		assert.Equal(t, project.ProjectStatusArchived, p.Status)
		assert.Equal(t, "Project A", p.Name) // unchanged
		return nil
	})

	status := "archived"
	out, err := uc.Update(context.Background(), "p1", ucadmin.UpdateProjectInput{
		Status: &status,
	})
	require.NoError(t, err)
	assert.Equal(t, "archived", out.Status)
}
