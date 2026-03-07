package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/domain/report"
)

type petReportRepo struct {
	db *sqlx.DB
}

// NewPetReportRepository creates a PetReportRepository.
func NewPetReportRepository(db *sqlx.DB) PetReportRepository {
	return &petReportRepo{db: db}
}

func applyPetFilters(q string, f report.PetReportFilter, args []any, ix int) (string, []any, int) {
	if f.DateFrom != nil {
		q += fmt.Sprintf(" AND created_at >= $%d", ix)
		args = append(args, *f.DateFrom)
		ix++
	}
	if f.DateTo != nil {
		q += fmt.Sprintf(" AND created_at <= $%d", ix)
		args = append(args, *f.DateTo)
		ix++
	}
	if f.Status != nil {
		q += fmt.Sprintf(" AND status = $%d", ix)
		args = append(args, *f.Status)
		ix++
	}
	return q, args, ix
}

func (r *petReportRepo) CountByStatus(ctx context.Context, f report.PetReportFilter) ([]report.CountResult, error) {
	q := `SELECT status AS label, COUNT(*) AS count FROM pets WHERE project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = applyPetFilters(q, f, args, 2)
	q += ` GROUP BY status ORDER BY count DESC`

	var rows []report.CountResult
	if err := r.db.SelectContext(ctx, &rows, q, args...); err != nil {
		return nil, fmt.Errorf("count pets by status: %w", err)
	}
	return rows, nil
}

func (r *petReportRepo) CountByOwnership(ctx context.Context, f report.PetReportFilter) ([]report.CountResult, error) {
	q := `SELECT CASE WHEN owner_id IS NOT NULL THEN 'with_owner' ELSE 'without_owner' END AS label,
		COUNT(*) AS count FROM pets WHERE project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = applyPetFilters(q, f, args, 2)
	q += ` GROUP BY label ORDER BY count DESC`

	var rows []report.CountResult
	if err := r.db.SelectContext(ctx, &rows, q, args...); err != nil {
		return nil, fmt.Errorf("count pets by ownership: %w", err)
	}
	return rows, nil
}

func (r *petReportRepo) CountByMonth(ctx context.Context, f report.PetReportFilter) ([]report.CountResult, error) {
	q := `SELECT to_char(created_at, 'YYYY-MM') AS label, COUNT(*) AS count
		FROM pets WHERE project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = applyPetFilters(q, f, args, 2)
	q += ` GROUP BY label ORDER BY label`

	var rows []report.CountResult
	if err := r.db.SelectContext(ctx, &rows, q, args...); err != nil {
		return nil, fmt.Errorf("count pets by month: %w", err)
	}
	return rows, nil
}

func (r *petReportRepo) CountByStatusByMonth(ctx context.Context, f report.PetReportFilter) ([]report.MonthlyStatusCount, error) {
	q := `SELECT to_char(created_at, 'YYYY-MM') AS month, status, COUNT(*) AS count
		FROM pets WHERE project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = applyPetFilters(q, f, args, 2)
	q += ` GROUP BY month, status ORDER BY month, status`

	var rows []report.MonthlyStatusCount
	if err := r.db.SelectContext(ctx, &rows, q, args...); err != nil {
		return nil, fmt.Errorf("count pets by status by month: %w", err)
	}
	return rows, nil
}
