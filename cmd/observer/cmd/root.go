package cmd

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// NewRootCmd builds the Observer CLI explicitly instead of relying on package init side effects.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "observer",
		Short: "Observer - IDP case management platform",
		Long: `Observer is a self-hosted case management platform for organizations
working with internally displaced persons (IDPs) and refugees.

It provides project-scoped people tracking, support records, migration
history, household management, document storage, and reporting — all
behind a dual-level RBAC system (platform role + project role).`,
		PersistentPreRun: func(*cobra.Command, []string) {
			_ = godotenv.Load()
		},
	}

	rootCmd.AddCommand(NewServeCmd())
	rootCmd.AddCommand(NewMigrateCmd())
	rootCmd.AddCommand(NewKeygenCmd())
	rootCmd.AddCommand(NewCreateAdminCmd())
	rootCmd.AddCommand(NewSeedCmd())
	rootCmd.AddCommand(NewSetupCmd())

	return rootCmd
}
