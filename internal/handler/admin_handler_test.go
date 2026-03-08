package handler_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/handler"
)

func newAdminHandler(d *adminTestDeps) *handler.AdminHandler {
	return handler.NewAdminHandler(d.userUseCase(), d.loginAttempts)
}

// --- ListUsers ---

func TestAdminHandler_ListUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	uid := testID()
	now := time.Now().UTC()
	d.userRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*user.User{
		{ID: uid, FirstName: "Alice", Email: "alice@test.com", Role: user.RoleStaff, CreatedAt: now, UpdatedAt: now},
	}, 1, nil)

	c, w := newTestContext(http.MethodGet, "/admin/users?page=1&per_page=10", nil)
	h.ListUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	users := resp["users"].([]any)
	assert.Len(t, users, 1)
	assert.Equal(t, float64(1), resp["total"])
}

func TestAdminHandler_ListUsers_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	d.userRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, 0, fmt.Errorf("db error"))

	c, w := newTestContext(http.MethodGet, "/admin/users", nil)
	h.ListUsers(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- GetUser ---

func TestAdminHandler_GetUser_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	c, w := newTestContextWithParams(http.MethodGet, "/admin/users/bad-id", nil, gin.Params{
		{Key: "id", Value: "bad-id"},
	})
	h.GetUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_GetUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	uid := testID()
	d.userRepo.EXPECT().GetByID(gomock.Any(), uid).Return(nil, user.ErrUserNotFound)

	c, w := newTestContextWithParams(http.MethodGet, "/admin/users/"+uid.String(), nil, gin.Params{
		{Key: "id", Value: uid.String()},
	})
	h.GetUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAdminHandler_GetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	uid := testID()
	now := time.Now().UTC()
	d.userRepo.EXPECT().GetByID(gomock.Any(), uid).Return(&user.User{
		ID:        uid,
		FirstName: "Alice",
		Email:     "alice@test.com",
		Role:      user.RoleStaff,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	c, w := newTestContextWithParams(http.MethodGet, "/admin/users/"+uid.String(), nil, gin.Params{
		{Key: "id", Value: uid.String()},
	})
	h.GetUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, uid.String(), resp["id"])
}

// --- CreateUser ---

func TestAdminHandler_CreateUser_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	c, w := newTestContext(http.MethodPost, "/admin/users", map[string]string{})
	h.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_CreateUser_DuplicateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	d.userRepo.EXPECT().GetByEmail(gomock.Any(), "taken@test.com").Return(&user.User{}, nil)

	c, w := newTestContext(http.MethodPost, "/admin/users", map[string]any{
		"first_name": "Alice",
		"email":      "taken@test.com",
		"password":   "password123",
		"role":       "staff",
	})
	h.CreateUser(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestAdminHandler_CreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	d.userRepo.EXPECT().GetByEmail(gomock.Any(), "new@test.com").Return(nil, user.ErrUserNotFound)
	d.hasher.EXPECT().Hash("password123").Return("hash", "salt", nil)
	d.userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	d.credRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContext(http.MethodPost, "/admin/users", map[string]any{
		"first_name": "Alice",
		"email":      "new@test.com",
		"password":   "password123",
		"role":       "staff",
	})
	h.CreateUser(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["id"])
	assert.Equal(t, "new@test.com", resp["email"])
}

// --- UpdateUser ---

func TestAdminHandler_UpdateUser_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	c, w := newTestContextWithParams(http.MethodPatch, "/admin/users/bad", map[string]any{
		"first_name": "Updated",
	}, gin.Params{{Key: "id", Value: "bad"}})
	h.UpdateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_UpdateUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	uid := testID()
	d.userRepo.EXPECT().GetByID(gomock.Any(), uid).Return(nil, user.ErrUserNotFound)

	firstName := "Updated"
	c, w := newTestContextWithParams(http.MethodPatch, "/admin/users/"+uid.String(), map[string]any{
		"first_name": firstName,
	}, gin.Params{{Key: "id", Value: uid.String()}})
	h.UpdateUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAdminHandler_UpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	uid := testID()
	now := time.Now().UTC()
	existingUser := &user.User{
		ID:        uid,
		FirstName: "Old",
		Email:     "alice@test.com",
		Role:      user.RoleStaff,
		CreatedAt: now,
		UpdatedAt: now,
	}

	d.userRepo.EXPECT().GetByID(gomock.Any(), uid).Return(existingUser, nil)
	d.userRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	firstName := "Updated"
	c, w := newTestContextWithParams(http.MethodPatch, "/admin/users/"+uid.String(), map[string]any{
		"first_name": firstName,
	}, gin.Params{{Key: "id", Value: uid.String()}})
	h.UpdateUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, uid.String(), resp["id"])
}

// --- ResetPassword ---

func TestAdminHandler_ResetPassword_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	c, w := newTestContextWithParams(http.MethodPost, "/admin/users/bad/reset-password", map[string]any{
		"new_password": "newpassword123",
	}, gin.Params{{Key: "id", Value: "bad"}})
	h.ResetPassword(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_ResetPassword_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	uid := testID()
	c, w := newTestContextWithParams(http.MethodPost, "/admin/users/"+uid.String()+"/reset-password",
		map[string]any{}, gin.Params{{Key: "id", Value: uid.String()}})
	h.ResetPassword(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_ResetPassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	uid := testID()
	cred := &user.Credentials{
		UserID:       uid,
		PasswordHash: "oldhash",
		Salt:         "oldsalt",
	}

	d.credRepo.EXPECT().GetByUserID(gomock.Any(), uid).Return(cred, nil)
	d.hasher.EXPECT().Hash("newpassword123").Return("newhash", "newsalt", nil)
	d.credRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContextWithParams(http.MethodPost, "/admin/users/"+uid.String()+"/reset-password",
		map[string]any{"new_password": "newpassword123"}, gin.Params{{Key: "id", Value: uid.String()}})
	h.ResetPassword(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- UnlockAccount ---

func TestAdminHandler_UnlockAccount_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	c, w := newTestContextWithParams(http.MethodPost, "/admin/users/bad/unlock", nil,
		gin.Params{{Key: "id", Value: "bad"}})
	h.UnlockAccount(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_UnlockAccount_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	uid := testID()
	d.userRepo.EXPECT().GetByID(gomock.Any(), uid).Return(nil, user.ErrUserNotFound)

	c, w := newTestContextWithParams(http.MethodPost, "/admin/users/"+uid.String()+"/unlock", nil,
		gin.Params{{Key: "id", Value: uid.String()}})
	h.UnlockAccount(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAdminHandler_UnlockAccount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAdminTestDeps(ctrl)
	h := newAdminHandler(d)

	uid := testID()
	now := time.Now().UTC()
	u := &user.User{
		ID:        uid,
		Email:     "locked@test.com",
		Role:      user.RoleStaff,
		CreatedAt: now,
		UpdatedAt: now,
	}

	d.userRepo.EXPECT().GetByID(gomock.Any(), uid).Return(u, nil)
	d.loginAttempts.EXPECT().ClearAttempts(gomock.Any(), "locked@test.com").Return(nil)

	c, w := newTestContextWithParams(http.MethodPost, "/admin/users/"+uid.String()+"/unlock", nil,
		gin.Params{{Key: "id", Value: uid.String()}})
	h.UnlockAccount(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "account unlocked", resp["message"])
}
