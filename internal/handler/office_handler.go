package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/reference"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

// OfficeHandler exposes office CRUD HTTP endpoints.
type OfficeHandler struct {
	uc *ucadmin.OfficeUseCase
}

// NewOfficeHandler creates an OfficeHandler.
func NewOfficeHandler(uc *ucadmin.OfficeUseCase) *OfficeHandler {
	return &OfficeHandler{uc: uc}
}

// List handles GET /admin/offices.
// @Summary List all offices
// @Tags admin-offices
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object "List of offices"
// @Failure 500 {object} ErrorResponse
// @Router /admin/offices [get]
func (h *OfficeHandler) List(c *gin.Context) {
	out, err := h.uc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"offices": out})
}

// Get handles GET /admin/offices/:id.
// @Summary Get an office by ID
// @Tags admin-offices
// @Produce json
// @Security BearerAuth
// @Param id path string true "Office ID"
// @Success 200 {object} ucadmin.OfficeDTO
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/offices/{id} [get]
func (h *OfficeHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /admin/offices.
// @Summary Create a new office
// @Tags admin-offices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body ucadmin.CreateOfficeInput true "Office payload"
// @Success 201 {object} ucadmin.OfficeDTO
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/offices [post]
func (h *OfficeHandler) Create(c *gin.Context) {
	var input ucadmin.CreateOfficeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.Create(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /admin/offices/:id.
// @Summary Update an office
// @Tags admin-offices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Office ID"
// @Param input body ucadmin.UpdateOfficeInput true "Update payload"
// @Success 200 {object} ucadmin.OfficeDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/offices/{id} [patch]
func (h *OfficeHandler) Update(c *gin.Context) {
	var input ucadmin.UpdateOfficeInput
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

// Delete handles DELETE /admin/offices/:id.
// @Summary Delete an office
// @Tags admin-offices
// @Produce json
// @Security BearerAuth
// @Param id path string true "Office ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/offices/{id} [delete]
func (h *OfficeHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "office deleted"})
}

func (h *OfficeHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, reference.ErrOfficeNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
