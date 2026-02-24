// @title Observer API
// @version 1.0
// @description IDP case management platform API.
//
// @host localhost:9000
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer {token}"
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
	rootCmd.AddCommand(cmd.CreateAdminCmd)
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
