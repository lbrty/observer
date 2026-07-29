package storage_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lbrty/observer/internal/storage"
)

func TestLocalStorage_Save_WritesContent(t *testing.T) {
	dir := t.TempDir()
	ls, err := storage.NewLocalStorage(dir)
	require.NoError(t, err)

	err = ls.Save(context.Background(), "sub/file.txt", strings.NewReader("hello world"), "text/plain")
	require.NoError(t, err)

	rc, err := ls.Open(context.Background(), "sub/file.txt")
	require.NoError(t, err)
	defer rc.Close()

	got, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(got))

	info, err := os.Stat(filepath.Join(dir, "sub", "file.txt"))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestLocalStorage_Delete_MissingFileIsNoop(t *testing.T) {
	dir := t.TempDir()
	ls, err := storage.NewLocalStorage(dir)
	require.NoError(t, err)

	err = ls.Delete(context.Background(), "nonexistent/path.txt")
	assert.NoError(t, err)
}

func TestLocalStorage_Open_MissingFileReturnsError(t *testing.T) {
	dir := t.TempDir()
	ls, err := storage.NewLocalStorage(dir)
	require.NoError(t, err)

	_, err = ls.Open(context.Background(), "nonexistent.txt")
	assert.ErrorIs(t, err, storage.ErrNotFound)
}
