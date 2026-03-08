package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
		internalError(c, "list migration records", err)
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
		HandleError(c, err)
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
	projectID := c.Param("project_id")
	personID := c.Param("person_id")
	var input ucproject.CreateMigrationRecordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.Create(c.Request.Context(), projectID, personID, input)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /projects/:project_id/people/:person_id/migration-records/:id.
func (h *MigrationRecordHandler) Update(c *gin.Context) {
	var input ucproject.UpdateMigrationRecordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

