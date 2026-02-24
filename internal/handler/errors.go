package handler

import "github.com/gin-gonic/gin"

// errJSON builds a JSON error response with both a human-readable message and
// a machine-readable code that maps to a frontend i18n key.
func errJSON(code, msg string) gin.H {
	return gin.H{"error": msg, "code": code}
}
