package handler_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/middleware"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

type permissionTestDeps struct {
	ctrl     *gomock.Controller
	permRepo *repomock.MockPermissionRepository
	userRepo *repomock.MockUserRepository
}

func newPermissionTestDeps(ctrl *gomock.Controller) *permissionTestDeps {
	return &permissionTestDeps{
		ctrl:     ctrl,
		permRepo: repomock.NewMockPermissionRepository(ctrl),
		userRepo: repomock.NewMockUserRepository(ctrl),
	}
}

func (d *permissionTestDeps) permissionUseCase() *ucadmin.PermissionUseCase {
	return ucadmin.NewPermissionUseCase(d.permRepo, d.userRepo)
}

func newPermissionHandler(d *permissionTestDeps) *handler.PermissionHandler {
	return handler.NewPermissionHandler(d.permissionUseCase())
}

func setPermAdminAuth(c *gin.Context) ulid.ULID {
	uid := testID()
	setAuthContext(c, uid)
	c.Set(string(middleware.CtxUserRole), string(user.RoleAdmin))
	return uid
}

// --- ListPermissions ---

func TestPermissionHandler_ListPermissions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newPermissionTestDeps(ctrl)
	h := newPermissionHandler(d)

	now := time.Now().UTC()
	projectID := testID().String()
	userUID := ulid.Make()
	permID := testID().String()

	d.permRepo.EXPECT().List(gomock.Any(), projectID).Return([]*project.ProjectPermission{
		{
			ID:               permID,
			ProjectID:        projectID,
			UserID:           userUID.String(),
			Role:             project.ProjectRoleConsultant,
			CanViewContact:   true,
			CanViewPersonal:  false,
			CanViewDocuments: false,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
	}, nil)

	d.userRepo.EXPECT().GetByIDs(gomock.Any(), gomock.Any()).Return([]*user.User{
		{
			ID:        userUID,
			FirstName: "Alice",
			LastName:  "Test",
			Email:     "alice@test.com",
			Role:      user.RoleConsultant,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/admin/projects/"+projectID+"/permissions", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	setPermAdminAuth(c)
	h.ListPermissions(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	perms := resp["permissions"].([]any)
	assert.Len(t, perms, 1)
}

func TestPermissionHandler_ListPermissions_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newPermissionTestDeps(ctrl)
	h := newPermissionHandler(d)

	projectID := testID().String()
	d.permRepo.EXPECT().List(gomock.Any(), projectID).Return(nil, fmt.Errorf("db error"))

	c, w := newTestContextWithParams(http.MethodGet, "/admin/projects/"+projectID+"/permissions", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	setPermAdminAuth(c)
	h.ListPermissions(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- AssignPermission ---

func TestPermissionHandler_AssignPermission_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newPermissionTestDeps(ctrl)
	h := newPermissionHandler(d)

	projectID := testID().String()
	c, w := newTestContextWithParams(http.MethodPost, "/admin/projects/"+projectID+"/permissions",
		map[string]any{}, gin.Params{{Key: "project_id", Value: projectID}})
	h.AssignPermission(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPermissionHandler_AssignPermission_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newPermissionTestDeps(ctrl)
	h := newPermissionHandler(d)

	projectID := testID().String()
	c, w := newTestContextWithParams(http.MethodPost, "/admin/projects/"+projectID+"/permissions", map[string]any{
		"user_id": testID().String(),
		"role":    "invalid",
	}, gin.Params{{Key: "project_id", Value: projectID}})
	h.AssignPermission(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPermissionHandler_AssignPermission_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newPermissionTestDeps(ctrl)
	h := newPermissionHandler(d)

	projectID := testID().String()
	d.permRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(project.ErrPermissionExists)

	c, w := newTestContextWithParams(http.MethodPost, "/admin/projects/"+projectID+"/permissions", map[string]any{
		"user_id": testID().String(),
		"role":    "consultant",
	}, gin.Params{{Key: "project_id", Value: projectID}})
	h.AssignPermission(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestPermissionHandler_AssignPermission_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newPermissionTestDeps(ctrl)
	h := newPermissionHandler(d)

	projectID := testID().String()
	d.permRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPost, "/admin/projects/"+projectID+"/permissions", map[string]any{
		"user_id":            testID().String(),
		"role":               "consultant",
		"can_view_contact":   true,
		"can_view_personal":  false,
		"can_view_documents": false,
	}, gin.Params{{Key: "project_id", Value: projectID}})
	h.AssignPermission(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, "consultant", resp["role"])
}

// --- UpdatePermission ---

func TestPermissionHandler_UpdatePermission_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newPermissionTestDeps(ctrl)
	h := newPermissionHandler(d)

	permID := testID().String()
	d.permRepo.EXPECT().GetByID(gomock.Any(), permID).Return(nil, project.ErrPermissionNotFound)

	c, w := newTestContextWithParams(http.MethodPatch, "/admin/projects/x/permissions/"+permID, map[string]any{
		"role": "viewer",
	}, gin.Params{
		{Key: "project_id", Value: testID().String()},
		{Key: "id", Value: permID},
	})
	h.UpdatePermission(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPermissionHandler_UpdatePermission_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newPermissionTestDeps(ctrl)
	h := newPermissionHandler(d)

	now := time.Now().UTC()
	permID := testID().String()
	projectID := testID().String()
	existing := &project.ProjectPermission{
		ID:               permID,
		ProjectID:        projectID,
		UserID:           testID().String(),
		Role:             project.ProjectRoleConsultant,
		CanViewContact:   false,
		CanViewPersonal:  false,
		CanViewDocuments: false,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	d.permRepo.EXPECT().GetByID(gomock.Any(), permID).Return(existing, nil)
	d.permRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPatch, "/admin/projects/"+projectID+"/permissions/"+permID, map[string]any{
		"role": "viewer",
	}, gin.Params{
		{Key: "project_id", Value: projectID},
		{Key: "id", Value: permID},
	})
	h.UpdatePermission(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, permID, resp["id"])
	assert.Equal(t, "viewer", resp["role"])
}

// --- RevokePermission ---

func TestPermissionHandler_RevokePermission_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newPermissionTestDeps(ctrl)
	h := newPermissionHandler(d)

	permID := testID().String()
	d.permRepo.EXPECT().Delete(gomock.Any(), permID).Return(project.ErrPermissionNotFound)

	c, w := newTestContextWithParams(http.MethodDelete, "/admin/projects/x/permissions/"+permID, nil, gin.Params{
		{Key: "project_id", Value: testID().String()},
		{Key: "id", Value: permID},
	})
	h.RevokePermission(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPermissionHandler_RevokePermission_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newPermissionTestDeps(ctrl)
	h := newPermissionHandler(d)

	permID := testID().String()
	d.permRepo.EXPECT().Delete(gomock.Any(), permID).Return(nil)

	c, w := newTestContextWithParams(http.MethodDelete, "/admin/projects/x/permissions/"+permID, nil, gin.Params{
		{Key: "project_id", Value: testID().String()},
		{Key: "id", Value: permID},
	})
	h.RevokePermission(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "permission revoked", resp["message"])
}
