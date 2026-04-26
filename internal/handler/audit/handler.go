package audit

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
)

// AuditHandler exposes audit log HTTP endpoints.
type AuditHandler struct {
	auditUC *ucaudit.AuditUseCase
}

// NewAuditHandler creates an AuditHandler.
func NewAuditHandler(auditUC *ucaudit.AuditUseCase) *AuditHandler {
	return &AuditHandler{auditUC: auditUC}
}

// ListAll handles GET /admin/audit-logs (all projects).
func (h *AuditHandler) ListAll(c *gin.Context) {
	var input ucaudit.ListInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", err.Error()))
		return
	}

	out, err := h.auditUC.List(c.Request.Context(), input)
	if err != nil {
		handler.InternalError(c, "list audit logs", err)
		return
	}

	c.JSON(http.StatusOK, out)
}

// ListByProject handles GET /projects/:project_id/audit-logs (project-scoped).
func (h *AuditHandler) ListByProject(c *gin.Context) {
	projectID := c.Param("project_id")

	var input ucaudit.ListInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", err.Error()))
		return
	}

	input.ProjectID = &projectID

	out, err := h.auditUC.List(c.Request.Context(), input)
	if err != nil {
		handler.InternalError(c, "list project audit logs", err)
		return
	}

	c.JSON(http.StatusOK, out)
}
