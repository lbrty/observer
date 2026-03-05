package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/project"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

// PermissionHandler exposes project permission management HTTP endpoints.
type PermissionHandler struct {
	permUC *ucadmin.PermissionUseCase
}

// NewPermissionHandler creates a PermissionHandler.
func NewPermissionHandler(permUC *ucadmin.PermissionUseCase) *PermissionHandler {
	return &PermissionHandler{permUC: permUC}
}

// ListPermissions handles GET /admin/projects/:project_id/permissions.
// @Summary List project permissions
// @Tags admin-permissions
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {object} PermissionListResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/projects/{project_id}/permissions [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	projectID := c.Param("project_id")

	out, err := h.permUC.List(c.Request.Context(), projectID)
	if err != nil {
		internalError(c, "list permissions", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": out})
}

// AssignPermission handles POST /admin/projects/:project_id/permissions.
// @Summary Assign a project permission
// @Tags admin-permissions
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param input body ucadmin.AssignPermissionInput true "Permission assignment payload"
// @Success 201 {object} ucadmin.PermissionDTO
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/projects/{project_id}/permissions [post]
func (h *PermissionHandler) AssignPermission(c *gin.Context) {
	projectID := c.Param("project_id")

	var input ucadmin.AssignPermissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}

	out, err := h.permUC.Assign(c.Request.Context(), projectID, input)
	if err != nil {
		h.handlePermError(c, err)
		return
	}

	c.JSON(http.StatusCreated, out)
}

// UpdatePermission handles PATCH /admin/projects/:project_id/permissions/:id.
// @Summary Update a project permission
// @Tags admin-permissions
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param id path string true "Permission ID"
// @Param input body ucadmin.UpdatePermissionInput true "Permission update payload"
// @Success 200 {object} ucadmin.PermissionDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/projects/{project_id}/permissions/{id} [patch]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id := c.Param("id")

	var input ucadmin.UpdatePermissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}

	out, err := h.permUC.Update(c.Request.Context(), id, input)
	if err != nil {
		h.handlePermError(c, err)
		return
	}

	c.JSON(http.StatusOK, out)
}

// RevokePermission handles DELETE /admin/projects/:project_id/permissions/:id.
// @Summary Revoke a project permission
// @Tags admin-permissions
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param id path string true "Permission ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /admin/projects/{project_id}/permissions/{id} [delete]
func (h *PermissionHandler) RevokePermission(c *gin.Context) {
	id := c.Param("id")

	if err := h.permUC.Revoke(c.Request.Context(), id); err != nil {
		h.handlePermError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission revoked"})
}

func (h *PermissionHandler) handlePermError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, project.ErrPermissionNotFound):
		c.JSON(http.StatusNotFound, errJSON("errors.project.permissionNotFound", err.Error()))
	case errors.Is(err, project.ErrPermissionExists):
		c.JSON(http.StatusConflict, errJSON("errors.project.permissionExists", err.Error()))
	case errors.Is(err, project.ErrInvalidProjectRole):
		c.JSON(http.StatusBadRequest, errJSON("errors.project.invalidRole", err.Error()))
	default:
		internalError(c, "handle permission", err)
	}
}
