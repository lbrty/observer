package project

import (
	"time"

	"github.com/lbrty/observer/internal/domain/document"
)

// DocumentDTO is the project-scoped document metadata representation.
type DocumentDTO struct {
	ID               string    `json:"id"`
	PersonID         string    `json:"person_id"`
	ProjectID        string    `json:"project_id"`
	UploadedBy       *string   `json:"uploaded_by,omitempty"`
	EncryptionKeyRef *string   `json:"encryption_key_ref,omitempty"`
	Name             string    `json:"name"`
	Path             string    `json:"path"`
	MimeType         string    `json:"mime_type"`
	Size             int64     `json:"size"`
	CreatedAt        time.Time `json:"created_at"`
}

// CreateDocumentInput holds data for creating document metadata.
type CreateDocumentInput struct {
	PersonID         string  `json:"person_id" binding:"required"`
	Name             string  `json:"name" binding:"required"`
	Path             string  `json:"path" binding:"required"`
	MimeType         string  `json:"mime_type" binding:"required"`
	Size             int64   `json:"size" binding:"required"`
	EncryptionKeyRef *string `json:"encryption_key_ref"`
}

func documentToDTO(d *document.Document) DocumentDTO {
	return DocumentDTO{
		ID:               d.ID,
		PersonID:         d.PersonID,
		ProjectID:        d.ProjectID,
		UploadedBy:       d.UploadedBy,
		EncryptionKeyRef: d.EncryptionKeyRef,
		Name:             d.Name,
		Path:             d.Path,
		MimeType:         d.MimeType,
		Size:             d.Size,
		CreatedAt:        d.CreatedAt,
	}
}
