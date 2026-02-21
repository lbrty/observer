package health_test

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

func TestHealthHandler_Healthy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_database.NewMockDB(ctrl)
	mockDB.EXPECT().Ping(gomock.Any()).Return(nil)

	h := health.NewHandler(mockDB)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

	h.Health(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"ok"`)
}

func TestHealthHandler_Unhealthy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_database.NewMockDB(ctrl)
	mockDB.EXPECT().Ping(gomock.Any()).Return(assert.AnError)

	h := health.NewHandler(mockDB)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

	h.Health(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"not ok"`)
}
