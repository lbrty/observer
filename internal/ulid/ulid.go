package ulid

import (
	"crypto/rand"
	"sync"

	"github.com/oklog/ulid/v2"
)

var (
	entropy     = ulid.Monotonic(rand.Reader, 0)
	entropyLock sync.Mutex
)

// New generates a new unique ULID.
func New() ulid.ULID {
	entropyLock.Lock()
	defer entropyLock.Unlock()
	return ulid.MustNew(ulid.Now(), entropy)
}

// NewString generates a new ULID and returns it as a string.
func NewString() string {
	return New().String()
}

// Parse parses a ULID from string.
func Parse(s string) (ulid.ULID, error) {
	return ulid.Parse(s)
}

// IsValid reports whether s is a valid ULID string.
func IsValid(s string) bool {
	_, err := ulid.Parse(s)
	return err == nil
}
