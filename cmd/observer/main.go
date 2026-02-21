package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/lbrty/observer/cmd/observer/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "observer",
	Short: "Observer - IDP platform",
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(cmd.ServeCmd)
	rootCmd.AddCommand(cmd.MigrateCmd)
	rootCmd.AddCommand(cmd.KeygenCmd)
}

func initConfig() {
	_ = godotenv.Load()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
