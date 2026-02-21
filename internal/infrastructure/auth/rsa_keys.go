package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// RSAKeys holds the RSA key pair used for JWT signing.
type RSAKeys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// LoadRSAKeys loads RSA keys from PEM files.
func LoadRSAKeys(privateKeyPath, publicKeyPath string) (*RSAKeys, error) {
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}

	privateBlock, _ := pem.Decode(privateKeyData)
	if privateBlock == nil {
		return nil, fmt.Errorf("failed to decode private key PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read public key: %w", err)
	}

	publicBlock, _ := pem.Decode(publicKeyData)
	if publicBlock == nil {
		return nil, fmt.Errorf("failed to decode public key PEM")
	}

	publicKeyIface, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	publicKey, ok := publicKeyIface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return &RSAKeys{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}
