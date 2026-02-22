package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/reference"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

// PlaceHandler exposes place CRUD HTTP endpoints.
type PlaceHandler struct {
	uc *ucadmin.PlaceUseCase
}

// NewPlaceHandler creates a PlaceHandler.
func NewPlaceHandler(uc *ucadmin.PlaceUseCase) *PlaceHandler {
	return &PlaceHandler{uc: uc}
}

// List handles GET /admin/places.
func (h *PlaceHandler) List(c *gin.Context) {
	stateID := c.Query("state_id")
	if stateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state_id is required"})
		return
	}
	out, err := h.uc.List(c.Request.Context(), stateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"places": out})
}

// Get handles GET /admin/places/:id.
func (h *PlaceHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /admin/places.
func (h *PlaceHandler) Create(c *gin.Context) {
	stateID := c.Query("state_id")
	if stateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state_id is required"})
		return
	}
	var input ucadmin.CreatePlaceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.Create(c.Request.Context(), stateID, input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /admin/places/:id.
func (h *PlaceHandler) Update(c *gin.Context) {
	var input ucadmin.UpdatePlaceInput
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

// Delete handles DELETE /admin/places/:id.
func (h *PlaceHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "place deleted"})
}

func (h *PlaceHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, reference.ErrPlaceNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
