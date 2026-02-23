package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/domain/household"
)

type householdRepo struct {
	db *sqlx.DB
}

// NewHouseholdRepository creates a HouseholdRepository.
func NewHouseholdRepository(db *sqlx.DB) HouseholdRepository {
	return &householdRepo{db: db}
}

func (r *householdRepo) List(ctx context.Context, projectID string, page, perPage int) ([]*household.Household, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM households WHERE project_id = $1`, projectID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count households: %w", err)
	}

	offset := (page - 1) * perPage
	const q = `SELECT id, project_id, reference_number, head_person_id, created_at, updated_at
		FROM households WHERE project_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, q, projectID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list households: %w", err)
	}
	defer rows.Close()

	var out []*household.Household
	for rows.Next() {
		var h household.Household
		if err := rows.Scan(&h.ID, &h.ProjectID, &h.ReferenceNumber, &h.HeadPersonID, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan household: %w", err)
		}
		h.CreatedAt = h.CreatedAt.UTC()
		h.UpdatedAt = h.UpdatedAt.UTC()
		out = append(out, &h)
	}
	return out, total, rows.Err()
}

func (r *householdRepo) GetByID(ctx context.Context, id string) (*household.Household, error) {
	const q = `SELECT id, project_id, reference_number, head_person_id, created_at, updated_at FROM households WHERE id = $1`
	var h household.Household
	err := r.db.QueryRowContext(ctx, q, id).Scan(&h.ID, &h.ProjectID, &h.ReferenceNumber, &h.HeadPersonID, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, household.ErrHouseholdNotFound
		}
		return nil, fmt.Errorf("get household: %w", err)
	}
	h.CreatedAt = h.CreatedAt.UTC()
	h.UpdatedAt = h.UpdatedAt.UTC()
	return &h, nil
}

func (r *householdRepo) Create(ctx context.Context, h *household.Household) error {
	const q = `INSERT INTO households (id, project_id, reference_number, head_person_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	now := time.Now().UTC()
	h.CreatedAt = now
	h.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q, h.ID, h.ProjectID, h.ReferenceNumber, h.HeadPersonID, h.CreatedAt, h.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create household: %w", err)
	}
	return nil
}

func (r *householdRepo) Update(ctx context.Context, h *household.Household) error {
	const q = `UPDATE households SET reference_number=$2, head_person_id=$3, updated_at=$4 WHERE id=$1`
	h.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, h.ID, h.ReferenceNumber, h.HeadPersonID, h.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update household: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return household.ErrHouseholdNotFound
	}
	return nil
}

func (r *householdRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM households WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete household: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return household.ErrHouseholdNotFound
	}
	return nil
}
