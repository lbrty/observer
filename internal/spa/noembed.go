//go:build !production

package spa

import "io/fs"

// FS returns nil when the SPA is not embedded.
func FS() (fs.FS, error) {
	return nil, nil
}

// Enabled reports whether the SPA is embedded in this build.
func Enabled() bool {
	return false
}
