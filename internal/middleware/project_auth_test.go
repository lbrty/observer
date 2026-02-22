package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/middleware"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
)

func newProjectRouter(
	t *testing.T,
	userID ulid.ULID,
	role string,
	action project.Action,
	mockSetup func(*mock_repo.MockPermissionLoader),
) (*httptest.ResponseRecorder, *gin.Engine) {
	t.Helper()
	ctrl := gomock.NewController(t)

	mockPerm := mock_repo.NewMockPermissionLoader(ctrl)
	if mockSetup != nil {
		mockSetup(mockPerm)
	}

	paMW := middleware.NewProjectAuthMiddleware(mockPerm)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.GET("/projects/:project_id/people", func(c *gin.Context) {
		setupAuthContext(c, userID, role)
		c.Next()
	}, paMW.RequireProjectRole(action), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"project_role":       c.GetString(string(middleware.CtxProjectRole)),
			"can_view_contact":   c.GetBool(string(middleware.CtxCanViewContact)),
			"can_view_personal":  c.GetBool(string(middleware.CtxCanViewPersonal)),
			"can_view_documents": c.GetBool(string(middleware.CtxCanViewDocuments)),
		})
	})

	return w, r
}

func TestRequireProjectRole_AdminBypass(t *testing.T) {
	uid := ulid.Make()
	w, r := newProjectRouter(t, uid, "admin", project.ActionDelete, nil)

	req := httptest.NewRequest(http.MethodGet, "/projects/01PROJ000000000000000001/people", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"project_role":"owner"`)
	assert.Contains(t, w.Body.String(), `"can_view_contact":true`)
	assert.Contains(t, w.Body.String(), `"can_view_personal":true`)
	assert.Contains(t, w.Body.String(), `"can_view_documents":true`)
}

func TestRequireProjectRole_OwnerBypass(t *testing.T) {
	uid := ulid.Make()
	projID := "01PROJ000000000000000002"

	w, r := newProjectRouter(t, uid, "staff", project.ActionDelete, func(m *mock_repo.MockPermissionLoader) {
		m.EXPECT().IsProjectOwner(gomock.Any(), uid, projID).Return(true, nil)
	})

	req := httptest.NewRequest(http.MethodGet, "/projects/"+projID+"/people", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"project_role":"owner"`)
}

func TestRequireProjectRole_Sufficient(t *testing.T) {
	uid := ulid.Make()
	projID := "01PROJ000000000000000003"

	w, r := newProjectRouter(t, uid, "staff", project.ActionCreate, func(m *mock_repo.MockPermissionLoader) {
		m.EXPECT().IsProjectOwner(gomock.Any(), uid, projID).Return(false, nil)
		m.EXPECT().GetPermission(gomock.Any(), uid, projID).Return(&project.Permission{
			Role:             project.ProjectRoleConsultant,
			CanViewContact:   true,
			CanViewPersonal:  false,
			CanViewDocuments: true,
		}, nil)
	})

	req := httptest.NewRequest(http.MethodGet, "/projects/"+projID+"/people", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"project_role":"consultant"`)
	assert.Contains(t, w.Body.String(), `"can_view_contact":true`)
	assert.Contains(t, w.Body.String(), `"can_view_personal":false`)
	assert.Contains(t, w.Body.String(), `"can_view_documents":true`)
}

func TestRequireProjectRole_Insufficient(t *testing.T) {
	uid := ulid.Make()
	projID := "01PROJ000000000000000004"

	w, r := newProjectRouter(t, uid, "staff", project.ActionCreate, func(m *mock_repo.MockPermissionLoader) {
		m.EXPECT().IsProjectOwner(gomock.Any(), uid, projID).Return(false, nil)
		m.EXPECT().GetPermission(gomock.Any(), uid, projID).Return(&project.Permission{
			Role: project.ProjectRoleViewer,
		}, nil)
	})

	req := httptest.NewRequest(http.MethodGet, "/projects/"+projID+"/people", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "insufficient project permissions")
}

func TestRequireProjectRole_NoPermission(t *testing.T) {
	uid := ulid.Make()
	projID := "01PROJ000000000000000005"

	w, r := newProjectRouter(t, uid, "staff", project.ActionRead, func(m *mock_repo.MockPermissionLoader) {
		m.EXPECT().IsProjectOwner(gomock.Any(), uid, projID).Return(false, nil)
		m.EXPECT().GetPermission(gomock.Any(), uid, projID).Return(nil, project.ErrPermissionNotFound)
	})

	req := httptest.NewRequest(http.MethodGet, "/projects/"+projID+"/people", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "no project access")
}

func TestRequireProjectRole_ProjectNotFound(t *testing.T) {
	uid := ulid.Make()
	projID := "01PROJ000000000000000006"

	w, r := newProjectRouter(t, uid, "staff", project.ActionRead, func(m *mock_repo.MockPermissionLoader) {
		m.EXPECT().IsProjectOwner(gomock.Any(), uid, projID).Return(false, project.ErrProjectNotFound)
	})

	req := httptest.NewRequest(http.MethodGet, "/projects/"+projID+"/people", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "project not found")
}

func TestRequireProjectRole_SensitivityFlags(t *testing.T) {
	uid := ulid.Make()
	projID := "01PROJ000000000000000007"

	w, r := newProjectRouter(t, uid, "staff", project.ActionRead, func(m *mock_repo.MockPermissionLoader) {
		m.EXPECT().IsProjectOwner(gomock.Any(), uid, projID).Return(false, nil)
		m.EXPECT().GetPermission(gomock.Any(), uid, projID).Return(&project.Permission{
			Role:             project.ProjectRoleManager,
			CanViewContact:   false,
			CanViewPersonal:  true,
			CanViewDocuments: false,
		}, nil)
	})

	req := httptest.NewRequest(http.MethodGet, "/projects/"+projID+"/people", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"can_view_contact":false`)
	assert.Contains(t, w.Body.String(), `"can_view_personal":true`)
	assert.Contains(t, w.Body.String(), `"can_view_documents":false`)
}
