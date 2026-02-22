package project

import "errors"

var (
	ErrProjectNotFound    = errors.New("project not found")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrInvalidProjectRole = errors.New("invalid project role")
	ErrPermissionExists   = errors.New("permission already exists for this user and project")
)
