package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/crypto"
	cryptomock "github.com/lbrty/observer/internal/crypto/mock"
	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/middleware"
	repomock "github.com/lbrty/observer/internal/repository/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupAuthContext(c *gin.Context, userID ulid.ULID, role string) {
	c.Set(string(middleware.CtxUserID), userID)
	c.Set(string(middleware.CtxUserRole), role)
}

func TestRequireRole_Allow(t *testing.T) {
	mw := middleware.NewAuthMiddleware(nil, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.GET("/test", func(c *gin.Context) {
		setupAuthContext(c, ulid.Make(), "admin")
		c.Next()
	}, mw.RequireRole(user.RoleAdmin, user.RoleStaff), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_Deny(t *testing.T) {
	mw := middleware.NewAuthMiddleware(nil, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.GET("/test", func(c *gin.Context) {
		setupAuthContext(c, ulid.Make(), "guest")
		c.Next()
	}, mw.RequireRole(user.RoleAdmin), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequireRole_NoRole(t *testing.T) {
	mw := middleware.NewAuthMiddleware(nil, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.GET("/test", mw.RequireRole(user.RoleAdmin), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAuthenticateUsesCurrentDatabaseRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	tokenGen := cryptomock.NewMockTokenGenerator(ctrl)
	userRepo := repomock.NewMockUserRepository(ctrl)
	userID := ulid.Make()

	tokenGen.EXPECT().ValidateAccessToken("signed-token").Return(&crypto.Claims{
		RegisteredClaims: jwt.RegisteredClaims{Subject: userID.String()},
		Role:             string(user.RoleAdmin),
		Type:             "access",
	}, nil)
	userRepo.EXPECT().GetByID(gomock.Any(), userID).Return(&user.User{
		ID:       userID,
		Role:     user.RoleGuest,
		IsActive: true,
	}, nil)

	router := gin.New()
	router.GET("/test", middleware.NewAuthMiddleware(tokenGen, userRepo).Authenticate(), func(c *gin.Context) {
		role, ok := middleware.UserRoleFrom(c)
		assert.True(t, ok)
		assert.Equal(t, user.RoleGuest, role)
		c.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", "Bearer signed-token")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNoContent, response.Code)
}
