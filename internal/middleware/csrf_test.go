package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/lbrty/observer/internal/middleware"
)

func setupCSRFRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.CSRFProtection())
	r.POST("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func TestCSRF_MissingCookieRejected(t *testing.T) {
	r := setupCSRFRouter()
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCSRF_MissingHeaderRejected(t *testing.T) {
	r := setupCSRFRouter()
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}"))
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: "abc123"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCSRF_MismatchedTokenRejected(t *testing.T) {
	r := setupCSRFRouter()
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}"))
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: "abc123"})
	req.Header.Set("X-CSRF-Token", "wrong")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCSRF_MatchingTokenAllowed(t *testing.T) {
	r := setupCSRFRouter()
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}"))
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: "abc123"})
	req.Header.Set("X-CSRF-Token", "abc123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCSRF_GETAllowedWithoutToken(t *testing.T) {
	r := setupCSRFRouter()
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
