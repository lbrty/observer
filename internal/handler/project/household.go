package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
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
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.List(c.Request.Context(), projectID, input)
	if err != nil {
		handler.InternalError(c, "list households", err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Get handles GET /projects/:project_id/households/:id.
func (h *HouseholdHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("project_id"), c.Param("id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/households.
func (h *HouseholdHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	input, ok := handler.BindJSON[ucproject.CreateHouseholdInput](c)
	if !ok {
		return
	}
	out, err := h.uc.Create(c.Request.Context(), projectID, input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /projects/:project_id/households/:id.
func (h *HouseholdHandler) Update(c *gin.Context) {
	input, ok := handler.BindJSON[ucproject.UpdateHouseholdInput](c)
	if !ok {
		return
	}
	out, err := h.uc.Update(c.Request.Context(), c.Param("project_id"), c.Param("id"), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /projects/:project_id/households/:id.
func (h *HouseholdHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("project_id"), c.Param("id")); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "household deleted"})
}

// AddMember handles POST /projects/:project_id/households/:id/members.
func (h *HouseholdHandler) AddMember(c *gin.Context) {
	input, ok := handler.BindJSON[ucproject.AddMemberInput](c)
	if !ok {
		return
	}
	out, err := h.uc.AddMember(c.Request.Context(), c.Param("project_id"), c.Param("id"), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// RemoveMember handles DELETE /projects/:project_id/households/:id/members/:person_id.
func (h *HouseholdHandler) RemoveMember(c *gin.Context) {
	if err := h.uc.RemoveMember(c.Request.Context(), c.Param("project_id"), c.Param("id"), c.Param("person_id")); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "member removed"})
}
