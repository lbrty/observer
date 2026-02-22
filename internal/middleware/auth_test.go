package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"

	"github.com/lbrty/observer/internal/domain/user"
	"github.com/lbrty/observer/internal/middleware"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupAuthContext(c *gin.Context, userID ulid.ULID, role string) {
	c.Set(string(middleware.CtxUserID), userID)
	c.Set(string(middleware.CtxUserRole), role)
}

func TestRequireRole_Allow(t *testing.T) {
	mw := middleware.NewAuthMiddleware(nil)

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
	mw := middleware.NewAuthMiddleware(nil)

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
	mw := middleware.NewAuthMiddleware(nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.GET("/test", mw.RequireRole(user.RoleAdmin), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
