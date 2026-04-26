package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
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
		handler.InternalError(c, "list countries", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"countries": out})
}

// Get handles GET /admin/countries/:id.
func (h *CountryHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /admin/countries.
func (h *CountryHandler) Create(c *gin.Context) {
	input, ok := handler.BindJSON[ucadmin.CreateCountryInput](c)
	if !ok {
		return
	}
	out, err := h.uc.Create(c.Request.Context(), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /admin/countries/:id.
func (h *CountryHandler) Update(c *gin.Context) {
	input, ok := handler.BindJSON[ucadmin.UpdateCountryInput](c)
	if !ok {
		return
	}
	out, err := h.uc.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /admin/countries/:id.
func (h *CountryHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "country deleted"})
}
