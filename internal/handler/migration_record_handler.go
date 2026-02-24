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
// @Summary List migration records for a person
// @Tags project-migration-records
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param person_id path string true "Person ID"
// @Success 200 {object} object "Wrapper with records array"
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people/{person_id}/migration-records [get]
func (h *MigrationRecordHandler) List(c *gin.Context) {
	personID := c.Param("person_id")
	out, err := h.uc.ListByPerson(c.Request.Context(), personID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errJSON("errors.internal", "internal server error"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"records": out})
}

// Get handles GET /projects/:project_id/people/:person_id/migration-records/:id.
// @Summary Get a migration record by ID
// @Tags project-migration-records
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param person_id path string true "Person ID"
// @Param id path string true "Migration record ID"
// @Success 200 {object} ucproject.MigrationRecordDTO
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people/{person_id}/migration-records/{id} [get]
func (h *MigrationRecordHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/people/:person_id/migration-records.
// @Summary Create a migration record for a person
// @Tags project-migration-records
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param person_id path string true "Person ID"
// @Param input body ucproject.CreateMigrationRecordInput true "Migration record payload"
// @Success 201 {object} ucproject.MigrationRecordDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people/{person_id}/migration-records [post]
func (h *MigrationRecordHandler) Create(c *gin.Context) {
	personID := c.Param("person_id")
	var input ucproject.CreateMigrationRecordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
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
		c.JSON(http.StatusNotFound, errJSON("errors.migration.notFound", err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, errJSON("errors.internal", "internal server error"))
	}
}
