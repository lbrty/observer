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

	"github.com/lbrty/observer/cmd/observer/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
