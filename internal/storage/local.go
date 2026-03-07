package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage stores files on the local filesystem.
type LocalStorage struct {
	root string
}

// NewLocalStorage creates a LocalStorage rooted at the given directory.
// The directory is created if it does not exist.
func NewLocalStorage(root string) (*LocalStorage, error) {
	if err := os.MkdirAll(root, 0o750); err != nil {
		return nil, fmt.Errorf("create storage root %q: %w", root, err)
	}
	return &LocalStorage{root: root}, nil
}

func (s *LocalStorage) resolve(path string) string {
	return filepath.Join(s.root, filepath.Clean("/"+path))
}

func (s *LocalStorage) Save(_ context.Context, path string, r io.Reader) error {
	full := s.resolve(path)
	if err := os.MkdirAll(filepath.Dir(full), 0o750); err != nil {
		return fmt.Errorf("create parent dirs: %w", err)
	}

	f, err := os.Create(full)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		os.Remove(full)
		return fmt.Errorf("write file: %w", err)
	}
	return f.Close()
}

func (s *LocalStorage) Open(_ context.Context, path string) (io.ReadCloser, error) {
	f, err := os.Open(s.resolve(path))
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	return f, nil
}

func (s *LocalStorage) Delete(_ context.Context, path string) error {
	if err := os.Remove(s.resolve(path)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file: %w", err)
	}
	return nil
}
