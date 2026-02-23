package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/domain/document"
)

type documentRepo struct {
	db *sqlx.DB
}

// NewDocumentRepository creates a DocumentRepository.
func NewDocumentRepository(db *sqlx.DB) DocumentRepository {
	return &documentRepo{db: db}
}

func (r *documentRepo) List(ctx context.Context, personID string) ([]*document.Document, error) {
	const q = `SELECT id, person_id, project_id, uploaded_by, encryption_key_ref, name, path, mime_type, size, created_at
		FROM documents WHERE person_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, q, personID)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	defer rows.Close()

	var out []*document.Document
	for rows.Next() {
		var d document.Document
		if err := rows.Scan(&d.ID, &d.PersonID, &d.ProjectID, &d.UploadedBy, &d.EncryptionKeyRef,
			&d.Name, &d.Path, &d.MimeType, &d.Size, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan document: %w", err)
		}
		d.CreatedAt = d.CreatedAt.UTC()
		out = append(out, &d)
	}
	return out, rows.Err()
}

func (r *documentRepo) GetByID(ctx context.Context, id string) (*document.Document, error) {
	const q = `SELECT id, person_id, project_id, uploaded_by, encryption_key_ref, name, path, mime_type, size, created_at
		FROM documents WHERE id = $1`
	var d document.Document
	err := r.db.QueryRowContext(ctx, q, id).Scan(&d.ID, &d.PersonID, &d.ProjectID, &d.UploadedBy, &d.EncryptionKeyRef,
		&d.Name, &d.Path, &d.MimeType, &d.Size, &d.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, document.ErrDocumentNotFound
		}
		return nil, fmt.Errorf("get document: %w", err)
	}
	d.CreatedAt = d.CreatedAt.UTC()
	return &d, nil
}

func (r *documentRepo) Create(ctx context.Context, d *document.Document) error {
	const q = `INSERT INTO documents (id, person_id, project_id, uploaded_by, encryption_key_ref, name, path, mime_type, size, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	d.CreatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, q, d.ID, d.PersonID, d.ProjectID, d.UploadedBy, d.EncryptionKeyRef,
		d.Name, d.Path, d.MimeType, d.Size, d.CreatedAt)
	if err != nil {
		return fmt.Errorf("create document: %w", err)
	}
	return nil
}

func (r *documentRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM documents WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return document.ErrDocumentNotFound
	}
	return nil
}
