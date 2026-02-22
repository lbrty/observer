package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/reference"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

// StateHandler exposes state CRUD HTTP endpoints.
type StateHandler struct {
	uc *ucadmin.StateUseCase
}

// NewStateHandler creates a StateHandler.
func NewStateHandler(uc *ucadmin.StateUseCase) *StateHandler {
	return &StateHandler{uc: uc}
}

// List handles GET /admin/states.
func (h *StateHandler) List(c *gin.Context) {
	countryID := c.Query("country_id")
	if countryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "country_id is required"})
		return
	}
	out, err := h.uc.List(c.Request.Context(), countryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"states": out})
}

// Get handles GET /admin/states/:id.
func (h *StateHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /admin/states.
func (h *StateHandler) Create(c *gin.Context) {
	countryID := c.Query("country_id")
	if countryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "country_id is required"})
		return
	}
	var input ucadmin.CreateStateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.Create(c.Request.Context(), countryID, input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /admin/states/:id.
func (h *StateHandler) Update(c *gin.Context) {
	var input ucadmin.UpdateStateInput
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

// Delete handles DELETE /admin/states/:id.
func (h *StateHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "state deleted"})
}

func (h *StateHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, reference.ErrStateNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
