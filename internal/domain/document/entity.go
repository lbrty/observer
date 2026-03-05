package document

import "time"

// Document represents file metadata attached to a person.
type Document struct {
	ID               string
	PersonID         string
	ProjectID        string
	UploadedBy       *string
	EncryptionKeyRef *string
	Name             string
	Path             string
	MimeType         string
	Size             int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
