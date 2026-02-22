package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/reference"
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
func (h *CountryHandler) List(c *gin.Context) {
	out, err := h.uc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"countries": out})
}

// Get handles GET /admin/countries/:id.
func (h *CountryHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /admin/countries.
func (h *CountryHandler) Create(c *gin.Context) {
	var input ucadmin.CreateCountryInput
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

// Update handles PATCH /admin/countries/:id.
func (h *CountryHandler) Update(c *gin.Context) {
	var input ucadmin.UpdateCountryInput
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

// Delete handles DELETE /admin/countries/:id.
func (h *CountryHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "country deleted"})
}

func (h *CountryHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, reference.ErrCountryNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, reference.ErrCountryCodeExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
