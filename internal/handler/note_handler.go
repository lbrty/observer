package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/note"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notes": out})
}

// Create handles POST /projects/:project_id/people/:person_id/notes.
func (h *NoteHandler) Create(c *gin.Context) {
	personID := c.Param("person_id")
	var input ucproject.CreateNoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := middleware.UserIDFrom(c)
	out, err := h.uc.Create(c.Request.Context(), personID, userID.String(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Delete handles DELETE /projects/:project_id/people/:person_id/notes/:id.
func (h *NoteHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "note deleted"})
}

func (h *NoteHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, note.ErrNoteNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
