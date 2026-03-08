package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/cobra"

	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/migrations"
)

// MigrateCmd groups migration subcommands.
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration management",
	Long: `Manage database schema migrations.

Observer uses forward-only SQL migrations. In production builds, migrations
are embedded in the binary and applied automatically on server start.
For development, use the subcommands below to apply, create, or inspect
migration state.`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	Long: `Apply all pending database migrations in order.

Reads DATABASE_DSN from environment or .env file. In production builds,
uses embedded migrations; in development, reads from the migrations/
directory (configurable with --path).`,
	Example: `  # Apply all pending migrations
  observer migrate up

  # Use a custom migrations directory
  observer migrate up --path ./db/migrations`,
	RunE: runMigrateUp,
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new forward migration file",
	Long: `Create a new .up.sql migration file with the next sequence number.

The sequence number is auto-incremented from the highest existing
migration file, or can be specified explicitly with --seq.`,
	Example: `  # Create a migration (auto-numbered)
  observer migrate create add_audit_log

  # Create with explicit sequence number
  observer migrate create add_audit_log --seq 25`,
	Args: cobra.ExactArgs(1),
	RunE: runMigrateCreate,
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	Long:  `Display the current migration version and dirty state from the database.`,
	Example: `  observer migrate version`,
	RunE: runMigrateVersion,
}

func init() {
	MigrateCmd.AddCommand(migrateUpCmd)
	MigrateCmd.AddCommand(migrateCreateCmd)
	MigrateCmd.AddCommand(migrateVersionCmd)

	migrateUpCmd.Flags().String("path", "migrations", "Path to migrations directory")
	migrateCreateCmd.Flags().String("path", "migrations", "Path to migrations directory")
	migrateCreateCmd.Flags().Uint("seq", 0, "Explicit sequence number (default: auto-increment from highest existing)")
	migrateVersionCmd.Flags().String("path", "migrations", "Path to migrations directory")
}

func newMigrate(cmd *cobra.Command) (*migrate.Migrate, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Use embedded migrations in production builds.
	if migrations.Embedded() {
		fsys, err := migrations.FS()
		if err != nil {
			return nil, fmt.Errorf("embedded migrations: %w", err)
		}
		d, err := iofs.New(fsys, ".")
		if err != nil {
			return nil, fmt.Errorf("iofs source: %w", err)
		}
		return migrate.NewWithSourceInstance("iofs", d, cfg.Database.DSN)
	}

	dir, _ := cmd.Flags().GetString("path")
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	return migrate.New(fmt.Sprintf("file://%s", absDir), cfg.Database.DSN)
}

func runMigrateUp(cmd *cobra.Command, _ []string) error {
	m, err := newMigrate(cmd)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	v, dirty, _ := m.Version()
	fmt.Printf("Migrated to version %d (dirty=%v)\n", v, dirty)
	return nil
}

func runMigrateCreate(cmd *cobra.Command, args []string) error {
	dir, _ := cmd.Flags().GetString("path")
	name := args[0]

	seq, _ := cmd.Flags().GetUint("seq")
	if seq == 0 {
		next, err := nextMigrationSeq(dir)
		if err != nil {
			return err
		}
		seq = next
	}

	filename := filepath.Join(dir, fmt.Sprintf("%06d_%s.up.sql", seq, name))
	if err := os.WriteFile(filename, []byte(""), 0644); err != nil {
		return fmt.Errorf("create %s: %w", filename, err)
	}
	fmt.Printf("Created: %s\n", filename)
	return nil
}

// nextMigrationSeq scans the migrations directory and returns the next sequence number.
func nextMigrationSeq(dir string) (uint, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, fmt.Errorf("read migrations dir: %w", err)
	}

	var max uint
	for _, e := range entries {
		var n uint
		if _, err := fmt.Sscanf(e.Name(), "%06d", &n); err == nil && n > max {
			max = n
		}
	}
	return max + 1, nil
}

func runMigrateVersion(cmd *cobra.Command, _ []string) error {
	m, err := newMigrate(cmd)
	if err != nil {
		return err
	}
	defer m.Close()

	v, dirty, err := m.Version()
	if err != nil {
		return err
	}

	fmt.Printf("Version: %d, Dirty: %v\n", v, dirty)
	return nil
}
