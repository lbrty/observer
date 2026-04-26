package storage

import (
	"context"
	"errors"
	"io"
)

// ErrNotFound is returned by Open when the requested file does not exist.
var ErrNotFound = errors.New("file not found")

// FileStorage defines operations for storing and retrieving files.
type FileStorage interface {
	// Save writes r to path. contentType is used as the object's media type
	// (e.g. for Content-Type headers on S3); local storage ignores it.
	Save(ctx context.Context, path string, r io.Reader, contentType string) error
	// Open returns a reader for the file at path. Returns ErrNotFound if absent.
	Open(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
}
