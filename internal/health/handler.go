package health

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/database"
)

//go:generate mockgen -destination=mock/handler.go -package=mock github.com/lbrty/observer/internal/health Handler

// Handler defines the health check interface.
type Handler interface {
	Health(c *gin.Context)
}

type handler struct {
	db database.DB
}

// NewHandler creates a health Handler.
func NewHandler(db database.DB) Handler {
	return &handler{db: db}
}

// Health responds with the current health status.
// @Summary Health check
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *handler) Health(c *gin.Context) {
	if err := h.db.Ping(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ok"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
