package handler_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/auth"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/handler"
)

func newAuthHandler(d *authTestDeps) *handler.AuthHandler {
	return handler.NewAuthHandler(d.authUseCase(), d.userRepo, d.loginAttempts, testCookieConfig(), testJWTConfig())
}

// --- Register ---

func TestAuthHandler_Register_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	// Missing required fields
	c, w := newTestContext(http.MethodPost, "/auth/register", map[string]string{})
	h.Register(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "errors.validation", resp["code"])
}

func TestAuthHandler_Register_DuplicateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	d.userRepo.EXPECT().GetByEmail(gomock.Any(), "taken@test.com").Return(&user.User{}, nil)

	c, w := newTestContext(http.MethodPost, "/auth/register", map[string]string{
		"email":    "taken@test.com",
		"password": "password123",
		"role":     "staff",
	})
	h.Register(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "errors.user.emailExists", resp["code"])
}

func TestAuthHandler_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	d.userRepo.EXPECT().GetByEmail(gomock.Any(), "new@test.com").Return(nil, user.ErrUserNotFound)
	d.hasher.EXPECT().Hash("password123").Return("hashed", "salt", nil)
	d.userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	d.credRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContext(http.MethodPost, "/auth/register", map[string]string{
		"email":    "new@test.com",
		"password": "password123",
		"role":     "staff",
	})
	h.Register(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.NotEmpty(t, resp["user_id"])
}

// --- Login ---

func TestAuthHandler_Login_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	c, w := newTestContext(http.MethodPost, "/auth/login", map[string]string{})
	h.Login(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login_AccountLocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	d.loginAttempts.EXPECT().IsLocked(gomock.Any(), "locked@test.com").Return(5*time.Minute, nil)

	c, w := newTestContext(http.MethodPost, "/auth/login", map[string]string{
		"email":    "locked@test.com",
		"password": "password123",
	})
	h.Login(c)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "errors.auth.accountLocked", resp["code"])
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	d.loginAttempts.EXPECT().IsLocked(gomock.Any(), "user@test.com").Return(time.Duration(0), nil)
	d.userRepo.EXPECT().GetByEmail(gomock.Any(), "user@test.com").Return(nil, user.ErrUserNotFound)
	d.loginAttempts.EXPECT().RecordFailure(gomock.Any(), "user@test.com").Return(time.Duration(0), nil)

	c, w := newTestContext(http.MethodPost, "/auth/login", map[string]string{
		"email":    "user@test.com",
		"password": "wrong",
	})
	h.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_Login_InactiveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	inactiveUser := &user.User{
		ID:       testID(),
		Email:    "inactive@test.com",
		Role:     user.RoleStaff,
		IsActive: false,
	}

	d.loginAttempts.EXPECT().IsLocked(gomock.Any(), "inactive@test.com").Return(time.Duration(0), nil)
	d.userRepo.EXPECT().GetByEmail(gomock.Any(), "inactive@test.com").Return(inactiveUser, nil)

	c, w := newTestContext(http.MethodPost, "/auth/login", map[string]string{
		"email":    "inactive@test.com",
		"password": "password123",
	})
	h.Login(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	uid := testID()
	activeUser := &user.User{
		ID:       uid,
		Email:    "active@test.com",
		Role:     user.RoleStaff,
		IsActive: true,
	}
	cred := &user.Credentials{
		UserID:       uid,
		PasswordHash: "hash",
		Salt:         "salt",
	}
	expiresAt := time.Now().Add(15 * time.Minute)

	d.loginAttempts.EXPECT().IsLocked(gomock.Any(), "active@test.com").Return(time.Duration(0), nil)
	d.userRepo.EXPECT().GetByEmail(gomock.Any(), "active@test.com").Return(activeUser, nil)
	d.credRepo.EXPECT().GetByUserID(gomock.Any(), uid).Return(cred, nil)
	d.hasher.EXPECT().Verify("password123", "hash", "salt").Return(nil)
	d.mfaRepo.EXPECT().GetByUserID(gomock.Any(), uid).Return(nil, user.ErrUserNotFound)
	d.tokenGen.EXPECT().GenerateAccessToken(uid, "staff").Return("access-token", expiresAt, nil)
	d.tokenGen.EXPECT().GenerateRefreshToken().Return("refresh-token", nil)
	d.sessionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	d.loginAttempts.EXPECT().ClearAttempts(gomock.Any(), "active@test.com").Return(nil)

	c, w := newTestContext(http.MethodPost, "/auth/login", map[string]string{
		"email":    "active@test.com",
		"password": "password123",
	})
	h.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify cookies are set
	cookies := w.Result().Cookies()
	var foundAccess, foundRefresh bool
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			foundAccess = true
			assert.Equal(t, "access-token", cookie.Value)
		}
		if cookie.Name == "refresh_token" {
			foundRefresh = true
			assert.Equal(t, "refresh-token", cookie.Value)
		}
	}
	assert.True(t, foundAccess, "access_token cookie should be set")
	assert.True(t, foundRefresh, "refresh_token cookie should be set")
}

// --- Me ---

