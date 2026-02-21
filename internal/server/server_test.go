package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_database "github.com/lbrty/observer/internal/database/mock"
	"github.com/lbrty/observer/internal/health"
)

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
