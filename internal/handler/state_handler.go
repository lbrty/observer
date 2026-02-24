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
// @Summary List states by country
// @Tags admin-states
// @Produce json
// @Security BearerAuth
// @Param country_id query string true "Country ID"
// @Success 200 {object} object "List of states"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/states [get]
func (h *StateHandler) List(c *gin.Context) {
	countryID := c.Query("country_id")
	var (
		out []ucadmin.StateDTO
		err error
	)
	if countryID != "" {
		out, err = h.uc.List(c.Request.Context(), countryID)
	} else {
		out, err = h.uc.ListAll(c.Request.Context())
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, errJSON("errors.internal", "internal server error"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"states": out})
}

// Get handles GET /admin/states/:id.
// @Summary Get a state by ID
// @Tags admin-states
// @Produce json
// @Security BearerAuth
// @Param id path string true "State ID"
// @Success 200 {object} ucadmin.StateDTO
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/states/{id} [get]
func (h *StateHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /admin/states.
// @Summary Create a new state
// @Tags admin-states
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param country_id query string true "Country ID"
// @Param input body ucadmin.CreateStateInput true "State payload"
// @Success 201 {object} ucadmin.StateDTO
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/states [post]
func (h *StateHandler) Create(c *gin.Context) {
	countryID := c.Query("country_id")
	if countryID == "" {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", "country_id is required"))
		return
	}
	var input ucadmin.CreateStateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
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
// @Summary Update a state
// @Tags admin-states
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "State ID"
// @Param input body ucadmin.UpdateStateInput true "Update payload"
// @Success 200 {object} ucadmin.StateDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/states/{id} [patch]
func (h *StateHandler) Update(c *gin.Context) {
	var input ucadmin.UpdateStateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
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
// @Summary Delete a state
// @Tags admin-states
// @Produce json
// @Security BearerAuth
// @Param id path string true "State ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/states/{id} [delete]
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
		c.JSON(http.StatusNotFound, errJSON("errors.reference.stateNotFound", err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, errJSON("errors.internal", "internal server error"))
	}
}
