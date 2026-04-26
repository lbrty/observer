package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
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
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", err.Error()))
		return
	}

	userID, _ := middleware.UserIDFrom(c)
	role, _ := middleware.UserRoleFrom(c)
	input.CallerID = userID.String()
	input.CallerRole = role

	out, err := h.uc.List(c.Request.Context(), input)
	if err != nil {
		handler.InternalError(c, "list projects", err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Get handles GET /admin/projects/:project_id.
func (h *ProjectHandler) Get(c *gin.Context) {
	userID, _ := middleware.UserIDFrom(c)
	role, _ := middleware.UserRoleFrom(c)

	out, err := h.uc.Get(c.Request.Context(), c.Param("project_id"), userID.String(), role)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /admin/projects.
func (h *ProjectHandler) Create(c *gin.Context) {
	input, ok := handler.BindJSON[ucadmin.CreateProjectInput](c)
	if !ok {
		return
	}
	userID, _ := middleware.UserIDFrom(c)
	out, err := h.uc.Create(c.Request.Context(), userID.String(), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /admin/projects/:project_id.
func (h *ProjectHandler) Update(c *gin.Context) {
	input, ok := handler.BindJSON[ucadmin.UpdateProjectInput](c)
	if !ok {
		return
	}
	out, err := h.uc.Update(c.Request.Context(), c.Param("project_id"), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}
