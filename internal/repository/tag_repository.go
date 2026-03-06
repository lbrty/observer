package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

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
	const q = `SELECT id, project_id, name, color, created_at FROM tags WHERE project_id = $1 ORDER BY name`
	rows, err := r.db.QueryContext(ctx, q, projectID)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	defer rows.Close()

	var out []*tag.Tag
	for rows.Next() {
		var t tag.Tag
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Name, &t.Color, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		TimesToUTC(&t.CreatedAt)
		out = append(out, &t)
	}
	return out, rows.Err()
}

func (r *tagRepo) GetByID(ctx context.Context, id string) (*tag.Tag, error) {
	const q = `SELECT id, project_id, name, color, created_at FROM tags WHERE id = $1`
	var t tag.Tag
	err := r.db.QueryRowContext(ctx, q, id).Scan(&t.ID, &t.ProjectID, &t.Name, &t.Color, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, tag.ErrTagNotFound
		}
		return nil, fmt.Errorf("get tag: %w", err)
	}
	TimesToUTC(&t.CreatedAt)
	return &t, nil
}

func (r *tagRepo) Create(ctx context.Context, t *tag.Tag) error {
	const q = `INSERT INTO tags (id, project_id, name, color, created_at) VALUES ($1, $2, $3, $4, $5)`
	t.CreatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, q, t.ID, t.ProjectID, t.Name, t.Color, t.CreatedAt)
	if err != nil {
		if IsUniqueViolation(err) {
			return tag.ErrTagNameExists
		}
		return fmt.Errorf("create tag: %w", err)
	}
	return nil
}

func (r *tagRepo) Update(ctx context.Context, t *tag.Tag) error {
	const q = `UPDATE tags SET name = $1, color = $2 WHERE id = $3`
	res, err := r.db.ExecContext(ctx, q, t.Name, t.Color, t.ID)
	if err != nil {
		if IsUniqueViolation(err) {
			return tag.ErrTagNameExists
		}
		return fmt.Errorf("update tag: %w", err)
	}
	return CheckRowsAffected(res, tag.ErrTagNotFound)
}

func (r *tagRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM tags WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}
	return CheckRowsAffected(res, tag.ErrTagNotFound)
}
