package spa

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler returns a Gin handler that serves a single-page application
// from the given filesystem. Requests matching a static file are served
// directly; everything else falls back to index.html for client-side routing.
func Handler(fsys fs.FS) gin.HandlerFunc {
	fileServer := http.FileServer(http.FS(fsys))

	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if p == "/" {
			p = "index.html"
		}

		// Try to open the exact file.
		if f, err := fsys.Open(p); err == nil {
			f.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		// Fallback to index.html for client-side routing.
		c.Request.URL.Path = "/"
		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
