package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/migration"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

// MigrationRecordHandler exposes migration record HTTP endpoints.
type MigrationRecordHandler struct {
	uc *ucproject.MigrationRecordUseCase
}

// NewMigrationRecordHandler creates a MigrationRecordHandler.
func NewMigrationRecordHandler(uc *ucproject.MigrationRecordUseCase) *MigrationRecordHandler {
	return &MigrationRecordHandler{uc: uc}
}

// List handles GET /projects/:project_id/people/:person_id/migration-records.
func (h *MigrationRecordHandler) List(c *gin.Context) {
	personID := c.Param("person_id")
	out, err := h.uc.ListByPerson(c.Request.Context(), personID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"records": out})
}

// Get handles GET /projects/:project_id/people/:person_id/migration-records/:id.
func (h *MigrationRecordHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/people/:person_id/migration-records.
func (h *MigrationRecordHandler) Create(c *gin.Context) {
	personID := c.Param("person_id")
	var input ucproject.CreateMigrationRecordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.Create(c.Request.Context(), personID, input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *MigrationRecordHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, migration.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
