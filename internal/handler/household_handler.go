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
func (h *HouseholdHandler) List(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.ListHouseholdsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.uc.List(c.Request.Context(), projectID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, out)
}

// Get handles GET /projects/:project_id/households/:id.
func (h *HouseholdHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/households.
func (h *HouseholdHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.CreateHouseholdInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
func (h *HouseholdHandler) Update(c *gin.Context) {
	var input ucproject.UpdateHouseholdInput
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

// Delete handles DELETE /projects/:project_id/households/:id.
func (h *HouseholdHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "household deleted"})
}

// AddMember handles POST /projects/:project_id/households/:id/members.
func (h *HouseholdHandler) AddMember(c *gin.Context) {
	householdID := c.Param("id")
	var input ucproject.AddMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, household.ErrMemberNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, household.ErrMemberExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
