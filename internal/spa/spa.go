package spa

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Mount registers a NoRoute handler that serves the embedded SPA:
//   - known static files are served with correct MIME types via http.FileServer
//   - everything else falls back to index.html for client-side routing
func Mount(r *gin.Engine, fsys fs.FS) {
	index, _ := fs.ReadFile(fsys, "index.html")
	fileServer := http.FileServer(http.FS(fsys))

	r.NoRoute(func(c *gin.Context) {
		// fs.FS requires relative paths (no leading slash).
		p := strings.TrimPrefix(c.Request.URL.Path, "/")
		if p == "" {
			p = "index.html"
		}

		if f, err := fsys.Open(p); err == nil {
			f.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// Unknown path → serve index.html so the SPA router takes over.
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})
}
