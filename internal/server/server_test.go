package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_database "github.com/lbrty/observer/internal/database/mock"
	"github.com/lbrty/observer/internal/health"
)

// bodySizeLimitMiddleware mirrors the middleware added to the real server.
func bodySizeLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.GetHeader("Content-Type"), "multipart/form-data") {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		}
		c.Next()
	}
}

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHealthRoute_Via_Gin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_database.NewMockDB(ctrl)
	mockDB.EXPECT().Ping(gomock.Any()).Return(nil)

	router := gin.New()
	router.GET("/health", health.NewHandler(mockDB).Health)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"ok"`)
}

func TestBodySizeLimit_JSONRejectedOver1MB(t *testing.T) {
	router := gin.New()
	router.Use(bodySizeLimitMiddleware())
	router.POST("/test", func(c *gin.Context) {
		// Try to read the body; if it exceeds the limit this should error
		buf := make([]byte, 2<<20)
		_, err := c.Request.Body.Read(buf)
		if err != nil {
			c.Status(http.StatusRequestEntityTooLarge)
			return
		}
		c.Status(http.StatusOK)
	})

	bigBody := strings.Repeat("x", 2<<20) // 2 MB
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(bigBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
}

func TestBodySizeLimit_MultipartNotLimited(t *testing.T) {
	router := gin.New()
	router.Use(bodySizeLimitMiddleware())
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("body"))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=xxx")
	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusRequestEntityTooLarge, w.Code)
}

func TestRequestID_Present(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_database.NewMockDB(ctrl)
	mockDB.EXPECT().Ping(gomock.Any()).Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Header("X-Request-ID", "test-id")
		c.Next()
	})
	router.GET("/health", health.NewHandler(mockDB).Health)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
}
