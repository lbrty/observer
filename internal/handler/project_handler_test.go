package handler_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/middleware"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

type projectTestDeps struct {
	ctrl        *gomock.Controller
	projectRepo *repomock.MockProjectRepository
	permRepo    *repomock.MockPermissionRepository
}

func newProjectTestDeps(ctrl *gomock.Controller) *projectTestDeps {
	return &projectTestDeps{
		ctrl:        ctrl,
		projectRepo: repomock.NewMockProjectRepository(ctrl),
		permRepo:    repomock.NewMockPermissionRepository(ctrl),
	}
}

func (d *projectTestDeps) projectUseCase() *ucadmin.ProjectUseCase {
	return ucadmin.NewProjectUseCase(d.projectRepo, d.permRepo)
}

func newProjectHandler(d *projectTestDeps) *handler.ProjectHandler {
	return handler.NewProjectHandler(d.projectUseCase())
}

func setAdminAuth(c *gin.Context) {
	uid := testID()
	setAuthContext(c, uid)
	c.Set(string(middleware.CtxUserRole), string(user.RoleAdmin))
}

// --- List ---

func TestProjectHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newProjectTestDeps(ctrl)
	h := newProjectHandler(d)

	now := time.Now().UTC()
	ownerID := testID().String()
	d.projectRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*project.Project{
		{ID: testID().String(), Name: "Test Project", OwnerID: ownerID, Status: project.ProjectStatusActive, CreatedAt: now, UpdatedAt: now},
	}, 1, nil)

	c, w := newTestContext(http.MethodGet, "/admin/projects?page=1&per_page=10", nil)
	setAdminAuth(c)
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	projects := resp["projects"].([]any)
	assert.Len(t, projects, 1)
	assert.Equal(t, float64(1), resp["total"])
}

func TestProjectHandler_List_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newProjectTestDeps(ctrl)
	h := newProjectHandler(d)

	d.projectRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, 0, fmt.Errorf("db error"))

	c, w := newTestContext(http.MethodGet, "/admin/projects", nil)
	setAdminAuth(c)
	h.List(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- Get ---

func TestProjectHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newProjectTestDeps(ctrl)
	h := newProjectHandler(d)

	projectID := testID().String()
	d.projectRepo.EXPECT().GetByID(gomock.Any(), projectID).Return(nil, project.ErrProjectNotFound)

	c, w := newTestContextWithParams(http.MethodGet, "/admin/projects/"+projectID, nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	setAdminAuth(c)
	h.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProjectHandler_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newProjectTestDeps(ctrl)
	h := newProjectHandler(d)

	now := time.Now().UTC()
	projectID := testID().String()
	ownerID := testID().String()
	d.projectRepo.EXPECT().GetByID(gomock.Any(), projectID).Return(&project.Project{
		ID:      projectID,
		Name:    "Test Project",
		OwnerID: ownerID,
		Status:  project.ProjectStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/admin/projects/"+projectID, nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	setAdminAuth(c)
	h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, projectID, resp["id"])
	assert.Equal(t, "Test Project", resp["name"])
}

// --- Create ---

func TestProjectHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newProjectTestDeps(ctrl)
	h := newProjectHandler(d)

	c, w := newTestContext(http.MethodPost, "/admin/projects", map[string]any{})
	setAdminAuth(c)
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProjectHandler_Create_NameExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newProjectTestDeps(ctrl)
	h := newProjectHandler(d)

	d.projectRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(project.ErrProjectNameExists)

	c, w := newTestContext(http.MethodPost, "/admin/projects", map[string]any{
		"name": "Existing Project",
	})
	setAdminAuth(c)
	h.Create(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestProjectHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newProjectTestDeps(ctrl)
	h := newProjectHandler(d)

	d.projectRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContext(http.MethodPost, "/admin/projects", map[string]any{
		"name":        "New Project",
		"description": "A test project",
	})
	setAdminAuth(c)
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, "New Project", resp["name"])
	assert.Equal(t, "active", resp["status"])
}

// --- Update ---

func TestProjectHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newProjectTestDeps(ctrl)
	h := newProjectHandler(d)

	projectID := testID().String()
	d.projectRepo.EXPECT().GetByID(gomock.Any(), projectID).Return(nil, project.ErrProjectNotFound)

	c, w := newTestContextWithParams(http.MethodPatch, "/admin/projects/"+projectID, map[string]any{
		"name": "Updated",
	}, gin.Params{{Key: "project_id", Value: projectID}})
	h.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProjectHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newProjectTestDeps(ctrl)
	h := newProjectHandler(d)

	now := time.Now().UTC()
	projectID := testID().String()
	ownerID := testID().String()
	existing := &project.Project{
		ID:        projectID,
		Name:      "Old Name",
		OwnerID:   ownerID,
		Status:    project.ProjectStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	d.projectRepo.EXPECT().GetByID(gomock.Any(), projectID).Return(existing, nil)
	d.projectRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPatch, "/admin/projects/"+projectID, map[string]any{
		"name": "Updated Name",
	}, gin.Params{{Key: "project_id", Value: projectID}})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, projectID, resp["id"])
}
