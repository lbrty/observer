package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/lbrty/observer/internal/domain/tag"
)

type tagRepo struct {
	db *sqlx.DB
}

// NewTagRepository creates a TagRepository.
func NewTagRepository(db *sqlx.DB) TagRepository {
	return &tagRepo{db: db}
}

func (r *tagRepo) List(ctx context.Context, projectID string) ([]*tag.Tag, error) {
	const q = `SELECT id, project_id, name, created_at FROM tags WHERE project_id = $1 ORDER BY name`
	rows, err := r.db.QueryContext(ctx, q, projectID)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	defer rows.Close()

	var out []*tag.Tag
	for rows.Next() {
		var t tag.Tag
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Name, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		t.CreatedAt = t.CreatedAt.UTC()
		out = append(out, &t)
	}
	return out, rows.Err()
}

func (r *tagRepo) GetByID(ctx context.Context, id string) (*tag.Tag, error) {
	const q = `SELECT id, project_id, name, created_at FROM tags WHERE id = $1`
	var t tag.Tag
	err := r.db.QueryRowContext(ctx, q, id).Scan(&t.ID, &t.ProjectID, &t.Name, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, tag.ErrTagNotFound
		}
		return nil, fmt.Errorf("get tag: %w", err)
	}
	t.CreatedAt = t.CreatedAt.UTC()
	return &t, nil
}

func (r *tagRepo) Create(ctx context.Context, t *tag.Tag) error {
	const q = `INSERT INTO tags (id, project_id, name, created_at) VALUES ($1, $2, $3, $4)`
	t.CreatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, q, t.ID, t.ProjectID, t.Name, t.CreatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return tag.ErrTagNameExists
		}
		return fmt.Errorf("create tag: %w", err)
	}
	return nil
}

func (r *tagRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM tags WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return tag.ErrTagNotFound
	}
	return nil
}
