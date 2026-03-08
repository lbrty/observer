package handler_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/middleware"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucmy "github.com/lbrty/observer/internal/usecase/my"
)

func newMyHandler(ctrl *gomock.Controller) (*handler.MyHandler, *repomock.MockPermissionRepository, *repomock.MockProjectRepository) {
	permRepo := repomock.NewMockPermissionRepository(ctrl)
	projectRepo := repomock.NewMockProjectRepository(ctrl)
	uc := ucmy.NewMyProjectsUseCase(permRepo, projectRepo)
	return handler.NewMyHandler(uc), permRepo, projectRepo
}

func TestMyHandler_Projects_NoAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, _ := newMyHandler(ctrl)

	c, w := newTestContext(http.MethodGet, "/my/projects", nil)
	h.Projects(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMyHandler_Projects_Admin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, projectRepo := newMyHandler(ctrl)

	userID := testID()
	now := time.Now().UTC()
	projectID := testID().String()

	projectRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*project.Project{
		{ID: projectID, Name: "Test Project", Status: project.ProjectStatusActive, OwnerID: userID.String(), CreatedAt: now, UpdatedAt: now},
	}, 1, nil)

	c, w := newTestContext(http.MethodGet, "/my/projects", nil)
	setAuthContext(c, userID)
	c.Set(string(middleware.CtxUserRole), string(user.RoleAdmin))
	h.Projects(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	projects := resp["projects"].([]any)
	assert.Len(t, projects, 1)
	p := projects[0].(map[string]any)
	assert.Equal(t, "Test Project", p["name"])
	assert.Equal(t, "owner", p["role"])
}

func TestMyHandler_Projects_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, _, projectRepo := newMyHandler(ctrl)

	userID := testID()
	projectRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, 0, fmt.Errorf("db error"))

	c, w := newTestContext(http.MethodGet, "/my/projects", nil)
	setAuthContext(c, userID)
	c.Set(string(middleware.CtxUserRole), string(user.RoleAdmin))
	h.Projects(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
