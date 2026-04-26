package my_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	domainproject "github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/handler/handlertest"
	"github.com/lbrty/observer/internal/handler/my"
	"github.com/lbrty/observer/internal/middleware"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucmy "github.com/lbrty/observer/internal/usecase/my"
)

func newMyHandler(ctrl *gomock.Controller) (*my.MyHandler, *repomock.MockPermissionRepository, *repomock.MockProjectRepository) {
	permRepo := repomock.NewMockPermissionRepository(ctrl)
	projectRepo := repomock.NewMockProjectRepository(ctrl)
	uc := ucmy.NewMyProjectsUseCase(permRepo, projectRepo)
	return my.NewMyHandler(uc), permRepo, projectRepo
}

func TestMyHandler_Projects_NoAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, _ := newMyHandler(ctrl)

	c, w := handlertest.NewTestContext(http.MethodGet, "/my/projects", nil)
	h.Projects(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMyHandler_Projects_Admin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, projectRepo := newMyHandler(ctrl)

	userID := handlertest.TestID()
	now := time.Now().UTC()
	projectID := handlertest.TestID().String()

	projectRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*domainproject.Project{
		{ID: projectID, Name: "Test Project", Status: domainproject.ProjectStatusActive, OwnerID: userID.String(), CreatedAt: now, UpdatedAt: now},
	}, 1, nil)

	c, w := handlertest.NewTestContext(http.MethodGet, "/my/projects", nil)
	handlertest.SetAuthContext(c, userID)
	c.Set(string(middleware.CtxUserRole), string(user.RoleAdmin))
	h.Projects(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	projects := resp["projects"].([]any)
	assert.Len(t, projects, 1)
	p := projects[0].(map[string]any)
	assert.Equal(t, "Test Project", p["name"])
	assert.Equal(t, "owner", p["role"])
}
