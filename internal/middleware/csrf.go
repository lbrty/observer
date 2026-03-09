package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CSRFTokenHeader = "X-CSRF-Token"
	CSRFTokenCookie = "csrf_token"
)

var csrfSafeMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodOptions: true,
}

// csrfExemptPaths are auth bootstrap endpoints that run before a CSRF cookie exists.
// They are protected by credentials and rate limiting instead.
var csrfExemptPaths = map[string]bool{
	"/auth/login":    true,
	"/auth/register": true,
	"/auth/refresh":  true,
	"/auth/mfa":      true,
}

// CSRFProtection validates state-changing requests carry an X-CSRF-Token header
// matching the csrf_token cookie (double-submit cookie pattern).
// The csrf_token cookie is set at login (HttpOnly: false) so JS can read it.
// A cross-site attacker can force cookie sending but cannot read cookies or set custom headers.
func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		if csrfSafeMethods[c.Request.Method] || csrfExemptPaths[c.FullPath()] {
			c.Next()
			return
		}

		cookie, err := c.Cookie(CSRFTokenCookie)
		header := c.GetHeader(CSRFTokenHeader)

		if err != nil || cookie == "" || header == "" ||
			subtle.ConstantTimeCompare([]byte(cookie), []byte(header)) != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid or missing CSRF token",
				"code":  "errors.auth.invalidCSRFToken",
			})
			return
		}

		c.Next()
	}
}
