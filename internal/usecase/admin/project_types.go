package admin

import (
	"time"

	"github.com/lbrty/observer/internal/domain/project"
	"github.com/lbrty/observer/internal/domain/user"
)

// ProjectDTO is the admin-facing project representation.
type ProjectDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	OwnerID     string    `json:"owner_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProjectInput holds data for creating a project.
type CreateProjectInput struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

// UpdateProjectInput holds data for updating a project.
type UpdateProjectInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

// ListProjectsInput holds filter parameters.
type ListProjectsInput struct {
	OwnerID    *string   `form:"owner_id"`
	Status     *string   `form:"status"`
	Page       int       `form:"page"`
	PerPage    int       `form:"per_page"`
	CallerID   string    `form:"-"`
	CallerRole user.Role `form:"-"`
}

// ListProjectsOutput holds paginated results.
type ListProjectsOutput struct {
	Projects []ProjectDTO `json:"projects"`
	Total    int          `json:"total"`
	Page     int          `json:"page"`
	PerPage  int          `json:"per_page"`
}

func projectToDTO(p *project.Project) ProjectDTO {
	return ProjectDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		Status:      string(p.Status),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
