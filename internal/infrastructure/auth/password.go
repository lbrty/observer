package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// PasswordHasher defines hashing and verification for passwords.
type PasswordHasher interface {
	Hash(password string) (hash, salt string, err error)
	Verify(password, hash, salt string) error
}

// ArgonHasher implements PasswordHasher using argon2id.
type ArgonHasher struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

// NewArgonHasher creates an ArgonHasher with sensible defaults.
func NewArgonHasher() *ArgonHasher {
	return &ArgonHasher{
		time:    1,
		memory:  64 * 1024, // 64 MB
		threads: 4,
		keyLen:  32,
	}
}

// Hash derives a base64-encoded hash and salt from the given password.
func (h *ArgonHasher) Hash(password string) (string, string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", "", fmt.Errorf("generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, h.time, h.memory, h.threads, h.keyLen)

	return base64.RawStdEncoding.EncodeToString(hash),
		base64.RawStdEncoding.EncodeToString(salt),
		nil
}

// Verify checks that password matches the stored hash+salt.
func (h *ArgonHasher) Verify(password, hashStr, saltStr string) error {
	salt, err := base64.RawStdEncoding.DecodeString(saltStr)
	if err != nil {
		return fmt.Errorf("decode salt: %w", err)
	}

	expected, err := base64.RawStdEncoding.DecodeString(hashStr)
	if err != nil {
		return fmt.Errorf("decode hash: %w", err)
	}

	computed := argon2.IDKey([]byte(password), salt, h.time, h.memory, h.threads, h.keyLen)

	if subtle.ConstantTimeCompare(computed, expected) == 1 {
		return nil
	}

	return fmt.Errorf("invalid password")
}
