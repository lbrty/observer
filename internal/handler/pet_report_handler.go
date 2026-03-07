package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}

	out, err := h.uc.Generate(c.Request.Context(), projectID, input)
	if err != nil {
		internalError(c, "generate pet report", err)
		return
	}

	c.JSON(http.StatusOK, out)
}
