package project

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/gif"
	_ "image/png"
	"io"
	"path"
	"strings"

	"golang.org/x/image/draw"

	"github.com/lbrty/observer/internal/domain/document"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/storage"
	"github.com/lbrty/observer/internal/ulid"
)

// DocumentUseCase handles document metadata and file storage.
type DocumentUseCase struct {
	repo repository.DocumentRepository
	fs   storage.FileStorage
}

// NewDocumentUseCase creates a DocumentUseCase.
func NewDocumentUseCase(repo repository.DocumentRepository, fs storage.FileStorage) *DocumentUseCase {
	return &DocumentUseCase{repo: repo, fs: fs}
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

func storagePath(projectID, personID, docID, name string) string {
	return path.Join(projectID, personID, docID+"_"+name)
}

func thumbnailPath(originalPath string) string {
	return originalPath + "_thumb.jpg"
}

func isImage(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}

const thumbnailWidth = 300

func (uc *DocumentUseCase) generateThumbnail(ctx context.Context, originalPath, thumbPath string) error {
	rc, err := uc.fs.Open(ctx, originalPath)
	if err != nil {
		return fmt.Errorf("open original: %w", err)
	}
	defer rc.Close()

	src, _, err := image.Decode(rc)
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	bounds := src.Bounds()
	origW := bounds.Dx()
	origH := bounds.Dy()
	newW := thumbnailWidth
	newH := origH * newW / origW
	if newH < 1 {
		newH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 80}); err != nil {
		return fmt.Errorf("encode thumbnail: %w", err)
	}

	if err := uc.fs.Save(ctx, thumbPath, &buf); err != nil {
		return fmt.Errorf("save thumbnail: %w", err)
	}
	return nil
}

// Thumbnail returns a reader for the document's thumbnail image.
func (uc *DocumentUseCase) Thumbnail(ctx context.Context, id string) (*DocumentDTO, io.ReadCloser, error) {
	d, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("get document: %w", err)
	}
	if !isImage(d.MimeType) {
		return nil, nil, document.ErrNotImage
	}

	thumbPath := thumbnailPath(d.Path)
	rc, err := uc.fs.Open(ctx, thumbPath)
	if err == nil {
		dto := documentToDTO(d)
		return &dto, rc, nil
	}

	if err := uc.generateThumbnail(ctx, d.Path, thumbPath); err != nil {
		return nil, nil, fmt.Errorf("generate thumbnail: %w", err)
	}

	rc, err = uc.fs.Open(ctx, thumbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("open thumbnail: %w", err)
	}
	dto := documentToDTO(d)
	return &dto, rc, nil
}

// Upload stores the file and creates document metadata.
func (uc *DocumentUseCase) Upload(ctx context.Context, projectID, personID, uploadedBy, filename, mimeType string, size int64, body io.Reader) (*DocumentDTO, error) {
	docID := ulid.NewString()
	fpath := storagePath(projectID, personID, docID, filename)

	if err := uc.fs.Save(ctx, fpath, body); err != nil {
		return nil, fmt.Errorf("save file: %w", err)
	}

	d := &document.Document{
		ID:         docID,
		PersonID:   personID,
		ProjectID:  projectID,
		UploadedBy: &uploadedBy,
		Name:       filename,
		Path:       fpath,
		MimeType:   mimeType,
		Size:       size,
	}
	if err := uc.repo.Create(ctx, d); err != nil {
		_ = uc.fs.Delete(ctx, fpath)
		return nil, fmt.Errorf("create document: %w", err)
	}
	dto := documentToDTO(d)
	return &dto, nil
}

// Download opens the file for a given document.
func (uc *DocumentUseCase) Download(ctx context.Context, id string) (*DocumentDTO, io.ReadCloser, error) {
	d, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("get document: %w", err)
	}
	rc, err := uc.fs.Open(ctx, d.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("open file: %w", err)
	}
	dto := documentToDTO(d)
	return &dto, rc, nil
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

// Delete removes document metadata and the stored file.
func (uc *DocumentUseCase) Delete(ctx context.Context, id string) error {
	d, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get document for delete: %w", err)
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	_ = uc.fs.Delete(ctx, d.Path)
	return nil
}
