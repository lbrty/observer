package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/lbrty/observer/internal/domain/household"
)

type householdMemberRepo struct {
	db *sqlx.DB
}

// NewHouseholdMemberRepository creates a HouseholdMemberRepository.
func NewHouseholdMemberRepository(db *sqlx.DB) HouseholdMemberRepository {
	return &householdMemberRepo{db: db}
}

func (r *householdMemberRepo) List(ctx context.Context, householdID string) ([]*household.Member, error) {
	const q = `SELECT household_id, person_id, relationship FROM household_members WHERE household_id = $1`
	rows, err := r.db.QueryContext(ctx, q, householdID)
	if err != nil {
		return nil, fmt.Errorf("list household members: %w", err)
	}
	defer rows.Close()

	var out []*household.Member
	for rows.Next() {
		var m household.Member
		if err := rows.Scan(&m.HouseholdID, &m.PersonID, &m.Relationship); err != nil {
			return nil, fmt.Errorf("scan household member: %w", err)
		}
		out = append(out, &m)
	}
	return out, rows.Err()
}

func (r *householdMemberRepo) Add(ctx context.Context, m *household.Member) error {
	const q = `INSERT INTO household_members (household_id, person_id, relationship) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, q, m.HouseholdID, m.PersonID, m.Relationship)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return household.ErrMemberExists
		}
		return fmt.Errorf("add household member: %w", err)
	}
	return nil
}

func (r *householdMemberRepo) Remove(ctx context.Context, householdID, personID string) error {
	const q = `DELETE FROM household_members WHERE household_id = $1 AND person_id = $2`
	res, err := r.db.ExecContext(ctx, q, householdID, personID)
	if err != nil {
		return fmt.Errorf("remove household member: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return household.ErrMemberNotFound
	}
	return nil
}
