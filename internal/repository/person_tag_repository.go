package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type personTagRepo struct {
	db *sqlx.DB
}

// NewPersonTagRepository creates a PersonTagRepository.
func NewPersonTagRepository(db *sqlx.DB) PersonTagRepository {
	return &personTagRepo{db: db}
}

func (r *personTagRepo) List(ctx context.Context, personID string) ([]string, error) {
	const q = `SELECT tag_id FROM person_tags WHERE person_id = $1 ORDER BY tag_id`
	rows, err := r.db.QueryContext(ctx, q, personID)
	if err != nil {
		return nil, fmt.Errorf("list person tags: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan tag id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *personTagRepo) ReplaceAll(ctx context.Context, personID string, tagIDs []string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM person_tags WHERE person_id = $1`, personID); err != nil {
		return fmt.Errorf("delete person tags: %w", err)
	}

	for _, tagID := range tagIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO person_tags (person_id, tag_id) VALUES ($1, $2)`, personID, tagID); err != nil {
			return fmt.Errorf("insert person tag: %w", err)
		}
	}

	return tx.Commit()
}
