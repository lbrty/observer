package my

import "time"

// MyProjectDTO represents a project the current user has access to.
type MyProjectDTO struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Description      *string   `json:"description,omitempty"`
	Status           string    `json:"status"`
	Role             string    `json:"role"`
	CanViewContact   bool      `json:"can_view_contact"`
	CanViewPersonal  bool      `json:"can_view_personal"`
	CanViewDocuments bool      `json:"can_view_documents"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// MyProjectsOutput is the response for listing the current user's projects.
type MyProjectsOutput struct {
	Projects []MyProjectDTO `json:"projects"`
}
