package project

import "errors"

var (
	ErrProjectNotFound    = errors.New("project not found")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrPermissionNotFound = errors.New("permission not found")
)
