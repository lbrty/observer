package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/middleware"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

const maxUploadSize = 50 << 20 // 50 MB

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
	if !middleware.CanViewDocumentsFrom(c) {
		c.JSON(http.StatusForbidden, errJSON("errors.document.insufficientPermissions", "insufficient permissions to view documents"))
		return
	}
	personID := c.Param("person_id")
	out, err := h.uc.List(c.Request.Context(), personID)
	if err != nil {
		internalError(c, "list documents", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"documents": out})
}

// Get handles GET /projects/:project_id/documents/:id.
func (h *DocumentHandler) Get(c *gin.Context) {
	if !middleware.CanViewDocumentsFrom(c) {
		c.JSON(http.StatusForbidden, errJSON("errors.document.insufficientPermissions", "insufficient permissions to view documents"))
		return
	}
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Upload handles POST /projects/:project_id/people/:person_id/documents (multipart).
func (h *DocumentHandler) Upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", "file is required"))
		return
	}
	defer file.Close()

	projectID := c.Param("project_id")
	personID := c.Param("person_id")
	userID, _ := middleware.UserIDFrom(c)

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	out, err := h.uc.Upload(
		c.Request.Context(),
		projectID,
		personID,
		userID.String(),
		header.Filename,
		mimeType,
		header.Size,
		file,
	)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Download handles GET /projects/:project_id/documents/:id/download.
func (h *DocumentHandler) Download(c *gin.Context) {
	if !middleware.CanViewDocumentsFrom(c) {
		c.JSON(http.StatusForbidden, errJSON("errors.document.insufficientPermissions", "insufficient permissions to view documents"))
		return
	}

	doc, rc, err := h.uc.Download(c.Request.Context(), c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	defer rc.Close()

	c.Header("Content-Disposition", "attachment; filename=\""+doc.Name+"\"")
	c.DataFromReader(http.StatusOK, doc.Size, doc.MimeType, rc, nil)
}

// Stream handles GET /projects/:project_id/documents/:id/stream.
func (h *DocumentHandler) Stream(c *gin.Context) {
	if !middleware.CanViewDocumentsFrom(c) {
		c.JSON(http.StatusForbidden, errJSON("errors.document.insufficientPermissions", "insufficient permissions to view documents"))
		return
	}

	doc, rc, err := h.uc.Download(c.Request.Context(), c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	defer rc.Close()

	c.Header("Content-Disposition", "inline; filename=\""+doc.Name+"\"")
	c.DataFromReader(http.StatusOK, doc.Size, doc.MimeType, rc, nil)
}

// Thumbnail handles GET /projects/:project_id/documents/:id/thumbnail.
func (h *DocumentHandler) Thumbnail(c *gin.Context) {
	if !middleware.CanViewDocumentsFrom(c) {
		c.JSON(http.StatusForbidden, errJSON("errors.document.insufficientPermissions", "insufficient permissions to view documents"))
		return
	}

	_, rc, err := h.uc.Thumbnail(c.Request.Context(), c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	defer rc.Close()

	c.Header("Cache-Control", "public, max-age=86400")
	c.DataFromReader(http.StatusOK, -1, "image/jpeg", rc, nil)
}

// Update handles PATCH /projects/:project_id/documents/:id.
func (h *DocumentHandler) Update(c *gin.Context) {
	var input ucproject.UpdateDocumentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	out, err := h.uc.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /projects/:project_id/documents/:id.
func (h *DocumentHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "document deleted"})
}
