package cmd

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// SetupCmd runs first-time project setup.
var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "First-time project setup",
	Long: `Run first-time setup for Observer.

Generates a .env file with sensible defaults, creates RSA keys for JWT
signing, and creates required directories. After setup, start Postgres
and Redis, run migrations, create an admin user, and start the server.`,
	Example: `  # Run interactive setup
  observer setup

  # Then start the server
  observer migrate up
  observer create-admin --email admin@example.com --password "s3cure-p4ss"
  observer serve`,
	RunE: runSetup,
}

const defaultEnvContent = `# Server
SERVER_HOST=localhost
SERVER_PORT=9000

# Database
DATABASE_DSN=postgres://observer:observer@localhost:5432/observer?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379/0

# JWT
JWT_PRIVATE_KEY_PATH=keys/private_key.pem
JWT_PUBLIC_KEY_PATH=keys/public_key.pem
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
JWT_ISSUER=observer

# CORS
CORS_ORIGINS=http://localhost:5173

# Cookies
COOKIE_DOMAIN=
COOKIE_SECURE=false
COOKIE_SAME_SITE=lax

# Storage
STORAGE_PATH=data/uploads

# Logging
LOG_LEVEL=info
`

func runSetup(cmd *cobra.Command, _ []string) error {
	reader := bufio.NewReader(cmd.InOrStdin())
	out := cmd.OutOrStdout()

	// Step 1: Write .env
	if err := writeEnvFile(reader, out); err != nil {
		return err
	}

	// Step 2: Create required directories
	for _, dir := range []string{"keys", "data/uploads"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
		fmt.Fprintf(out, "Created directory: %s\n", dir)
	}

	// Step 3: Generate RSA keys
	if err := generateSetupKeys(out); err != nil {
		return err
	}

	// Step 4: Print next steps
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Setup complete!")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Next steps:")
	fmt.Fprintln(out, "  1. Start Postgres and Redis:")
	fmt.Fprintln(out, "     docker compose up -d")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "  2. Run database migrations:")
	fmt.Fprintln(out, "     observer migrate up")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "  3. Create an admin user:")
	fmt.Fprintln(out, `     observer create-admin --email admin@example.com --password "your-password"`)
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "  4. Start the server:")
	fmt.Fprintln(out, "     observer serve")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Useful commands:")
	fmt.Fprintln(out, "  just --list          # Show all available tasks")
	fmt.Fprintln(out, "  observer seed        # Populate with development data")
	fmt.Fprintln(out, "  observer --help      # CLI reference")

	return nil
}

func writeEnvFile(reader *bufio.Reader, out io.Writer) error {
	if _, err := os.Stat(".env"); err == nil {
		fmt.Fprint(out, ".env already exists. Overwrite? [y/N] ")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Fprintln(out, "Skipping .env generation.")
			return nil
		}
	}

	if err := os.WriteFile(".env", []byte(defaultEnvContent), 0644); err != nil {
		return fmt.Errorf("write .env: %w", err)
	}
	fmt.Fprintln(out, "Created .env with default configuration.")
	return nil
}

func generateSetupKeys(out io.Writer) error {
	privPath := "keys/private_key.pem"

	if _, err := os.Stat(privPath); err == nil {
		fmt.Fprintln(out, "RSA keys already exist, skipping generation.")
		return nil
	}

	fmt.Fprintln(out, "Generating 4096-bit RSA key pair...")

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("generate key: %w", err)
	}

	if err := writePrivateKey(privateKey, "keys/private_key.pem"); err != nil {
		return err
	}
	fmt.Fprintln(out, "Private key written to: keys/private_key.pem")

	if err := writePublicKey(&privateKey.PublicKey, "keys/public_key.pem"); err != nil {
		return err
	}
	fmt.Fprintln(out, "Public key written to: keys/public_key.pem")

	return nil
}
