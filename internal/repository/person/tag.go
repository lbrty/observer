package person

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/repository"
)

type personTagRepo struct {
	db *sqlx.DB
}

// NewTag creates a PersonTagRepository.
func NewTag(db *sqlx.DB) repository.PersonTagRepository {
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

func (r *personTagRepo) ListBulk(ctx context.Context, entityIDs []string) (map[string][]string, error) {
	if len(entityIDs) == 0 {
		return map[string][]string{}, nil
	}

	q, args, err := repository.BuildBulkTagQuery("person_tags", "person_id", entityIDs)
	if err != nil {
		return nil, fmt.Errorf("build list bulk: %w", err)
	}
	return repository.QueryBulkTags(ctx, r.db, q, args)
}

func (r *personTagRepo) ReplaceAll(ctx context.Context, personID string, tagIDs []string) error {
	return repository.ReplaceAllJunction(ctx, r.db, "person_tags", "person_id", "tag_id", personID, tagIDs)
}
