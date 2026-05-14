package project

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/middleware"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

const maxUploadSize = 50 << 20 // 50 MB

// requireDocPerm writes a 403 and returns false when the request lacks document view permission.
func requireDocPerm(c *gin.Context) bool {
	if !middleware.CanViewDocumentsFrom(c) {
		c.JSON(http.StatusForbidden, handler.ErrJSON("errors.document.insufficientPermissions", "insufficient permissions to view documents"))
		return false
	}
	return true
}

// DocumentHandler exposes document HTTP endpoints.
type DocumentHandler struct {
	uc *ucproject.DocumentUseCase
}

// NewDocumentHandler creates a DocumentHandler.
func NewDocumentHandler(uc *ucproject.DocumentUseCase) *DocumentHandler {
	return &DocumentHandler{uc: uc}
}

// List handles GET /projects/:project_id/people/:person_id/documents.
func (h *DocumentHandler) List(c *gin.Context) {
	if !requireDocPerm(c) {
		return
	}
	personID := c.Param("person_id")
	projectID := c.Param("project_id")
	out, err := h.uc.List(c.Request.Context(), projectID, personID)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"documents": out})
}

// Get handles GET /projects/:project_id/documents/:id.
func (h *DocumentHandler) Get(c *gin.Context) {
	if !requireDocPerm(c) {
		return
	}
	out, err := h.uc.Get(c.Request.Context(), c.Param("project_id"), c.Param("id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Upload handles POST /projects/:project_id/people/:person_id/documents (multipart).
func (h *DocumentHandler) Upload(c *gin.Context) {
	if !requireDocPerm(c) {
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", "file is required"))
		return
	}
	defer file.Close()

	filename := filepath.Base(header.Filename)
	if filename == "." || filename == "" {
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", "invalid filename"))
		return
	}

	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	detectedMIME := http.DetectContentType(buf[:n])

	switch detectedMIME {
	case "text/html; charset=utf-8", "application/xhtml+xml", "text/xml; charset=utf-8":
		c.JSON(http.StatusBadRequest, handler.ErrJSON("errors.validation", "file type not permitted"))
		return
	}

	body := io.MultiReader(bytes.NewReader(buf[:n]), file)

	projectID := c.Param("project_id")
	personID := c.Param("person_id")
	userID, _ := middleware.UserIDFrom(c)

	out, err := h.uc.Upload(
		c.Request.Context(),
		projectID,
		personID,
		userID.String(),
		filename,
		detectedMIME,
		header.Size,
		body,
	)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Download handles GET /projects/:project_id/documents/:id/download.
func (h *DocumentHandler) Download(c *gin.Context) {
	if !requireDocPerm(c) {
		return
	}

	doc, rc, err := h.uc.Download(c.Request.Context(), c.Param("project_id"), c.Param("id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	defer rc.Close()

	c.Header("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": doc.Name}))
	c.DataFromReader(http.StatusOK, doc.Size, doc.MimeType, rc, nil)
}

// Stream handles GET /projects/:project_id/documents/:id/stream.
func (h *DocumentHandler) Stream(c *gin.Context) {
	if !requireDocPerm(c) {
		return
	}

	doc, rc, err := h.uc.Download(c.Request.Context(), c.Param("project_id"), c.Param("id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	defer rc.Close()

	c.Header("Content-Disposition", mime.FormatMediaType("inline", map[string]string{"filename": doc.Name}))
	c.Header("Content-Security-Policy", "frame-ancestors 'self'")
	c.Header("X-Frame-Options", "SAMEORIGIN")

	safeMIME := doc.MimeType
	switch {
	case strings.HasPrefix(safeMIME, "image/"),
		strings.HasPrefix(safeMIME, "video/"),
		safeMIME == "application/pdf":
		// keep as-is
	default:
		safeMIME = "application/octet-stream"
	}
	c.DataFromReader(http.StatusOK, doc.Size, safeMIME, rc, nil)
}

// Thumbnail handles GET /projects/:project_id/documents/:id/thumbnail.
func (h *DocumentHandler) Thumbnail(c *gin.Context) {
	if !requireDocPerm(c) {
		return
	}

	_, rc, err := h.uc.Thumbnail(c.Request.Context(), c.Param("project_id"), c.Param("id"))
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	defer rc.Close()

	c.Header("Cache-Control", "public, max-age=86400")
	c.DataFromReader(http.StatusOK, -1, "image/jpeg", rc, nil)
}

// Update handles PATCH /projects/:project_id/documents/:id.
func (h *DocumentHandler) Update(c *gin.Context) {
	if !requireDocPerm(c) {
		return
	}
	input, ok := handler.BindJSON[ucproject.UpdateDocumentInput](c)
	if !ok {
		return
	}
	out, err := h.uc.Update(c.Request.Context(), c.Param("project_id"), c.Param("id"), input)
	if err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /projects/:project_id/documents/:id.
func (h *DocumentHandler) Delete(c *gin.Context) {
	if !requireDocPerm(c) {
		return
	}
	if err := h.uc.Delete(c.Request.Context(), c.Param("project_id"), c.Param("id")); err != nil {
		handler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "document deleted"})
}
