package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ucreport "github.com/lbrty/observer/internal/usecase/report"
)

// ReportHandler exposes report HTTP endpoints.
type ReportHandler struct {
	uc *ucreport.ReportUseCase
}

// NewReportHandler creates a ReportHandler.
func NewReportHandler(uc *ucreport.ReportUseCase) *ReportHandler {
	return &ReportHandler{uc: uc}
}

// Generate handles GET /projects/:project_id/reports.
func (h *ReportHandler) Generate(c *gin.Context) {
	projectID := c.Param("project_id")

	var input ucreport.ReportInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}

	out, err := h.uc.Generate(c.Request.Context(), projectID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errJSON("errors.internal", "internal server error"))
		return
	}

	c.JSON(http.StatusOK, out)
}
