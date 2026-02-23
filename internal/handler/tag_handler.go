package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/tag"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

// TagHandler exposes tag CRUD HTTP endpoints.
type TagHandler struct {
	uc *ucproject.TagUseCase
}

// NewTagHandler creates a TagHandler.
func NewTagHandler(uc *ucproject.TagUseCase) *TagHandler {
	return &TagHandler{uc: uc}
}

// List handles GET /projects/:project_id/tags.
func (h *TagHandler) List(c *gin.Context) {
	projectID := c.Param("project_id")
	out, err := h.uc.List(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": out})
}

// Create handles POST /projects/:project_id/tags.
func (h *TagHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.CreateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.Create(c.Request.Context(), projectID, input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Delete handles DELETE /projects/:project_id/tags/:id.
func (h *TagHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tag deleted"})
}

func (h *TagHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, tag.ErrTagNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, tag.ErrTagNameExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
