package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/domain/reference"
)

// categoryRepo is a PostgreSQL-backed category repository.
type categoryRepo struct {
	db *sqlx.DB
}

// NewCategoryRepository creates a CategoryRepository.
func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) List(ctx context.Context) ([]*reference.Category, error) {
	const q = `SELECT id, name, description, created_at, updated_at FROM categories ORDER BY name`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	var out []*reference.Category
	for rows.Next() {
		var c reference.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		TimesToUTC(&c.CreatedAt, &c.UpdatedAt)
		out = append(out, &c)
	}
	return out, rows.Err()
}

func (r *categoryRepo) GetByID(ctx context.Context, id string) (*reference.Category, error) {
	const q = `SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1`
	var c reference.Category
	err := r.db.QueryRowContext(ctx, q, id).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, reference.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("get category: %w", err)
	}
	TimesToUTC(&c.CreatedAt, &c.UpdatedAt)
	return &c, nil
}

func (r *categoryRepo) Create(ctx context.Context, c *reference.Category) error {
	const q = `INSERT INTO categories (id, name, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q, c.ID, c.Name, c.Description, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		if IsUniqueViolation(err) {
			return reference.ErrCategoryNameExists
		}
		return fmt.Errorf("create category: %w", err)
	}
	return nil
}

func (r *categoryRepo) Update(ctx context.Context, c *reference.Category) error {
	const q = `UPDATE categories SET name=$2, description=$3, updated_at=$4 WHERE id=$1`
	c.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, c.ID, c.Name, c.Description, c.UpdatedAt)
	if err != nil {
		if IsUniqueViolation(err) {
			return reference.ErrCategoryNameExists
		}
		return fmt.Errorf("update category: %w", err)
	}
	return CheckRowsAffected(res, reference.ErrCategoryNotFound)
}

func (r *categoryRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM categories WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return CheckRowsAffected(res, reference.ErrCategoryNotFound)
}
