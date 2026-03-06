//go:build production

package migrations

import (
	"embed"
	"io/fs"
)

//go:embed *.sql
var sqlFS embed.FS

// FS returns the embedded migrations filesystem.
func FS() (fs.FS, error) {
	return sqlFS, nil
}

// Embedded reports whether migrations are embedded in this build.
func Embedded() bool {
	return true
}
