package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/middleware"
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
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	projectID := c.Param("project_id")
	userID, _ := middleware.UserIDFrom(c)
	role, _ := middleware.UserRoleFrom(c)

	out, err := h.permUC.List(c.Request.Context(), projectID, userID.String(), role)
	if err != nil {
		handler.InternalError(c, "list permissions", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": out})
}

// AssignPermission handles POST /admin/projects/:project_id/permissions.
func (h *PermissionHandler) AssignPermission(c *gin.Context) {
	projectID := c.Param("project_id")

	input, ok := handler.BindJSON[ucadmin.AssignPermissionInput](c)
	if !ok {
		return
	}

	out, err := h.permUC.Assign(c.Request.Context(), projectID, input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, out)
}

// UpdatePermission handles PATCH /admin/projects/:project_id/permissions/:id.
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id := c.Param("id")

	input, ok := handler.BindJSON[ucadmin.UpdatePermissionInput](c)
	if !ok {
		return
	}

	out, err := h.permUC.Update(c.Request.Context(), c.Param("project_id"), id, input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, out)
}

// RevokePermission handles DELETE /admin/projects/:project_id/permissions/:id.
func (h *PermissionHandler) RevokePermission(c *gin.Context) {
	projectID := c.Param("project_id")
	id := c.Param("id")

	if err := h.permUC.Revoke(c.Request.Context(), projectID, id); err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission revoked"})
}
