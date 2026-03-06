package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

// CountryHandler exposes country CRUD HTTP endpoints.
type CountryHandler struct {
	uc *ucadmin.CountryUseCase
}

// NewCountryHandler creates a CountryHandler.
func NewCountryHandler(uc *ucadmin.CountryUseCase) *CountryHandler {
	return &CountryHandler{uc: uc}
}

// List handles GET /admin/countries.
// @Summary List all countries
// @Tags admin-countries
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object "List of countries"
// @Failure 500 {object} ErrorResponse
// @Router /admin/countries [get]
func (h *CountryHandler) List(c *gin.Context) {
	out, err := h.uc.List(c.Request.Context())
	if err != nil {
		internalError(c, "list countries", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"countries": out})
}

// Get handles GET /admin/countries/:id.
// @Summary Get a country by ID
// @Tags admin-countries
// @Produce json
// @Security BearerAuth
// @Param id path string true "Country ID"
// @Success 200 {object} ucadmin.CountryDTO
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/countries/{id} [get]
func (h *CountryHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /admin/countries.
// @Summary Create a new country
// @Tags admin-countries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body ucadmin.CreateCountryInput true "Country payload"
// @Success 201 {object} ucadmin.CountryDTO
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/countries [post]
func (h *CountryHandler) Create(c *gin.Context) {
	var input ucadmin.CreateCountryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.Create(c.Request.Context(), input)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /admin/countries/:id.
// @Summary Update a country
// @Tags admin-countries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Country ID"
// @Param input body ucadmin.UpdateCountryInput true "Update payload"
// @Success 200 {object} ucadmin.CountryDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/countries/{id} [patch]
func (h *CountryHandler) Update(c *gin.Context) {
	var input ucadmin.UpdateCountryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /admin/countries/:id.
// @Summary Delete a country
// @Tags admin-countries
// @Produce json
// @Security BearerAuth
// @Param id path string true "Country ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/countries/{id} [delete]
func (h *CountryHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "country deleted"})
}

