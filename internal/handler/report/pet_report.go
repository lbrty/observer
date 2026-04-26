package report

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
	ucreport "github.com/lbrty/observer/internal/usecase/report"
)

// PetReportHandler exposes pet report HTTP endpoints.
type PetReportHandler struct {
	uc *ucreport.PetReportUseCase
}

// NewPetReportHandler creates a PetReportHandler.
func NewPetReportHandler(uc *ucreport.PetReportUseCase) *PetReportHandler {
	return &PetReportHandler{uc: uc}
}

// Generate handles GET /projects/:project_id/reports/pets.
func (h *PetReportHandler) Generate(c *gin.Context) {
	projectID := c.Param("project_id")

	var input ucreport.PetReportInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", err.Error()))
		return
	}

	out, err := h.uc.Generate(c.Request.Context(), projectID, input)
	if err != nil {
		handler.InternalError(c, "generate pet report", err)
		return
	}

	c.JSON(http.StatusOK, out)
}
