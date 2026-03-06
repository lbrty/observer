package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/middleware"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

// NoteHandler exposes person note HTTP endpoints.
type NoteHandler struct {
	uc *ucproject.NoteUseCase
}

// NewNoteHandler creates a NoteHandler.
func NewNoteHandler(uc *ucproject.NoteUseCase) *NoteHandler {
	return &NoteHandler{uc: uc}
}

// List handles GET /projects/:project_id/people/:person_id/notes.
// @Summary List notes for a person
// @Tags project-notes
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param person_id path string true "Person ID"
// @Success 200 {object} object "Wrapper with notes array"
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people/{person_id}/notes [get]
func (h *NoteHandler) List(c *gin.Context) {
	personID := c.Param("person_id")
	out, err := h.uc.List(c.Request.Context(), personID)
	if err != nil {
		internalError(c, "list notes", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"notes": out})
}

// Create handles POST /projects/:project_id/people/:person_id/notes.
// @Summary Create a note for a person
// @Tags project-notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param person_id path string true "Person ID"
// @Param input body ucproject.CreateNoteInput true "Note payload"
// @Success 201 {object} ucproject.NoteDTO
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people/{person_id}/notes [post]
func (h *NoteHandler) Create(c *gin.Context) {
	personID := c.Param("person_id")
	var input ucproject.CreateNoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	userID, _ := middleware.UserIDFrom(c)
	out, err := h.uc.Create(c.Request.Context(), personID, userID.String(), input)
	if err != nil {
		internalError(c, "create note", err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /projects/:project_id/people/:person_id/notes/:id.
func (h *NoteHandler) Update(c *gin.Context) {
	var input ucproject.UpdateNoteInput
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

// Delete handles DELETE /projects/:project_id/people/:person_id/notes/:id.
// @Summary Delete a note
// @Tags project-notes
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param person_id path string true "Person ID"
// @Param id path string true "Note ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/people/{person_id}/notes/{id} [delete]
func (h *NoteHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "note deleted"})
}

