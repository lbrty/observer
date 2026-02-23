package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type personCategoryRepo struct {
	db *sqlx.DB
}

// NewPersonCategoryRepository creates a PersonCategoryRepository.
func NewPersonCategoryRepository(db *sqlx.DB) PersonCategoryRepository {
	return &personCategoryRepo{db: db}
}

func (r *personCategoryRepo) List(ctx context.Context, personID string) ([]string, error) {
	const q = `SELECT category_id FROM person_categories WHERE person_id = $1 ORDER BY category_id`
	rows, err := r.db.QueryContext(ctx, q, personID)
	if err != nil {
		return nil, fmt.Errorf("list person categories: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan category id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *personCategoryRepo) ReplaceAll(ctx context.Context, personID string, categoryIDs []string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM person_categories WHERE person_id = $1`, personID); err != nil {
		return fmt.Errorf("delete person categories: %w", err)
	}

	for _, catID := range categoryIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO person_categories (person_id, category_id) VALUES ($1, $2)`, personID, catID); err != nil {
			return fmt.Errorf("insert person category: %w", err)
		}
	}

	return tx.Commit()
}
