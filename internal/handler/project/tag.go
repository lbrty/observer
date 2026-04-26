package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
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
		handler.InternalError(c, "list tags", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": out})
}

// Create handles POST /projects/:project_id/tags.
func (h *TagHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	input, ok := handler.BindJSON[ucproject.CreateTagInput](c)
	if !ok {
		return
	}
	out, err := h.uc.Create(c.Request.Context(), projectID, input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PUT /projects/:project_id/tags/:id.
func (h *TagHandler) Update(c *gin.Context) {
	input, ok := handler.BindJSON[ucproject.UpdateTagInput](c)
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

// Delete handles DELETE /projects/:project_id/tags/:id.
func (h *TagHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("project_id"), c.Param("id")); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tag deleted"})
}
