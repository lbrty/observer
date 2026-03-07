package document

import "errors"

var (
	ErrDocumentNotFound = errors.New("document not found")
	ErrNotImage         = errors.New("document is not an image")
)
