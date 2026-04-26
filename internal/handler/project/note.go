package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
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
func (h *NoteHandler) List(c *gin.Context) {
	personID := c.Param("person_id")
	out, err := h.uc.List(c.Request.Context(), personID)
	if err != nil {
		handler.InternalError(c, "list notes", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"notes": out})
}

// Create handles POST /projects/:project_id/people/:person_id/notes.
func (h *NoteHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	personID := c.Param("person_id")
	input, ok := handler.BindJSON[ucproject.CreateNoteInput](c)
	if !ok {
		return
	}
	userID, _ := middleware.UserIDFrom(c)
	out, err := h.uc.Create(c.Request.Context(), projectID, personID, userID.String(), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /projects/:project_id/people/:person_id/notes/:id.
func (h *NoteHandler) Update(c *gin.Context) {
	input, ok := handler.BindJSON[ucproject.UpdateNoteInput](c)
	if !ok {
		return
	}
	out, err := h.uc.Update(c.Request.Context(), c.Param("project_id"), c.Param("person_id"), c.Param("id"), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /projects/:project_id/people/:person_id/notes/:id.
func (h *NoteHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("project_id"), c.Param("person_id"), c.Param("id")); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "note deleted"})
}
