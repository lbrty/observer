package tag

import "errors"

var (
	ErrTagNotFound   = errors.New("tag not found")
	ErrTagNameExists = errors.New("tag name already exists in this project")
)
