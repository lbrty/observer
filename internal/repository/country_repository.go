package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/lbrty/observer/internal/domain/reference"
)

// countryRepo is a PostgreSQL-backed country repository.
type countryRepo struct {
	db *sqlx.DB
}

// NewCountryRepository creates a CountryRepository.
func NewCountryRepository(db *sqlx.DB) CountryRepository {
	return &countryRepo{db: db}
}

func (r *countryRepo) List(ctx context.Context) ([]*reference.Country, error) {
	const q = `SELECT id, name, code, created_at, updated_at FROM countries ORDER BY name`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list countries: %w", err)
	}
	defer rows.Close()

	var out []*reference.Country
	for rows.Next() {
		var c reference.Country
		if err := rows.Scan(&c.ID, &c.Name, &c.Code, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan country: %w", err)
		}
		c.CreatedAt = c.CreatedAt.UTC()
		c.UpdatedAt = c.UpdatedAt.UTC()
		out = append(out, &c)
	}
	return out, rows.Err()
}

func (r *countryRepo) GetByID(ctx context.Context, id string) (*reference.Country, error) {
	const q = `SELECT id, name, code, created_at, updated_at FROM countries WHERE id = $1`
	var c reference.Country
	err := r.db.QueryRowContext(ctx, q, id).Scan(&c.ID, &c.Name, &c.Code, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, reference.ErrCountryNotFound
		}
		return nil, fmt.Errorf("get country: %w", err)
	}
	c.CreatedAt = c.CreatedAt.UTC()
	c.UpdatedAt = c.UpdatedAt.UTC()
	return &c, nil
}

func (r *countryRepo) Create(ctx context.Context, c *reference.Country) error {
	const q = `INSERT INTO countries (id, name, code, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q, c.ID, c.Name, c.Code, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return reference.ErrCountryCodeExists
		}
		return fmt.Errorf("create country: %w", err)
	}
	return nil
}

func (r *countryRepo) Update(ctx context.Context, c *reference.Country) error {
	const q = `UPDATE countries SET name=$2, code=$3, updated_at=$4 WHERE id=$1`
	c.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, c.ID, c.Name, c.Code, c.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return reference.ErrCountryCodeExists
		}
		return fmt.Errorf("update country: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrCountryNotFound
	}
	return nil
}

func (r *countryRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM countries WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete country: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrCountryNotFound
	}
	return nil
}
