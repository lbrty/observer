package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/config"
	cryptomock "github.com/lbrty/observer/internal/crypto/mock"
	"github.com/lbrty/observer/internal/middleware"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
	ucauth "github.com/lbrty/observer/internal/usecase/auth"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestContext(method, path string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	var req *http.Request
	if body != nil {
		b, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

func newTestContextWithParams(method, path string, body any, params gin.Params) (*gin.Context, *httptest.ResponseRecorder) {
	c, w := newTestContext(method, path, body)
	c.Params = params
	return c, w
}

func setAuthContext(c *gin.Context, userID ulid.ULID) {
	c.Set(string(middleware.CtxUserID), userID)
}

func testID() ulid.ULID {
	return ulid.Make()
}

func parseResponse[T any](w *httptest.ResponseRecorder) T {
	var result T
	_ = json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

func testCookieConfig() config.CookieConfig {
	return config.CookieConfig{
		Domain:   "localhost",
		Secure:   false,
		SameSite: "lax",
		MaxAge:   2 * time.Hour,
	}
}

func testJWTConfig() config.JWTConfig {
	return config.JWTConfig{
		AccessTTL:  15 * time.Minute,
		RefreshTTL: 168 * time.Hour,
		Issuer:     "observer-test",
	}
}

// authTestDeps holds all mocks needed to construct an AuthHandler.
type authTestDeps struct {
	ctrl          *gomock.Controller
	userRepo      *repomock.MockUserRepository
	credRepo      *repomock.MockCredentialsRepository
	sessionRepo   *repomock.MockSessionRepository
	mfaRepo       *repomock.MockMFARepository
	hasher        *cryptomock.MockPasswordHasher
	tokenGen      *cryptomock.MockTokenGenerator
	loginAttempts *repomock.MockLoginAttemptStore
}

func newAuthTestDeps(ctrl *gomock.Controller) *authTestDeps {
	return &authTestDeps{
		ctrl:          ctrl,
		userRepo:      repomock.NewMockUserRepository(ctrl),
		credRepo:      repomock.NewMockCredentialsRepository(ctrl),
		sessionRepo:   repomock.NewMockSessionRepository(ctrl),
		mfaRepo:       repomock.NewMockMFARepository(ctrl),
		hasher:        cryptomock.NewMockPasswordHasher(ctrl),
		tokenGen:      cryptomock.NewMockTokenGenerator(ctrl),
		loginAttempts: repomock.NewMockLoginAttemptStore(ctrl),
	}
}

func (d *authTestDeps) authUseCase() *ucauth.AuthUseCase {
	return ucauth.NewAuthUseCase(d.userRepo, d.credRepo, d.sessionRepo, d.mfaRepo, d.hasher, d.tokenGen)
}

// adminTestDeps holds all mocks needed to construct an AdminHandler.
type adminTestDeps struct {
	ctrl          *gomock.Controller
	userRepo      *repomock.MockUserRepository
	credRepo      *repomock.MockCredentialsRepository
	hasher        *cryptomock.MockPasswordHasher
	loginAttempts *repomock.MockLoginAttemptStore
}

func newAdminTestDeps(ctrl *gomock.Controller) *adminTestDeps {
	return &adminTestDeps{
		ctrl:          ctrl,
		userRepo:      repomock.NewMockUserRepository(ctrl),
		credRepo:      repomock.NewMockCredentialsRepository(ctrl),
		hasher:        cryptomock.NewMockPasswordHasher(ctrl),
		loginAttempts: repomock.NewMockLoginAttemptStore(ctrl),
	}
}

func (d *adminTestDeps) userUseCase() *ucadmin.UserUseCase {
	return ucadmin.NewUserUseCase(d.userRepo, d.credRepo, d.hasher)
}
