package admin_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	domainproject "github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/handler/admin"
	"github.com/lbrty/observer/internal/handler/handlertest"
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
	return ucadmin.NewProjectUseCase(d.projectRepo, d.permRepo, nil)
}

func newProjectHandler(d *projectTestDeps) *admin.ProjectHandler {
	return admin.NewProjectHandler(d.projectUseCase())
}

func setAdminAuth(c *gin.Context) {
	uid := handlertest.TestID()
	handlertest.SetAuthContext(c, uid)
	c.Set(string(middleware.CtxUserRole), string(user.RoleAdmin))
}

// --- List ---

func TestProjectHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newProjectTestDeps(ctrl)
	h := newProjectHandler(d)

	now := time.Now().UTC()
	ownerID := handlertest.TestID().String()
	d.projectRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*domainproject.Project{
		{ID: handlertest.TestID().String(), Name: "Test Project", OwnerID: ownerID, Status: domainproject.ProjectStatusActive, CreatedAt: now, UpdatedAt: now},
	}, 1, nil)

	c, w := handlertest.NewTestContext(http.MethodGet, "/admin/projects?page=1&per_page=10", nil)
	setAdminAuth(c)
	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	projects := resp["projects"].([]any)
	assert.Len(t, projects, 1)
	assert.Equal(t, float64(1), resp["total"])
}

// --- Get ---

func TestProjectHandler_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newProjectTestDeps(ctrl)
	h := newProjectHandler(d)

	projectID := handlertest.TestID().String()
	d.projectRepo.EXPECT().GetByID(gomock.Any(), projectID).Return(nil, domainproject.ErrProjectNotFound)

	c, w := handlertest.NewTestContextWithParams(http.MethodGet, "/admin/projects/"+projectID, nil, gin.Params{
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
	projectID := handlertest.TestID().String()
	ownerID := handlertest.TestID().String()
	d.projectRepo.EXPECT().GetByID(gomock.Any(), projectID).Return(&domainproject.Project{
		ID:        projectID,
		Name:      "Test Project",
		OwnerID:   ownerID,
		Status:    domainproject.ProjectStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodGet, "/admin/projects/"+projectID, nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	setAdminAuth(c)
	h.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.Equal(t, projectID, resp["id"])
	assert.Equal(t, "Test Project", resp["name"])
}

// --- Create ---

func TestProjectHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newProjectTestDeps(ctrl)
	h := newProjectHandler(d)

	c, w := handlertest.NewTestContext(http.MethodPost, "/admin/projects", map[string]any{})
	setAdminAuth(c)
	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProjectHandler_Create_NameExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newProjectTestDeps(ctrl)
	h := newProjectHandler(d)

	d.projectRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domainproject.ErrProjectNameExists)

	c, w := handlertest.NewTestContext(http.MethodPost, "/admin/projects", map[string]any{
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

	c, w := handlertest.NewTestContext(http.MethodPost, "/admin/projects", map[string]any{
		"name":        "New Project",
		"description": "A test project",
	})
	setAdminAuth(c)
	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, "New Project", resp["name"])
	assert.Equal(t, "active", resp["status"])
}

// --- Update ---

func TestProjectHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newProjectTestDeps(ctrl)
	h := newProjectHandler(d)

	projectID := handlertest.TestID().String()
	d.projectRepo.EXPECT().GetByID(gomock.Any(), projectID).Return(nil, domainproject.ErrProjectNotFound)

	c, w := handlertest.NewTestContextWithParams(http.MethodPatch, "/admin/projects/"+projectID, map[string]any{
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
	projectID := handlertest.TestID().String()
	ownerID := handlertest.TestID().String()
	existing := &domainproject.Project{
		ID:        projectID,
		Name:      "Old Name",
		OwnerID:   ownerID,
		Status:    domainproject.ProjectStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	d.projectRepo.EXPECT().GetByID(gomock.Any(), projectID).Return(existing, nil)
	d.projectRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodPatch, "/admin/projects/"+projectID, map[string]any{
		"name": "Updated Name",
	}, gin.Params{{Key: "project_id", Value: projectID}})
	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.Equal(t, projectID, resp["id"])
}
