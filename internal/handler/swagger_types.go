package handler

import ucadmin "github.com/lbrty/observer/internal/usecase/admin"

// ErrorResponse is a generic error response.
type ErrorResponse struct {
	Error string `json:"error" example:"resource not found"`
}

// MessageResponse is a generic success message.
type MessageResponse struct {
	Message string `json:"message" example:"deleted"`
}

// IDListResponse wraps a list of IDs.
type IDListResponse struct {
	IDs []string `json:"ids"`
}

// PermissionListResponse wraps a list of permissions.
type PermissionListResponse struct {
	Permissions []ucadmin.PermissionMemberDTO `json:"permissions"`
}
