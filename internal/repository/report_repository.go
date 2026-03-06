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

// applyPeopleFilters appends WHERE clauses for people-based queries.
// dateCol is the column used for date range filtering (e.g. "p.registered_at").
func applyPeopleFilters(q string, f report.ReportFilter, dateCol string, args []any, ix int) (string, []any, int) {
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
	if f.CaseStatus != nil {
		q += fmt.Sprintf(" AND p.case_status = $%d", ix)
		args = append(args, *f.CaseStatus)
		ix++
	}
	if f.Sex != nil {
		q += fmt.Sprintf(" AND p.sex = $%d", ix)
		args = append(args, *f.Sex)
		ix++
	}
	if f.CategoryID != nil {
		q += fmt.Sprintf(" AND p.id IN (SELECT person_id FROM person_categories WHERE category_id = $%d)", ix)
		args = append(args, *f.CategoryID)
		ix++
	}
	if f.OfficeID != nil {
		q += fmt.Sprintf(" AND p.office_id = $%d", ix)
		args = append(args, *f.OfficeID)
		ix++
	}
	if f.ConsultantID != nil {
		q += fmt.Sprintf(" AND p.consultant_id = $%d", ix)
		args = append(args, *f.ConsultantID)
		ix++
	}
	if f.SupportType != nil {
		q += fmt.Sprintf(" AND p.id IN (SELECT sr2.person_id FROM support_records sr2 WHERE sr2.project_id = $1 AND sr2.type = $%d)", ix)
		args = append(args, *f.SupportType)
		ix++
	}
	return q, args, ix
}

// applySupportFilters appends WHERE clauses for support_records-based queries.
// dateCol is the column used for date range filtering (e.g. "sr.provided_at").
func applySupportFilters(q string, f report.ReportFilter, dateCol string, args []any, ix int) (string, []any, int) {
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
	if f.OfficeID != nil {
		q += fmt.Sprintf(" AND sr.office_id = $%d", ix)
		args = append(args, *f.OfficeID)
		ix++
	}
	if f.ConsultantID != nil {
		q += fmt.Sprintf(" AND sr.consultant_id = $%d", ix)
		args = append(args, *f.ConsultantID)
		ix++
	}
	if f.SupportType != nil {
		q += fmt.Sprintf(" AND sr.type = $%d", ix)
		args = append(args, *f.SupportType)
		ix++
	}

	// For person-related filters on support queries, use a subquery on people.
	var personClauses []string
	var personArgs []any
	if f.CategoryID != nil {
		personClauses = append(personClauses, fmt.Sprintf("id IN (SELECT person_id FROM person_categories WHERE category_id = $%d)", ix))
		personArgs = append(personArgs, *f.CategoryID)
		ix++
	}
	if f.Sex != nil {
		personClauses = append(personClauses, fmt.Sprintf("sex = $%d", ix))
		personArgs = append(personArgs, *f.Sex)
		ix++
	}
	if f.CaseStatus != nil {
		personClauses = append(personClauses, fmt.Sprintf("case_status = $%d", ix))
		personArgs = append(personArgs, *f.CaseStatus)
		ix++
	}
	if len(personClauses) > 0 {
		sub := " AND sr.person_id IN (SELECT id FROM people WHERE "
		for i, clause := range personClauses {
			if i > 0 {
				sub += " AND "
			}
			sub += clause
		}
		sub += ")"
		q += sub
		args = append(args, personArgs...)
	}

	return q, args, ix
}

func (r *reportRepo) CountConsultations(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT sr.type AS label, COUNT(*) AS count
		FROM support_records sr WHERE sr.project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = applySupportFilters(q, f, "sr.provided_at", args, 2)
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
	q, args, _ = applyPeopleFilters(q, f, "p.registered_at", args, 2)
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
	q, args, _ = applyPeopleFilters(q, f, "p.registered_at", args, 2)
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
	q, args, _ = applyPeopleFilters(q, f, "p.registered_at", args, 2)
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
	q, args, _ = applyPeopleFilters(q, f, "p.registered_at", args, 2)
	q += " GROUP BY label ORDER BY count DESC"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count by region: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountBySphere(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT COALESCE(sr.sphere::text, 'unspecified') AS label, COUNT(*) AS count
		FROM support_records sr WHERE sr.project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = applySupportFilters(q, f, "sr.provided_at", args, 2)
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
	q, args, _ = applySupportFilters(q, f, "sr.provided_at", args, 2)
	q += " GROUP BY label ORDER BY count DESC"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count by office: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByAgeGroup(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT
		COALESCE(p.age_group::text,
			CASE
				WHEN p.birth_date IS NULL THEN 'unknown'
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.birth_date)) < 1  THEN 'infant'
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.birth_date)) < 3  THEN 'toddler'
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.birth_date)) < 6  THEN 'pre_school'
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.birth_date)) < 12 THEN 'middle_childhood'
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.birth_date)) < 14 THEN 'young_teen'
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.birth_date)) < 18 THEN 'teenager'
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.birth_date)) < 25 THEN 'young_adult'
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.birth_date)) < 35 THEN 'early_adult'
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.birth_date)) < 55 THEN 'middle_aged_adult'
				ELSE 'old_adult'
			END
		) AS label,
		COUNT(DISTINCT p.id) AS count
	FROM people p
	WHERE p.project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = applyPeopleFilters(q, f, "p.registered_at", args, 2)
	q += " GROUP BY label ORDER BY label"

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
	q, args, _ = applyPeopleFilters(q, f, "p.registered_at", args, 2)
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
	ix := 2

	if f.SupportType != nil {
		q += fmt.Sprintf(` AND h.id IN (
			SELECT hm.household_id FROM household_members hm
			JOIN support_records sr ON sr.person_id = hm.person_id
			WHERE sr.project_id = $1 AND sr.type = $%d)`, ix)
		args = append(args, *f.SupportType)
		ix++
	}
	_ = ix

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count family units: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByCaseStatus(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT p.case_status AS label, COUNT(*) AS count
		FROM people p WHERE p.project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = applyPeopleFilters(q, f, "p.registered_at", args, 2)
	q += " GROUP BY p.case_status ORDER BY count DESC"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count by case status: %w", err)
	}
	return results, nil
}

func (r *reportRepo) StatusFlowReport(ctx context.Context, f report.ReportFilter) ([]report.StatusFlow, error) {
	q := `SELECT from_status, to_status, COUNT(*) AS count,
		   COALESCE(AVG(EXTRACT(EPOCH FROM (h.changed_at - COALESCE(prev.changed_at, p.registered_at, p.created_at))) / 86400)::numeric(10,1), 0) AS avg_days
	FROM person_status_history h
	JOIN people p ON h.person_id = p.id
	LEFT JOIN LATERAL (
		SELECT changed_at FROM person_status_history h2
		WHERE h2.person_id = h.person_id AND h2.changed_at < h.changed_at
		ORDER BY h2.changed_at DESC LIMIT 1
	) prev ON true
	WHERE p.project_id = $1`

	args := []any{f.ProjectID}
	q, args, _ = applyPeopleFilters(q, f, "h.changed_at", args, 2)
	q += " GROUP BY from_status, to_status ORDER BY count DESC"

	var results []report.StatusFlow
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("status flow report: %w", err)
	}
	return results, nil
}
