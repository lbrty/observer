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

const docCols = `id, person_id, project_id, uploaded_by, encryption_key_ref, name, path, mime_type, size, created_at, updated_at`

func scanDocument(row interface{ Scan(dest ...any) error }) (*document.Document, error) {
	var d document.Document
	var updatedAt sql.NullTime
	err := row.Scan(&d.ID, &d.PersonID, &d.ProjectID, &d.UploadedBy, &d.EncryptionKeyRef,
		&d.Name, &d.Path, &d.MimeType, &d.Size, &d.CreatedAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	TimesToUTC(&d.CreatedAt)
	if updatedAt.Valid {
		t := updatedAt.Time.UTC()
		d.UpdatedAt = t
	}
	return &d, nil
}

func (r *documentRepo) List(ctx context.Context, personID string) ([]*document.Document, error) {
	q := "SELECT " + docCols + " FROM documents WHERE person_id = $1 ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, q, personID)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	defer rows.Close()

	var out []*document.Document
	for rows.Next() {
		d, err := scanDocument(rows)
		if err != nil {
			return nil, fmt.Errorf("scan document: %w", err)
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *documentRepo) GetByID(ctx context.Context, id string) (*document.Document, error) {
	q := "SELECT " + docCols + " FROM documents WHERE id = $1"
	d, err := scanDocument(r.db.QueryRowContext(ctx, q, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, document.ErrDocumentNotFound
		}
		return nil, fmt.Errorf("get document: %w", err)
	}
	return d, nil
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

func (r *documentRepo) Update(ctx context.Context, d *document.Document) error {
	const q = `UPDATE documents SET name=$2, updated_at=$3 WHERE id=$1`
	d.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, d.ID, d.Name, d.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update document: %w", err)
	}
	return CheckRowsAffected(res, document.ErrDocumentNotFound)
}

func (r *documentRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM documents WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	return CheckRowsAffected(res, document.ErrDocumentNotFound)
}
