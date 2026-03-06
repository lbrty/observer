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

// officeRepo is a PostgreSQL-backed office repository.
type officeRepo struct {
	db *sqlx.DB
}

// NewOfficeRepository creates an OfficeRepository.
func NewOfficeRepository(db *sqlx.DB) OfficeRepository {
	return &officeRepo{db: db}
}

func (r *officeRepo) List(ctx context.Context) ([]*reference.Office, error) {
	const q = `SELECT id, name, place_id, created_at, updated_at FROM offices ORDER BY name`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list offices: %w", err)
	}
	defer rows.Close()

	var out []*reference.Office
	for rows.Next() {
		var o reference.Office
		if err := rows.Scan(&o.ID, &o.Name, &o.PlaceID, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan office: %w", err)
		}
		TimesToUTC(&o.CreatedAt, &o.UpdatedAt)
		out = append(out, &o)
	}
	return out, rows.Err()
}

func (r *officeRepo) GetByID(ctx context.Context, id string) (*reference.Office, error) {
	const q = `SELECT id, name, place_id, created_at, updated_at FROM offices WHERE id = $1`
	var o reference.Office
	err := r.db.QueryRowContext(ctx, q, id).Scan(&o.ID, &o.Name, &o.PlaceID, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, reference.ErrOfficeNotFound
		}
		return nil, fmt.Errorf("get office: %w", err)
	}
	TimesToUTC(&o.CreatedAt, &o.UpdatedAt)
	return &o, nil
}

func (r *officeRepo) Create(ctx context.Context, o *reference.Office) error {
	const q = `INSERT INTO offices (id, name, place_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	now := time.Now().UTC()
	o.CreatedAt = now
	o.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q, o.ID, o.Name, o.PlaceID, o.CreatedAt, o.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create office: %w", err)
	}
	return nil
}

func (r *officeRepo) Update(ctx context.Context, o *reference.Office) error {
	const q = `UPDATE offices SET name=$2, place_id=$3, updated_at=$4 WHERE id=$1`
	o.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, o.ID, o.Name, o.PlaceID, o.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update office: %w", err)
	}
	return CheckRowsAffected(res, reference.ErrOfficeNotFound)
}

func (r *officeRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM offices WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete office: %w", err)
	}
	return CheckRowsAffected(res, reference.ErrOfficeNotFound)
}
