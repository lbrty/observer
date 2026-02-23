package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/domain/document"
	"github.com/lbrty/observer/internal/middleware"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

// DocumentHandler exposes document metadata HTTP endpoints.
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
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to view documents"})
		return
	}
	personID := c.Param("person_id")
	out, err := h.uc.List(c.Request.Context(), personID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"documents": out})
}

// Get handles GET /projects/:project_id/documents/:id.
func (h *DocumentHandler) Get(c *gin.Context) {
	if !middleware.CanViewDocumentsFrom(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to view documents"})
		return
	}
	out, err := h.uc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Create handles POST /projects/:project_id/documents.
func (h *DocumentHandler) Create(c *gin.Context) {
	projectID := c.Param("project_id")
	var input ucproject.CreateDocumentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := middleware.UserIDFrom(c)
	out, err := h.uc.Create(c.Request.Context(), projectID, userID.String(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Delete handles DELETE /projects/:project_id/documents/:id.
func (h *DocumentHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "document deleted"})
}

func (h *DocumentHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, document.ErrDocumentNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
