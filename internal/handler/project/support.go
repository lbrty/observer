package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/middleware"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

// SupportRecordHandler exposes support record HTTP endpoints.
type SupportRecordHandler struct {
	uc *ucproject.SupportRecordUseCase
}

// NewSupportRecordHandler creates a SupportRecordHandler.
func NewSupportRecordHandler(uc *ucproject.SupportRecordUseCase) *SupportRecordHandler {
	return &SupportRecordHandler{uc: uc}
}

// List handles GET /projects/:project_id/support-records.
func (h *SupportRecordHandler) List(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.ListSupportRecordsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.List(c.Request.Context(), projectID, input)
	if err != nil {
		handler.InternalError(c, "list support records", err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Get handles GET /projects/:project_id/support-records/:id.
func (h *SupportRecordHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("project_id"), c.Param("id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/support-records.
func (h *SupportRecordHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	input, ok := handler.BindJSON[ucproject.CreateSupportRecordInput](c)
	if !ok {
		return
	}
	userID, _ := middleware.UserIDFrom(c)
	out, err := h.uc.Create(c.Request.Context(), projectID, userID.String(), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /projects/:project_id/support-records/:id.
func (h *SupportRecordHandler) Update(c *gin.Context) {
	input, ok := handler.BindJSON[ucproject.UpdateSupportRecordInput](c)
	if !ok {
		return
	}
	out, err := h.uc.Update(c.Request.Context(), c.Param("project_id"), c.Param("id"), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /projects/:project_id/support-records/:id.
func (h *SupportRecordHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("project_id"), c.Param("id")); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "support record deleted"})
}
