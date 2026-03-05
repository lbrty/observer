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
// @Summary List projects
// @Tags admin-projects
// @Accept json
// @Produce json
// @Param owner_id query string false "Filter by owner ID"
// @Param status query string false "Filter by project status"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} ucadmin.ListProjectsOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/projects [get]
func (h *ProjectHandler) List(c *gin.Context) {
	var input ucadmin.ListProjectsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.List(c.Request.Context(), input)
	if err != nil {
		internalError(c, "list projects", err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Get handles GET /admin/projects/:project_id.
// @Summary Get project by ID
// @Tags admin-projects
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {object} ucadmin.ProjectDTO
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/projects/{project_id} [get]
func (h *ProjectHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("project_id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /admin/projects.
// @Summary Create a project
// @Tags admin-projects
// @Accept json
// @Produce json
// @Param input body ucadmin.CreateProjectInput true "Project creation payload"
// @Success 201 {object} ucadmin.ProjectDTO
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	var input ucadmin.CreateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
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

// Update handles PATCH /admin/projects/:project_id.
// @Summary Update a project
// @Tags admin-projects
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param input body ucadmin.UpdateProjectInput true "Project update payload"
// @Success 200 {object} ucadmin.ProjectDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/projects/{project_id} [patch]
func (h *ProjectHandler) Update(c *gin.Context) {
	var input ucadmin.UpdateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.Update(c.Request.Context(), c.Param("project_id"), input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *ProjectHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, project.ErrProjectNotFound):
		c.JSON(http.StatusNotFound, errJSON("errors.project.notFound", err.Error()))
	case errors.Is(err, project.ErrProjectNameExists):
		c.JSON(http.StatusConflict, errJSON("errors.project.nameExists", err.Error()))
	default:
		internalError(c, "handle project", err)
	}
}
