package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/household"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

// HouseholdHandler exposes household HTTP endpoints.
type HouseholdHandler struct {
	uc *ucproject.HouseholdUseCase
}

// NewHouseholdHandler creates a HouseholdHandler.
func NewHouseholdHandler(uc *ucproject.HouseholdUseCase) *HouseholdHandler {
	return &HouseholdHandler{uc: uc}
}

// List handles GET /projects/:project_id/households.
// @Summary List households in a project
// @Tags project-households
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} ucproject.ListHouseholdsOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/households [get]
func (h *HouseholdHandler) List(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.ListHouseholdsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.List(c.Request.Context(), projectID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errJSON("errors.internal", "internal server error"))
		return
	}
	c.JSON(http.StatusOK, out)
}

// Get handles GET /projects/:project_id/households/:id.
// @Summary Get a household by ID
// @Tags project-households
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param id path string true "Household ID"
// @Success 200 {object} ucproject.HouseholdDTO
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/households/{id} [get]
func (h *HouseholdHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/households.
// @Summary Create a household in a project
// @Tags project-households
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param input body ucproject.CreateHouseholdInput true "Household payload"
// @Success 201 {object} ucproject.HouseholdDTO
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/households [post]
func (h *HouseholdHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.CreateHouseholdInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.Create(c.Request.Context(), projectID, input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /projects/:project_id/households/:id.
// @Summary Update a household
// @Tags project-households
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param id path string true "Household ID"
// @Param input body ucproject.UpdateHouseholdInput true "Update payload"
// @Success 200 {object} ucproject.HouseholdDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/households/{id} [patch]
func (h *HouseholdHandler) Update(c *gin.Context) {
	var input ucproject.UpdateHouseholdInput
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

// Delete handles DELETE /projects/:project_id/households/:id.
// @Summary Delete a household
// @Tags project-households
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param id path string true "Household ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/households/{id} [delete]
func (h *HouseholdHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "household deleted"})
}

// AddMember handles POST /projects/:project_id/households/:id/members.
// @Summary Add a member to a household
// @Tags project-households
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param id path string true "Household ID"
// @Param input body ucproject.AddMemberInput true "Member payload"
// @Success 201 {object} ucproject.HouseholdMemberDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/households/{id}/members [post]
func (h *HouseholdHandler) AddMember(c *gin.Context) {
	householdID := c.Param("id")
	var input ucproject.AddMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.AddMember(c.Request.Context(), householdID, input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// RemoveMember handles DELETE /projects/:project_id/households/:id/members/:person_id.
// @Summary Remove a member from a household
// @Tags project-households
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param id path string true "Household ID"
// @Param person_id path string true "Person ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/households/{id}/members/{person_id} [delete]
func (h *HouseholdHandler) RemoveMember(c *gin.Context) {
	householdID := c.Param("id")
	personID := c.Param("person_id")
	if err := h.uc.RemoveMember(c.Request.Context(), householdID, personID); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "member removed"})
}

func (h *HouseholdHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, household.ErrHouseholdNotFound):
		c.JSON(http.StatusNotFound, errJSON("errors.household.notFound", err.Error()))
	case errors.Is(err, household.ErrMemberNotFound):
		c.JSON(http.StatusNotFound, errJSON("errors.household.memberNotFound", err.Error()))
	case errors.Is(err, household.ErrMemberExists):
		c.JSON(http.StatusConflict, errJSON("errors.household.memberExists", err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, errJSON("errors.internal", "internal server error"))
	}
}
