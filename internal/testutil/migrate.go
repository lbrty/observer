package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateUp applies all migrations to the given DSN.
// It locates the migrations directory relative to this source file.
func MigrateUp(t *testing.T, dsn string) {
	t.Helper()

	_, thisFile, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(thisFile), "..", "..", "migrations")
	absDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		t.Fatalf("resolve migrations dir: %v", err)
	}
	if _, err := os.Stat(absDir); err != nil {
		t.Fatalf("migrations dir not found: %v", err)
	}

	m, err := migrate.New(fmt.Sprintf("file://%s", absDir), dsn)
	if err != nil {
		t.Fatalf("create migrate instance: %v", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("run migrations: %v", err)
	}
}
