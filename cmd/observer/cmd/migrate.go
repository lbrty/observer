package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"

	"github.com/lbrty/observer/internal/config"
)

// MigrateCmd groups migration subcommands.
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration management",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	RunE:  runMigrateUp,
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new forward migration file",
	Args:  cobra.ExactArgs(1),
	RunE:  runMigrateCreate,
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	RunE:  runMigrateVersion,
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
