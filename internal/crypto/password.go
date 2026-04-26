package crypto

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
	saltLen uint32
}

// NewArgonHasher creates an ArgonHasher with sensible defaults.
func NewArgonHasher() *ArgonHasher {
	return &ArgonHasher{
		time:    2,         // OWASP recommends t≥2 at 64 MB memory
		memory:  64 * 1024, // 64 MB
		threads: 4,
		keyLen:  32,
		saltLen: 16,
	}
}

// Hash derives a base64-encoded hash and salt from the given password.
func (h *ArgonHasher) Hash(password string) (string, string, error) {
	salt := make([]byte, h.saltLen)
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
	salt, saltErr := base64.RawStdEncoding.DecodeString(saltStr)
	expected, hashErr := base64.RawStdEncoding.DecodeString(hashStr)

	// Always run argon2 before returning any error. Decode failures are
	// programmer errors (corrupted DB values), not user-supplied inputs, so
	// this isn't a real timing oracle today — but keeping the slow path
	// unconditional ensures that stays true even if the call site changes.
	if saltErr != nil {
		salt = make([]byte, h.saltLen)
	}

	if hashErr != nil {
		expected = make([]byte, h.keyLen)
	}

	computed := argon2.IDKey([]byte(password), salt, h.time, h.memory, h.threads, h.keyLen)
	match := subtle.ConstantTimeCompare(computed, expected) == 1

	if saltErr != nil {
		return fmt.Errorf("decode salt: %w", saltErr)
	}

	if hashErr != nil {
		return fmt.Errorf("decode hash: %w", hashErr)
	}

	if !match {
		return fmt.Errorf("invalid password")
	}

	return nil
}
