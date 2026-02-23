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
	listUC   *ucadmin.ListPermissionsUseCase
	assignUC *ucadmin.AssignPermissionUseCase
	updateUC *ucadmin.UpdatePermissionUseCase
	revokeUC *ucadmin.RevokePermissionUseCase
}

// NewPermissionHandler creates a PermissionHandler.
func NewPermissionHandler(
	listUC *ucadmin.ListPermissionsUseCase,
	assignUC *ucadmin.AssignPermissionUseCase,
	updateUC *ucadmin.UpdatePermissionUseCase,
	revokeUC *ucadmin.RevokePermissionUseCase,
) *PermissionHandler {
	return &PermissionHandler{
		listUC:   listUC,
		assignUC: assignUC,
		updateUC: updateUC,
		revokeUC: revokeUC,
	}
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

	out, err := h.listUC.Execute(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.assignUC.Execute(c.Request.Context(), projectID, input)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.updateUC.Execute(c.Request.Context(), id, input)
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

	if err := h.revokeUC.Execute(c.Request.Context(), id); err != nil {
		h.handlePermError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission revoked"})
}

func (h *PermissionHandler) handlePermError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, project.ErrPermissionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, project.ErrPermissionExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, project.ErrInvalidProjectRole):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
