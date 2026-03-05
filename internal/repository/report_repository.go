package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/domain/report"
)

type reportRepo struct {
	db *sqlx.DB
}

// NewReportRepository creates a ReportRepository backed by the given DB.
func NewReportRepository(db *sqlx.DB) ReportRepository {
	return &reportRepo{db: db}
}

// dateFilter appends date range conditions to a query and returns the updated query and args.
func dateFilter(q string, f report.ReportFilter, dateCol string, args []any, ix int) (string, []any, int) {
	if f.DateFrom != nil {
		q += fmt.Sprintf(" AND %s >= $%d", dateCol, ix)
		args = append(args, *f.DateFrom)
		ix++
	}
	if f.DateTo != nil {
		q += fmt.Sprintf(" AND %s <= $%d", dateCol, ix)
		args = append(args, *f.DateTo)
		ix++
	}
	return q, args, ix
}

func (r *reportRepo) CountConsultations(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT sr.type AS label, COUNT(*) AS count
		FROM support_records sr WHERE sr.project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = dateFilter(q, f, "sr.provided_at", args, 2)
	q += " GROUP BY sr.type ORDER BY count DESC"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count consultations: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountBySex(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT COALESCE(p.sex, 'unknown') AS label, COUNT(*) AS count
		FROM people p WHERE p.project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = dateFilter(q, f, "p.registered_at", args, 2)
	q += " GROUP BY label ORDER BY count DESC"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count by sex: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByIDPStatus(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT COALESCE(s.conflict_zone, 'unknown') AS label, COUNT(DISTINCT p.id) AS count
		FROM people p
		LEFT JOIN places pl ON p.origin_place_id = pl.id
		LEFT JOIN states s ON pl.state_id = s.id
		WHERE p.project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = dateFilter(q, f, "p.registered_at", args, 2)
	q += " GROUP BY label ORDER BY count DESC"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count by idp status: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByCategory(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT c.name AS label, COUNT(DISTINCT pc.person_id) AS count
		FROM person_categories pc
		JOIN categories c ON pc.category_id = c.id
		JOIN people p ON pc.person_id = p.id
		WHERE p.project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = dateFilter(q, f, "p.registered_at", args, 2)
	q += " GROUP BY c.name ORDER BY count DESC"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count by category: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByCurrentRegion(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT COALESCE(s.name, 'unknown') AS label, COUNT(DISTINCT p.id) AS count
		FROM people p
		LEFT JOIN places pl ON p.current_place_id = pl.id
		LEFT JOIN states s ON pl.state_id = s.id
		WHERE p.project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = dateFilter(q, f, "p.registered_at", args, 2)
	q += " GROUP BY label ORDER BY count DESC"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count by region: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountBySphere(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT sr.sphere AS label, COUNT(*) AS count
		FROM support_records sr WHERE sr.project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = dateFilter(q, f, "sr.provided_at", args, 2)
	q += " GROUP BY sr.sphere ORDER BY count DESC"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count by sphere: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByOffice(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT COALESCE(o.name, 'unassigned') AS label, COUNT(*) AS count
		FROM support_records sr
		LEFT JOIN offices o ON sr.office_id = o.id
		WHERE sr.project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = dateFilter(q, f, "sr.provided_at", args, 2)
	q += " GROUP BY label ORDER BY count DESC"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count by office: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByAgeGroup(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT
		CASE
			WHEN age < 1 THEN 'infant'
			WHEN age < 3 THEN 'toddler'
			WHEN age < 6 THEN 'pre_school'
			WHEN age < 12 THEN 'middle_childhood'
			WHEN age < 15 THEN 'young_teen'
			WHEN age < 18 THEN 'teenager'
			WHEN age < 25 THEN 'young_adult'
			WHEN age < 45 THEN 'early_adult'
			WHEN age < 65 THEN 'middle_aged'
			ELSE 'older_adult'
		END AS label,
		COUNT(*) AS count
	FROM (
		SELECT EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.birth_date)) AS age
		FROM people p
		WHERE p.project_id = $1 AND p.birth_date IS NOT NULL
	) sub
	GROUP BY label ORDER BY label`
	args := []any{f.ProjectID}

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count by age group: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByTag(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT t.name AS label, COUNT(DISTINCT pt.person_id) AS count
		FROM person_tags pt
		JOIN tags t ON pt.tag_id = t.id
		JOIN people p ON pt.person_id = p.id
		WHERE p.project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = dateFilter(q, f, "p.registered_at", args, 2)
	q += " GROUP BY t.name ORDER BY count DESC"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count by tag: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountFamilyUnits(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT 'households' AS label, COUNT(DISTINCT h.id) AS count
		FROM households h WHERE h.project_id = $1`
	args := []any{f.ProjectID}

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count family units: %w", err)
	}
	return results, nil
}
