package admin

import (
	"time"

	"github.com/lbrty/observer/internal/domain/project"
)

// AssignPermissionInput holds data for assigning a project permission.
type AssignPermissionInput struct {
	UserID           string `json:"user_id" binding:"required"`
	Role             string `json:"role" binding:"required"`
	CanViewContact   bool   `json:"can_view_contact"`
	CanViewPersonal  bool   `json:"can_view_personal"`
	CanViewDocuments bool   `json:"can_view_documents"`
}

// UpdatePermissionInput holds fields for updating a project permission.
type UpdatePermissionInput struct {
	Role             *string `json:"role"`
	CanViewContact   *bool   `json:"can_view_contact"`
	CanViewPersonal  *bool   `json:"can_view_personal"`
	CanViewDocuments *bool   `json:"can_view_documents"`
}

// PermissionDTO is the admin-facing project permission representation.
type PermissionDTO struct {
	ID               string    `json:"id"`
	ProjectID        string    `json:"project_id"`
	UserID           string    `json:"user_id"`
	Role             string    `json:"role"`
	CanViewContact   bool      `json:"can_view_contact"`
	CanViewPersonal  bool      `json:"can_view_personal"`
	CanViewDocuments bool      `json:"can_view_documents"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func permToDTO(p *project.ProjectPermission) PermissionDTO {
	return PermissionDTO{
		ID:               p.ID,
		ProjectID:        p.ProjectID,
		UserID:           p.UserID,
		Role:             string(p.Role),
		CanViewContact:   p.CanViewContact,
		CanViewPersonal:  p.CanViewPersonal,
		CanViewDocuments: p.CanViewDocuments,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}
}

// PermissionMemberDTO is an enriched permission with user details.
type PermissionMemberDTO struct {
	ID               string    `json:"id"`
	ProjectID        string    `json:"project_id"`
	UserID           string    `json:"user_id"`
	UserFirstName    string    `json:"user_first_name"`
	UserLastName     string    `json:"user_last_name"`
	UserEmail        string    `json:"user_email"`
	UserRole         string    `json:"user_role"`
	Role             string    `json:"role"`
	CanViewContact   bool      `json:"can_view_contact"`
	CanViewPersonal  bool      `json:"can_view_personal"`
	CanViewDocuments bool      `json:"can_view_documents"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
