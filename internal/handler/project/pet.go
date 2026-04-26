package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
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
func (h *PetHandler) List(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.ListPetsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.List(c.Request.Context(), projectID, input)
	if err != nil {
		handler.InternalError(c, "list pets", err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Get handles GET /projects/:project_id/pets/:id.
func (h *PetHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("project_id"), c.Param("id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/pets.
func (h *PetHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	input, ok := handler.BindJSON[ucproject.CreatePetInput](c)
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

// Update handles PATCH /projects/:project_id/pets/:id.
func (h *PetHandler) Update(c *gin.Context) {
	input, ok := handler.BindJSON[ucproject.UpdatePetInput](c)
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

// Delete handles DELETE /projects/:project_id/pets/:id.
func (h *PetHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("project_id"), c.Param("id")); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pet deleted"})
}

// ListTags handles GET /projects/:project_id/pets/:id/tags.
func (h *PetHandler) ListTags(c *gin.Context) {
	ids, err := h.tagUC.List(c.Request.Context(), c.Param("id"))
	if err != nil {
		handler.InternalError(c, "list pet tags", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"tag_ids": ids})
}

// ReplaceTags handles PUT /projects/:project_id/pets/:id/tags.
func (h *PetHandler) ReplaceTags(c *gin.Context) {
	input, ok := handler.BindJSON[ucproject.ReplaceIDsInput](c)
	if !ok {
		return
	}
	if err := h.tagUC.Replace(c.Request.Context(), c.Param("project_id"), c.Param("id"), input.IDs); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"tag_ids": input.IDs})
}
