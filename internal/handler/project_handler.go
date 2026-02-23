package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/middleware"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

// ProjectHandler exposes project CRUD HTTP endpoints.
type ProjectHandler struct {
	uc *ucadmin.ProjectUseCase
}

// NewProjectHandler creates a ProjectHandler.
func NewProjectHandler(uc *ucadmin.ProjectUseCase) *ProjectHandler {
	return &ProjectHandler{uc: uc}
}

// List handles GET /admin/projects.
func (h *ProjectHandler) List(c *gin.Context) {
	var input ucadmin.ListProjectsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.List(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, out)
}

// Get handles GET /admin/projects/:id.
func (h *ProjectHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /admin/projects.
func (h *ProjectHandler) Create(c *gin.Context) {
	var input ucadmin.CreateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := middleware.UserIDFrom(c)
	out, err := h.uc.Create(c.Request.Context(), userID.String(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /admin/projects/:id.
func (h *ProjectHandler) Update(c *gin.Context) {
	var input ucadmin.UpdateProjectInput
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

func (h *ProjectHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, project.ErrProjectNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, project.ErrProjectNameExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
