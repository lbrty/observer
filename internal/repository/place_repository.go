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

// placeRepo is a PostgreSQL-backed place repository.
type placeRepo struct {
	db *sqlx.DB
}

// NewPlaceRepository creates a PlaceRepository.
func NewPlaceRepository(db *sqlx.DB) PlaceRepository {
	return &placeRepo{db: db}
}

func (r *placeRepo) List(ctx context.Context, stateID string) ([]*reference.Place, error) {
	const q = `
		SELECT id, state_id, name, lat, lon, created_at, updated_at
		FROM places WHERE state_id = $1 ORDER BY name
	`
	rows, err := r.db.QueryContext(ctx, q, stateID)
	if err != nil {
		return nil, fmt.Errorf("list places: %w", err)
	}
	defer rows.Close()

	var out []*reference.Place
	for rows.Next() {
		var p reference.Place
		if err := rows.Scan(&p.ID, &p.StateID, &p.Name, &p.Lat, &p.Lon, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan place: %w", err)
		}
		p.CreatedAt = p.CreatedAt.UTC()
		p.UpdatedAt = p.UpdatedAt.UTC()
		out = append(out, &p)
	}
	return out, rows.Err()
}

func (r *placeRepo) GetByID(ctx context.Context, id string) (*reference.Place, error) {
	const q = `SELECT id, state_id, name, lat, lon, created_at, updated_at FROM places WHERE id = $1`
	var p reference.Place
	err := r.db.QueryRowContext(ctx, q, id).Scan(&p.ID, &p.StateID, &p.Name, &p.Lat, &p.Lon, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, reference.ErrPlaceNotFound
		}
		return nil, fmt.Errorf("get place: %w", err)
	}
	p.CreatedAt = p.CreatedAt.UTC()
	p.UpdatedAt = p.UpdatedAt.UTC()
	return &p, nil
}

func (r *placeRepo) Create(ctx context.Context, p *reference.Place) error {
	const q = `
		INSERT INTO places (id, state_id, name, lat, lon, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q, p.ID, p.StateID, p.Name, p.Lat, p.Lon, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create place: %w", err)
	}
	return nil
}

func (r *placeRepo) Update(ctx context.Context, p *reference.Place) error {
	const q = `UPDATE places SET name=$2, lat=$3, lon=$4, updated_at=$5 WHERE id=$1`
	p.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, p.ID, p.Name, p.Lat, p.Lon, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update place: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrPlaceNotFound
	}
	return nil
}

func (r *placeRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM places WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete place: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrPlaceNotFound
	}
	return nil
}
