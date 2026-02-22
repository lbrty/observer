package admin

import "time"

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
