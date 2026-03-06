package spa

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Handler returns a Gin handler that serves a single-page application
// from the given directory. Known API prefixes are skipped. All other
// requests either serve a matching static file or fall back to index.html.
func Handler(dir string) gin.HandlerFunc {
	fs := http.Dir(dir)

	return func(c *gin.Context) {
		p := c.Request.URL.Path

		// Try to serve the exact file.
		if _, err := os.Stat(filepath.Join(dir, filepath.Clean(p))); err == nil {
			http.FileServer(fs).ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		// Fallback to index.html for client-side routing.
		c.File(filepath.Join(dir, "index.html"))
		c.Abort()
	}
}
