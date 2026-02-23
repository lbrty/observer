package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/pet"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

// PetHandler exposes pet CRUD HTTP endpoints.
type PetHandler struct {
	uc *ucproject.PetUseCase
}

// NewPetHandler creates a PetHandler.
func NewPetHandler(uc *ucproject.PetUseCase) *PetHandler {
	return &PetHandler{uc: uc}
}

// List handles GET /projects/:project_id/pets.
func (h *PetHandler) List(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.ListPetsInput
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

// Get handles GET /projects/:project_id/pets/:id.
func (h *PetHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/pets.
func (h *PetHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.CreatePetInput
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

// Update handles PATCH /projects/:project_id/pets/:id.
func (h *PetHandler) Update(c *gin.Context) {
	var input ucproject.UpdatePetInput
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

// Delete handles DELETE /projects/:project_id/pets/:id.
func (h *PetHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pet deleted"})
}

func (h *PetHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, pet.ErrPetNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
