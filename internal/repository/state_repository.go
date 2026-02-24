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

// stateRepo is a PostgreSQL-backed state repository.
type stateRepo struct {
	db *sqlx.DB
}

// NewStateRepository creates a StateRepository.
func NewStateRepository(db *sqlx.DB) StateRepository {
	return &stateRepo{db: db}
}

func (r *stateRepo) ListAll(ctx context.Context) ([]*reference.State, error) {
	const q = `
		SELECT id, country_id, name, code, conflict_zone, created_at, updated_at
		FROM states ORDER BY name
	`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list all states: %w", err)
	}
	defer rows.Close()

	var out []*reference.State
	for rows.Next() {
		var s reference.State
		if err := rows.Scan(&s.ID, &s.CountryID, &s.Name, &s.Code, &s.ConflictZone, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan state: %w", err)
		}
		s.CreatedAt = s.CreatedAt.UTC()
		s.UpdatedAt = s.UpdatedAt.UTC()
		out = append(out, &s)
	}
	return out, rows.Err()
}

func (r *stateRepo) List(ctx context.Context, countryID string) ([]*reference.State, error) {
	const q = `
		SELECT id, country_id, name, code, conflict_zone, created_at, updated_at
		FROM states WHERE country_id = $1 ORDER BY name
	`
	rows, err := r.db.QueryContext(ctx, q, countryID)
	if err != nil {
		return nil, fmt.Errorf("list states: %w", err)
	}
	defer rows.Close()

	var out []*reference.State
	for rows.Next() {
		var s reference.State
		if err := rows.Scan(&s.ID, &s.CountryID, &s.Name, &s.Code, &s.ConflictZone, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan state: %w", err)
		}
		s.CreatedAt = s.CreatedAt.UTC()
		s.UpdatedAt = s.UpdatedAt.UTC()
		out = append(out, &s)
	}
	return out, rows.Err()
}

func (r *stateRepo) GetByID(ctx context.Context, id string) (*reference.State, error) {
	const q = `SELECT id, country_id, name, code, conflict_zone, created_at, updated_at FROM states WHERE id = $1`
	var s reference.State
	err := r.db.QueryRowContext(ctx, q, id).Scan(&s.ID, &s.CountryID, &s.Name, &s.Code, &s.ConflictZone, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, reference.ErrStateNotFound
		}
		return nil, fmt.Errorf("get state: %w", err)
	}
	s.CreatedAt = s.CreatedAt.UTC()
	s.UpdatedAt = s.UpdatedAt.UTC()
	return &s, nil
}

func (r *stateRepo) Create(ctx context.Context, s *reference.State) error {
	const q = `
		INSERT INTO states (id, country_id, name, code, conflict_zone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now().UTC()
	s.CreatedAt = now
	s.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q, s.ID, s.CountryID, s.Name, s.Code, s.ConflictZone, s.CreatedAt, s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create state: %w", err)
	}
	return nil
}

func (r *stateRepo) Update(ctx context.Context, s *reference.State) error {
	const q = `UPDATE states SET name=$2, code=$3, conflict_zone=$4, updated_at=$5 WHERE id=$1`
	s.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, s.ID, s.Name, s.Code, s.ConflictZone, s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update state: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrStateNotFound
	}
	return nil
}

func (r *stateRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM states WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete state: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return reference.ErrStateNotFound
	}
	return nil
}