func TestAuthHandler_Me_NoAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	c, w := newTestContext(http.MethodGet, "/auth/me", nil)
	h.Me(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_Me_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	uid := testID()
	now := time.Now().UTC()
	u := &user.User{
		ID:         uid,
		FirstName:  "Test",
		LastName:   "User",
		Email:      "test@test.com",
		Phone:      "+1234567890",
		Role:       user.RoleStaff,
		IsVerified: true,
		CreatedAt:  now,
	}

	d.userRepo.EXPECT().GetByID(gomock.Any(), uid).Return(u, nil)

	c, w := newTestContext(http.MethodGet, "/auth/me", nil)
	setAuthContext(c, uid)
	h.Me(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, uid.String(), resp["id"])
	assert.Equal(t, "test@test.com", resp["email"])
}

// --- RefreshToken ---

func TestAuthHandler_RefreshToken_NoToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	c, w := newTestContext(http.MethodPost, "/auth/refresh", nil)
	h.RefreshToken(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_RefreshToken_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	d.sessionRepo.EXPECT().GetByRefreshToken(gomock.Any(), "bad-token").Return(nil, auth.ErrSessionNotFound)

	c, w := newTestContext(http.MethodPost, "/auth/refresh", map[string]string{
		"refresh_token": "bad-token",
	})
	h.RefreshToken(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	uid := testID()
	sessionID := testID()
	session := &auth.Session{
		ID:           sessionID,
		UserID:       uid,
		RefreshToken: "old-refresh",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
	activeUser := &user.User{
		ID:   uid,
		Role: user.RoleStaff,
	}
	expiresAt := time.Now().Add(15 * time.Minute)

	d.sessionRepo.EXPECT().GetByRefreshToken(gomock.Any(), "old-refresh").Return(session, nil)
	d.sessionRepo.EXPECT().Delete(gomock.Any(), sessionID).Return(nil)
	d.userRepo.EXPECT().GetByID(gomock.Any(), uid).Return(activeUser, nil)
	d.tokenGen.EXPECT().GenerateAccessToken(uid, "staff").Return("new-access", expiresAt, nil)
	d.tokenGen.EXPECT().GenerateRefreshToken().Return("new-refresh", nil)
	d.sessionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContext(http.MethodPost, "/auth/refresh", map[string]string{
		"refresh_token": "old-refresh",
	})
	h.RefreshToken(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Equal(t, "new-access", resp["access_token"])
	assert.Equal(t, "new-refresh", resp["refresh_token"])
}

// --- Logout ---

func TestAuthHandler_Logout_NoToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	c, w := newTestContext(http.MethodPost, "/auth/logout", nil)
	h.Logout(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	d.sessionRepo.EXPECT().DeleteByRefreshToken(gomock.Any(), "valid-refresh").Return(nil)

	c, w := newTestContext(http.MethodPost, "/auth/logout", map[string]string{
		"refresh_token": "valid-refresh",
	})
	h.Logout(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- UpdateProfile ---

func TestAuthHandler_UpdateProfile_NoAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	c, w := newTestContext(http.MethodPatch, "/auth/me", map[string]string{
		"first_name": "Updated",
	})
	h.UpdateProfile(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_UpdateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	uid := testID()
	now := time.Now().UTC()
	existingUser := &user.User{
		ID:        uid,
		FirstName: "Old",
		LastName:  "Name",
		Email:     "test@test.com",
		Role:      user.RoleStaff,
		CreatedAt: now,
	}

	d.userRepo.EXPECT().GetByID(gomock.Any(), uid).Return(existingUser, nil)
	d.userRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	firstName := "Updated"
	c, w := newTestContext(http.MethodPatch, "/auth/me", map[string]any{
		"first_name": firstName,
	})
	setAuthContext(c, uid)
	h.UpdateProfile(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- ChangePassword ---

func TestAuthHandler_ChangePassword_NoAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	c, w := newTestContext(http.MethodPost, "/auth/change-password", map[string]string{
		"current_password": "old",
		"new_password":     "newpassword123",
	})
	h.ChangePassword(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_ChangePassword_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	uid := testID()
	c, w := newTestContext(http.MethodPost, "/auth/change-password", map[string]string{})
	setAuthContext(c, uid)
	h.ChangePassword(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_ChangePassword_WrongCurrent(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	uid := testID()
	cred := &user.Credentials{
		UserID:       uid,
		PasswordHash: "hash",
		Salt:         "salt",
	}

	d.credRepo.EXPECT().GetByUserID(gomock.Any(), uid).Return(cred, nil)
	d.hasher.EXPECT().Verify("wrong", "hash", "salt").Return(user.ErrInvalidCredentials)

	c, w := newTestContext(http.MethodPost, "/auth/change-password", map[string]string{
		"current_password": "wrong",
		"new_password":     "newpassword123",
	})
	setAuthContext(c, uid)
	h.ChangePassword(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_ChangePassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := newAuthTestDeps(ctrl)
	h := newAuthHandler(d)

	uid := testID()
	cred := &user.Credentials{
		UserID:       uid,
		PasswordHash: "hash",
		Salt:         "salt",
	}

	d.credRepo.EXPECT().GetByUserID(gomock.Any(), uid).Return(cred, nil)
	d.hasher.EXPECT().Verify("oldpassword", "hash", "salt").Return(nil)
	d.hasher.EXPECT().Hash("newpassword123").Return("newhash", "newsalt", nil)
	d.credRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	c, w := newTestContext(http.MethodPost, "/auth/change-password", map[string]string{
		"current_password": "oldpassword",
		"new_password":     "newpassword123",
	})
	setAuthContext(c, uid)
	h.ChangePassword(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	require.Equal(t, "password changed successfully", resp["message"])
}
