//go:build production

package spa

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var distFS embed.FS

// FS returns the embedded frontend filesystem rooted at dist/.
func FS() (fs.FS, error) {
	return fs.Sub(distFS, "dist")
}

// Enabled reports whether the SPA is embedded in this build.
func Enabled() bool {
	return true
}
