package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
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

func scanHousehold(row interface{ Scan(dest ...any) error }) (*household.Household, error) {
	var h household.Household
	if err := row.Scan(&h.ID, &h.ProjectID, &h.ReferenceNumber, &h.HeadPersonID, &h.CreatedAt, &h.UpdatedAt); err != nil {
		return nil, err
	}
	TimesToUTC(&h.CreatedAt, &h.UpdatedAt)
	return &h, nil
}

func scanHouseholdWithCount(row interface{ Scan(dest ...any) error }) (*household.Household, error) {
	var h household.Household
	var firstName, lastName sql.NullString
	if err := row.Scan(&h.ID, &h.ProjectID, &h.ReferenceNumber, &h.HeadPersonID, &firstName, &lastName, &h.MemberCount, &h.CreatedAt, &h.UpdatedAt); err != nil {
		return nil, err
	}
	if firstName.Valid {
		h.HeadFirstName = &firstName.String
	}

	if lastName.Valid {
		h.HeadLastName = &lastName.String
	}

	TimesToUTC(&h.CreatedAt, &h.UpdatedAt)
	return &h, nil
}

func scanMember(row interface{ Scan(dest ...any) error }) (*household.Member, error) {
	var m household.Member
	if err := row.Scan(&m.HouseholdID, &m.PersonID, &m.Relationship); err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *householdRepo) List(ctx context.Context, filter household.HouseholdListFilter) ([]*household.Household, int, error) {
	perPage, off := normalizePagination(filter.Page, filter.PerPage)

	cond := sq.And{sq.Eq{"h.project_id": filter.ProjectID}}
	if filter.Search != nil && *filter.Search != "" {
		cond = append(cond, sq.ILike{"h.reference_number": "%" + *filter.Search + "%"})
	}

	if filter.CreatedFrom != nil {
		cond = append(cond, sq.GtOrEq{"h.created_at": *filter.CreatedFrom})
	}

	if filter.CreatedTo != nil {
		cond = append(cond, sq.LtOrEq{"h.created_at": *filter.CreatedTo})
	}

	countSQL, countArgs, err := psql.Select("COUNT(*)").From("households h").Where(cond).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build count households query: %w", err)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count households: %w", err)
	}

	offset := uint64(off)
	listSQL, listArgs, err := psql.
		Select(
			"h.id", "h.project_id", "h.reference_number", "h.head_person_id",
			"p.first_name", "p.last_name",
			"COUNT(hm.person_id) AS member_count",
			"h.created_at", "h.updated_at",
		).
		From("households h").
		LeftJoin("people p ON p.id = h.head_person_id").
		LeftJoin("household_members hm ON hm.household_id = h.id").
		Where(cond).
		GroupBy("h.id", "p.first_name", "p.last_name").
		OrderBy("h.created_at DESC").
		Limit(uint64(perPage)).
		Offset(offset).
		ToSql()

	if err != nil {
		return nil, 0, fmt.Errorf("build list households query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list households: %w", err)
	}
	defer rows.Close()

	var out []*household.Household
	for rows.Next() {
		h, err := scanHouseholdWithCount(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan household: %w", err)
		}
		out = append(out, h)
	}

	return out, total, rows.Err()
}

func (r *householdRepo) GetByID(ctx context.Context, id string) (*household.Household, error) {
	const q = `SELECT id, project_id, reference_number, head_person_id, created_at, updated_at FROM households WHERE id = $1`
	h, err := scanHousehold(r.db.QueryRowContext(ctx, q, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, household.ErrHouseholdNotFound
		}
		return nil, fmt.Errorf("get household: %w", err)
	}
	return h, nil
}

func (r *householdRepo) Create(ctx context.Context, h *household.Household) error {
	const q = `
		INSERT INTO households (id, project_id, reference_number, head_person_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
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
	return CheckRowsAffected(res, household.ErrHouseholdNotFound)
}

func (r *householdRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM households WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete household: %w", err)
	}
	return CheckRowsAffected(res, household.ErrHouseholdNotFound)
}

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
		m, err := scanMember(rows)
		if err != nil {
			return nil, fmt.Errorf("scan household member: %w", err)
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *householdMemberRepo) Add(ctx context.Context, m *household.Member) error {
	const q = `INSERT INTO household_members (household_id, person_id, relationship) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, q, m.HouseholdID, m.PersonID, m.Relationship)
	if err != nil {
		if IsUniqueViolation(err) {
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
	return CheckRowsAffected(res, household.ErrMemberNotFound)
}
