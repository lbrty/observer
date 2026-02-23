package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/support"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.List(c.Request.Context(), projectID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, out)
}

// Get handles GET /projects/:project_id/support-records/:id.
func (h *SupportRecordHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/support-records.
func (h *SupportRecordHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.CreateSupportRecordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := middleware.UserIDFrom(c)
	out, err := h.uc.Create(c.Request.Context(), projectID, userID.String(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /projects/:project_id/support-records/:id.
func (h *SupportRecordHandler) Update(c *gin.Context) {
	var input ucproject.UpdateSupportRecordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /projects/:project_id/support-records/:id.
func (h *SupportRecordHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "support record deleted"})
}

func (h *SupportRecordHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, support.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
