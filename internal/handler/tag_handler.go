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
// @Summary List tags in a project
// @Tags project-tags
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Success 200 {object} object "Wrapper with tags array"
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/tags [get]
func (h *TagHandler) List(c *gin.Context) {
	projectID := c.Param("project_id")
	out, err := h.uc.List(c.Request.Context(), projectID)
	if err != nil {
		internalError(c, "list tags", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": out})
}

// Create handles POST /projects/:project_id/tags.
// @Summary Create a tag in a project
// @Tags project-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param input body ucproject.CreateTagInput true "Tag payload"
// @Success 201 {object} ucproject.TagDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/tags [post]
func (h *TagHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.CreateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
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
// @Summary Delete a tag
// @Tags project-tags
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param id path string true "Tag ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/tags/{id} [delete]
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
		c.JSON(http.StatusNotFound, errJSON("errors.tag.notFound", err.Error()))
	case errors.Is(err, tag.ErrTagNameExists):
		c.JSON(http.StatusConflict, errJSON("errors.tag.nameExists", err.Error()))
	default:
		internalError(c, "handle tag", err)
	}
}
