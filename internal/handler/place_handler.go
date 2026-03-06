package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
// @Summary List places by state
// @Tags admin-places
// @Produce json
// @Security BearerAuth
// @Param state_id query string true "State ID"
// @Success 200 {object} object "List of places"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/places [get]
func (h *PlaceHandler) List(c *gin.Context) {
	stateID := c.Query("state_id")
	var (
		out []ucadmin.PlaceDTO
		err error
	)
	if stateID != "" {
		out, err = h.uc.List(c.Request.Context(), stateID)
	} else {
		out, err = h.uc.ListAll(c.Request.Context())
	}
	if err != nil {
		internalError(c, "list places", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"places": out})
}

// Get handles GET /admin/places/:id.
// @Summary Get a place by ID
// @Tags admin-places
// @Produce json
// @Security BearerAuth
// @Param id path string true "Place ID"
// @Success 200 {object} ucadmin.PlaceDTO
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/places/{id} [get]
func (h *PlaceHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /admin/places.
// @Summary Create a new place
// @Tags admin-places
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param state_id query string true "State ID"
// @Param input body ucadmin.CreatePlaceInput true "Place payload"
// @Success 201 {object} ucadmin.PlaceDTO
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/places [post]
func (h *PlaceHandler) Create(c *gin.Context) {
	stateID := c.Query("state_id")
	if stateID == "" {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", "state_id is required"))
		return
	}
	var input ucadmin.CreatePlaceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.Create(c.Request.Context(), stateID, input)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /admin/places/:id.
// @Summary Update a place
// @Tags admin-places
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Place ID"
// @Param input body ucadmin.UpdatePlaceInput true "Update payload"
// @Success 200 {object} ucadmin.PlaceDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/places/{id} [patch]
func (h *PlaceHandler) Update(c *gin.Context) {
	var input ucadmin.UpdatePlaceInput
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

// Delete handles DELETE /admin/places/:id.
// @Summary Delete a place
// @Tags admin-places
// @Produce json
// @Security BearerAuth
// @Param id path string true "Place ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/places/{id} [delete]
func (h *PlaceHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "place deleted"})
}

