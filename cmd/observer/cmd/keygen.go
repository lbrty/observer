package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// KeygenCmd generates RSA key pairs for JWT signing.
var KeygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate RSA key pair for JWT signing",
	Long: `Generate an RSA key pair for signing and verifying JWT tokens.

Writes private_key.pem and public_key.pem to the output directory.
Minimum key size is 4096 bits. Set JWT_PRIVATE_KEY_PATH and
JWT_PUBLIC_KEY_PATH in your .env to point to the generated files.`,
	Example: `  # Generate keys in the current directory
  observer keygen

  # Generate 8192-bit keys in the keys/ directory
  observer keygen --bits 8192 --output keys`,
	RunE: runKeygen,
}

func init() {
	KeygenCmd.Flags().Int("bits", 4096, "RSA key size (minimum 4096)")
	KeygenCmd.Flags().String("output", ".", "Output directory for key files")
}

func runKeygen(cmd *cobra.Command, _ []string) error {
	bits, _ := cmd.Flags().GetInt("bits")
	output, _ := cmd.Flags().GetString("output")

	if bits < 4096 {
		return fmt.Errorf("RSA key size must be at least 4096 bits (got %d)", bits)
	}

	fmt.Printf("Generating %d-bit RSA key pair...\n", bits)

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return fmt.Errorf("generate key: %w", err)
	}

	privateKeyPath := filepath.Join(output, "private_key.pem")
	if err := writePrivateKey(privateKey, privateKeyPath); err != nil {
		return err
	}
	fmt.Printf("Private key written to: %s\n", privateKeyPath)

	publicKeyPath := filepath.Join(output, "public_key.pem")
	if err := writePublicKey(&privateKey.PublicKey, publicKeyPath); err != nil {
		return err
	}
	fmt.Printf("Public key written to: %s\n", publicKeyPath)

	return nil
}

func writePrivateKey(key *rsa.PrivateKey, path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open private key file: %w", err)
	}
	defer f.Close()

	return pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
}

func writePublicKey(key *rsa.PublicKey, path string) error {
	pubBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return fmt.Errorf("marshal public key: %w", err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open public key file: %w", err)
	}
	defer f.Close()

	return pem.Encode(f, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})
}
