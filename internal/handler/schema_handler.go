package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SchemaStatus holds the migration drift result computed at startup.
type SchemaStatus struct {
	CurrentVersion uint `json:"current_version"`
	LatestVersion  uint `json:"latest_version"`
	Pending        int  `json:"pending"`
	Dirty          bool `json:"dirty"`
}

type schemaHandler struct {
	status SchemaStatus
}

// NewSchemaHandler creates a handler that returns the pre-computed schema status.
func NewSchemaHandler(status SchemaStatus) *schemaHandler {
	return &schemaHandler{status: status}
}

// Status returns the current migration drift status.
func (h *schemaHandler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, h.status)
}
