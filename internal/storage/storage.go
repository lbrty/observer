package storage

import (
	"context"
	"io"
)

// FileStorage defines operations for storing and retrieving files.
type FileStorage interface {
	Save(ctx context.Context, path string, r io.Reader) error
	Open(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
}
