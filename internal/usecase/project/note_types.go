package project

import (
	"time"

	"github.com/lbrty/observer/internal/domain/note"
)

// NoteDTO is the project-scoped person note representation.
type NoteDTO struct {
	ID        string     `json:"id"`
	PersonID  string     `json:"person_id"`
	AuthorID  *string    `json:"author_id,omitempty"`
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// CreateNoteInput holds data for creating a note.
type CreateNoteInput struct {
	Body string `json:"body" binding:"required"`
}

// UpdateNoteInput holds data for updating a note.
type UpdateNoteInput struct {
	Body string `json:"body" binding:"required"`
}

func noteToDTO(n *note.Note) NoteDTO {
	dto := NoteDTO{
		ID:        n.ID,
		PersonID:  n.PersonID,
		AuthorID:  n.AuthorID,
		Body:      n.Body,
		CreatedAt: n.CreatedAt,
	}
	if !n.UpdatedAt.IsZero() {
		t := n.UpdatedAt
		dto.UpdatedAt = &t
	}
	return dto
}
