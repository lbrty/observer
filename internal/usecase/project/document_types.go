package project

import (
	"time"

	"github.com/lbrty/observer/internal/domain/document"
)

// DocumentDTO is the project-scoped document metadata representation.
type DocumentDTO struct {
	ID               string     `json:"id"`
	PersonID         string     `json:"person_id"`
	ProjectID        string     `json:"project_id"`
	UploadedBy       *string    `json:"uploaded_by,omitempty"`
	EncryptionKeyRef *string    `json:"encryption_key_ref,omitempty"`
	Name             string     `json:"name"`
	MimeType         string     `json:"mime_type"`
	Size             int64      `json:"size"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

// UpdateDocumentInput holds data for updating document metadata.
type UpdateDocumentInput struct {
	Name *string `json:"name"`
}

func documentToDTO(d *document.Document) DocumentDTO {
	dto := DocumentDTO{
		ID:               d.ID,
		PersonID:         d.PersonID,
		ProjectID:        d.ProjectID,
		UploadedBy:       d.UploadedBy,
		EncryptionKeyRef: d.EncryptionKeyRef,
		Name:             d.Name,
		MimeType:         d.MimeType,
		Size:             d.Size,
		CreatedAt:        d.CreatedAt,
	}
	if !d.UpdatedAt.IsZero() {
		t := d.UpdatedAt
		dto.UpdatedAt = &t
	}
	return dto
}
