package project

import (
	"time"

	"github.com/lbrty/observer/internal/domain/tag"
)

// TagDTO is the project-scoped tag representation.
type TagDTO struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateTagInput holds data for creating a tag.
type CreateTagInput struct {
	Name string `json:"name" binding:"required"`
}

func tagToDTO(t *tag.Tag) TagDTO {
	return TagDTO{
		ID:        t.ID,
		ProjectID: t.ProjectID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
	}
}
