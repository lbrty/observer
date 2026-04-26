package person

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/repository"
)

type personCategoryRepo struct {
	db *sqlx.DB
}

// NewCategory creates a PersonCategoryRepository.
func NewCategory(db *sqlx.DB) repository.PersonCategoryRepository {
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

func (r *personCategoryRepo) ListBulk(ctx context.Context, personIDs []string) (map[string][]string, error) {
	if len(personIDs) == 0 {
		return map[string][]string{}, nil
	}

	sqlStr, args, err := repository.PSQL.
		Select("person_id, category_id").
		From("person_categories").
		Where(sq.Eq{"person_id": personIDs}).
		OrderBy("person_id, category_id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build list bulk: %w", err)
	}

	return repository.QueryBulkTags(ctx, r.db, sqlStr, args)
}

func (r *personCategoryRepo) ReplaceAll(ctx context.Context, personID string, categoryIDs []string) error {
	return repository.ReplaceAllJunction(ctx, r.db, "person_categories", "person_id", "category_id", personID, categoryIDs)
}
