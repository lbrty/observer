package pet

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/repository"
)

type petTagRepo struct {
	db *sqlx.DB
}

// NewTag creates a PetTagRepository.
func NewTag(db *sqlx.DB) repository.PetTagRepository {
	return &petTagRepo{db: db}
}

func (r *petTagRepo) List(ctx context.Context, petID string) ([]string, error) {
	const q = `SELECT tag_id FROM pet_tags WHERE pet_id = $1 ORDER BY tag_id`
	rows, err := r.db.QueryContext(ctx, q, petID)
	if err != nil {
		return nil, fmt.Errorf("list pet tags: %w", err)
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

func (r *petTagRepo) ListBulk(ctx context.Context, entityIDs []string) (map[string][]string, error) {
	if len(entityIDs) == 0 {
		return map[string][]string{}, nil
	}

	q, args, err := repository.BuildBulkTagQuery("pet_tags", "pet_id", entityIDs)
	if err != nil {
		return nil, fmt.Errorf("build list bulk: %w", err)
	}
	return repository.QueryBulkTags(ctx, r.db, q, args)
}

func (r *petTagRepo) ReplaceAll(ctx context.Context, petID string, tagIDs []string) error {
	return repository.ReplaceAllJunction(ctx, r.db, "pet_tags", "pet_id", "tag_id", petID, tagIDs)
}
