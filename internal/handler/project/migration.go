package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
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
		handler.InternalError(c, "list migration records", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"records": out})
}

// Get handles GET /projects/:project_id/people/:person_id/migration-records/:id.
func (h *MigrationRecordHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("person_id"), c.Param("id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/people/:person_id/migration-records.
func (h *MigrationRecordHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	personID := c.Param("person_id")
	input, ok := handler.BindJSON[ucproject.CreateMigrationRecordInput](c)
	if !ok {
		return
	}
	out, err := h.uc.Create(c.Request.Context(), projectID, personID, input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /projects/:project_id/people/:person_id/migration-records/:id.
func (h *MigrationRecordHandler) Update(c *gin.Context) {
	input, ok := handler.BindJSON[ucproject.UpdateMigrationRecordInput](c)
	if !ok {
		return
	}
	out, err := h.uc.Update(c.Request.Context(), c.Param("project_id"), c.Param("person_id"), c.Param("id"), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /projects/:project_id/people/:person_id/migration-records/:id.
func (h *MigrationRecordHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("project_id"), c.Param("person_id"), c.Param("id")); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "migration record deleted"})
}
