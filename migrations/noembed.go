//go:build !production

package migrations

import "io/fs"

// FS returns nil when migrations are not embedded.
func FS() (fs.FS, error) {
	return nil, nil
}

// Embedded reports whether migrations are embedded in this build.
func Embedded() bool {
	return false
}
