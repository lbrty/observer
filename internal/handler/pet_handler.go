package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

// PetHandler exposes pet CRUD HTTP endpoints.
type PetHandler struct {
	uc    *ucproject.PetUseCase
	tagUC *ucproject.PetTagUseCase
}

// NewPetHandler creates a PetHandler.
func NewPetHandler(uc *ucproject.PetUseCase, tagUC *ucproject.PetTagUseCase) *PetHandler {
	return &PetHandler{uc: uc, tagUC: tagUC}
}

// List handles GET /projects/:project_id/pets.
// @Summary List pets in a project
// @Tags project-pets
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} ucproject.ListPetsOutput
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/pets [get]
func (h *PetHandler) List(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.ListPetsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.List(c.Request.Context(), projectID, input)
	if err != nil {
		internalError(c, "list pets", err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Get handles GET /projects/:project_id/pets/:id.
// @Summary Get a pet by ID
// @Tags project-pets
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param id path string true "Pet ID"
// @Success 200 {object} ucproject.PetDTO
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/pets/{id} [get]
func (h *PetHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/pets.
// @Summary Create a pet in a project
// @Tags project-pets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param input body ucproject.CreatePetInput true "Pet payload"
// @Success 201 {object} ucproject.PetDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/pets [post]
func (h *PetHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.CreatePetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.Create(c.Request.Context(), projectID, input)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /projects/:project_id/pets/:id.
// @Summary Update a pet
// @Tags project-pets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param id path string true "Pet ID"
// @Param input body ucproject.UpdatePetInput true "Update payload"
// @Success 200 {object} ucproject.PetDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/pets/{id} [patch]
func (h *PetHandler) Update(c *gin.Context) {
	var input ucproject.UpdatePetInput
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

// Delete handles DELETE /projects/:project_id/pets/:id.
// @Summary Delete a pet
// @Tags project-pets
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param id path string true "Pet ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/pets/{id} [delete]
func (h *PetHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("project_id"), c.Param("id")); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pet deleted"})
}

// ListTags handles GET /projects/:project_id/pets/:id/tags.
func (h *PetHandler) ListTags(c *gin.Context) {
	ids, err := h.tagUC.List(c.Request.Context(), c.Param("id"))
	if err != nil {
		internalError(c, "list pet tags", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"tag_ids": ids})
}

// ReplaceTags handles PUT /projects/:project_id/pets/:id/tags.
func (h *PetHandler) ReplaceTags(c *gin.Context) {
	var input ucproject.ReplaceIDsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	if err := h.tagUC.Replace(c.Request.Context(), c.Param("id"), input.IDs); err != nil {
		internalError(c, "replace pet tags", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"tag_ids": input.IDs})
}

