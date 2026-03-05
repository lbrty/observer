package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// errJSON builds a JSON error response with both a human-readable message and
// a machine-readable code that maps to a frontend i18n key.
func errJSON(code, msg string) gin.H {
	return gin.H{"error": msg, "code": code}
}

// internalError logs the error with context and returns a 500 JSON response.
func internalError(c *gin.Context, msg string, err error) {
	slog.Error(msg, slog.String("error", err.Error()), slog.String("path", c.Request.URL.Path), slog.String("method", c.Request.Method))
	c.JSON(http.StatusInternalServerError, errJSON("errors.internal", "internal server error"))
}
