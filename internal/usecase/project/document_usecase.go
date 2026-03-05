package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/document"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// DocumentUseCase handles document metadata operations.
type DocumentUseCase struct {
	repo repository.DocumentRepository
}

// NewDocumentUseCase creates a DocumentUseCase.
func NewDocumentUseCase(repo repository.DocumentRepository) *DocumentUseCase {
	return &DocumentUseCase{repo: repo}
}

// List returns all documents for a person.
func (uc *DocumentUseCase) List(ctx context.Context, personID string) ([]DocumentDTO, error) {
	docs, err := uc.repo.List(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	dtos := make([]DocumentDTO, len(docs))
	for i, d := range docs {
		dtos[i] = documentToDTO(d)
	}
	return dtos, nil
}

// Get returns a document by ID.
func (uc *DocumentUseCase) Get(ctx context.Context, id string) (*DocumentDTO, error) {
	d, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get document: %w", err)
	}
	dto := documentToDTO(d)
	return &dto, nil
}

// Create creates document metadata. uploadedBy is auto-set from auth context.
func (uc *DocumentUseCase) Create(ctx context.Context, projectID, uploadedBy string, input CreateDocumentInput) (*DocumentDTO, error) {
	d := &document.Document{
		ID:               ulid.NewString(),
		PersonID:         input.PersonID,
		ProjectID:        projectID,
		UploadedBy:       &uploadedBy,
		EncryptionKeyRef: input.EncryptionKeyRef,
		Name:             input.Name,
		Path:             input.Path,
		MimeType:         input.MimeType,
		Size:             input.Size,
	}
	if err := uc.repo.Create(ctx, d); err != nil {
		return nil, fmt.Errorf("create document: %w", err)
	}
	dto := documentToDTO(d)
	return &dto, nil
}

// Update updates document metadata.
func (uc *DocumentUseCase) Update(ctx context.Context, id string, input UpdateDocumentInput) (*DocumentDTO, error) {
	d, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get document for update: %w", err)
	}
	if input.Name != nil {
		d.Name = *input.Name
	}
	if err := uc.repo.Update(ctx, d); err != nil {
		return nil, fmt.Errorf("update document: %w", err)
	}
	dto := documentToDTO(d)
	return &dto, nil
}

// Delete removes document metadata.
func (uc *DocumentUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	return nil
}
