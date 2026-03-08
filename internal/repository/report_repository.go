package repository

import (
	"context"
	"fmt"
	"strings"

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

func (r *reportRepo) CountPeopleBySphere(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
	q := `SELECT COALESCE(sr.sphere::text, 'unspecified') AS label, COUNT(DISTINCT sr.person_id) AS count
		FROM support_records sr WHERE sr.project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = applySupportFilters(q, f, "sr.provided_at", args, 2)
	q += " GROUP BY sr.sphere ORDER BY count DESC"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count people by sphere: %w", err)
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

func (r *reportRepo) CountConsultationsByAgeGroup(ctx context.Context, f report.ReportFilter) ([]report.CountResult, error) {
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
		COUNT(*) AS count
	FROM support_records sr
	JOIN people p ON sr.person_id = p.id
	WHERE sr.project_id = $1`
	args := []any{f.ProjectID}
	q, args, _ = applySupportFilters(q, f, "sr.provided_at", args, 2)
	q += " GROUP BY label ORDER BY label"

	var results []report.CountResult
	if err := r.db.SelectContext(ctx, &results, q, args...); err != nil {
		return nil, fmt.Errorf("count consultations by age group: %w", err)
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

// dimensionSQL maps allowed dimension names to SQL expressions.
var dimensionSQL = map[string]string{
	"sex": "p.sex",
	"age_group": `COALESCE(p.age_group::text,
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
		END)`,
	"region":        "COALESCE(st.name, 'unknown')",
	"conflict_zone": "COALESCE(st.conflict_zone, 'none')",
	"office":        "COALESCE(o.name, 'unknown')",
	"sphere":        "COALESCE(sr.sphere::text, 'unknown')",
	"category":      "COALESCE(cat.name, 'uncategorised')",
	"person_tag":    "COALESCE(t.name, 'untagged')",
	"pet_status":    "COALESCE(pets.status::text, 'unknown')",
	"pet_tag":       "COALESCE(pt_tag.name, 'untagged')",
}

// dimensionJoins returns SQL JOIN clauses needed for a given dimension and metric.
func dimensionJoins(dim, metric string) string {
	switch dim {
	case "region", "conflict_zone":
		return " LEFT JOIN places pl ON pl.id = p.current_place_id LEFT JOIN states st ON st.id = pl.state_id"
	case "office":
		if metric == "events" {
			return " LEFT JOIN offices o ON o.id = sr.office_id"
		}
		return " LEFT JOIN offices o ON o.id = p.office_id"
	case "category":
		return " LEFT JOIN person_categories pc ON pc.person_id = p.id LEFT JOIN categories cat ON cat.id = pc.category_id"
	case "person_tag":
		return " LEFT JOIN person_tags pt ON pt.person_id = p.id LEFT JOIN tags t ON t.id = pt.tag_id"
	case "pet_status":
		if metric == "pets" {
			return ""
		}
		return " LEFT JOIN pets ON pets.person_id = p.id"
	case "pet_tag":
		if metric == "pets" {
			return " LEFT JOIN pet_tags ptg ON ptg.pet_id = pets.id LEFT JOIN tags pt_tag ON pt_tag.id = ptg.tag_id"
		}
		return " LEFT JOIN pets ON pets.person_id = p.id LEFT JOIN pet_tags ptg ON ptg.pet_id = pets.id LEFT JOIN tags pt_tag ON pt_tag.id = ptg.tag_id"
	default:
		return ""
	}
}

func (r *reportRepo) CustomQuery(ctx context.Context, projectID string, metric string, groupBy []string, filter report.ReportFilter) ([]report.CustomResult, int, error) {
	// Validate dimensions
	for _, dim := range groupBy {
		if _, ok := dimensionSQL[dim]; !ok {
			return nil, 0, fmt.Errorf("unknown dimension: %s", dim)
		}
	}

	// Build base FROM clause per metric
	var from, countExpr, projectCol, dateCol string
	switch metric {
	case "events":
		from = "support_records sr JOIN people p ON p.id = sr.person_id"
		countExpr = "COUNT(sr.id)"
		projectCol = "sr.project_id"
		dateCol = "sr.provided_at"
	case "people":
		from = "people p"
		countExpr = "COUNT(DISTINCT p.id)"
		projectCol = "p.project_id"
		dateCol = "p.registered_at"
	case "units":
		from = "household_members hm JOIN households h ON h.id = hm.household_id JOIN people p ON p.id = hm.person_id"
		countExpr = "COUNT(DISTINCT hm.household_id)"
		projectCol = "h.project_id"
		dateCol = "p.registered_at"
	case "pets":
		from = "pets JOIN people p ON p.id = pets.person_id"
		countExpr = "COUNT(pets.id)"
		projectCol = "p.project_id"
		dateCol = "pets.created_at"
	default:
		return nil, 0, fmt.Errorf("unknown metric: %s", metric)
	}

	// Collect JOINs (deduplicate by tracking added tables)
	added := map[string]bool{}
	var joins strings.Builder
	for _, dim := range groupBy {
		j := dimensionJoins(dim, metric)
		if j != "" && !added[j] {
			joins.WriteString(j)
			added[j] = true
		}
	}

	// Need support_records join for "sphere" dimension when metric is "people"
	if metric == "people" {
		for _, dim := range groupBy {
			if dim == "sphere" {
				srJoin := " LEFT JOIN support_records sr ON sr.person_id = p.id AND sr.project_id = p.project_id"
				if !added[srJoin] {
					joins.WriteString(srJoin)
					added[srJoin] = true
				}
				break
			}
		}
	}

	// Build SELECT columns and GROUP BY
	selectCols := make([]string, len(groupBy))
	groupCols := make([]string, len(groupBy))
	for i, dim := range groupBy {
		alias := fmt.Sprintf("dim_%d", i)
		selectCols[i] = fmt.Sprintf("%s AS %s", dimensionSQL[dim], alias)
		groupCols[i] = alias
	}

	q := fmt.Sprintf("SELECT %s, %s AS count FROM %s%s WHERE %s = $1",
		strings.Join(selectCols, ", "),
		countExpr,
		from,
		joins.String(),
		projectCol,
	)

	args := []any{projectID}
	ix := 2

	// Apply date filters
	if filter.DateFrom != nil {
		q += fmt.Sprintf(" AND %s >= $%d", dateCol, ix)
		args = append(args, *filter.DateFrom)
		ix++
	}
	if filter.DateTo != nil {
		q += fmt.Sprintf(" AND %s <= $%d", dateCol, ix)
		args = append(args, *filter.DateTo)
		ix++
	}
	if filter.SupportType != nil {
		if metric == "events" {
			q += fmt.Sprintf(" AND sr.type = $%d", ix)
		} else {
			q += fmt.Sprintf(" AND p.id IN (SELECT sr2.person_id FROM support_records sr2 WHERE sr2.project_id = $1 AND sr2.type = $%d)", ix)
		}
		args = append(args, *filter.SupportType)
		ix++
	}
	if filter.Sex != nil {
		q += fmt.Sprintf(" AND p.sex = $%d", ix)
		args = append(args, *filter.Sex)
		ix++
	}
	if filter.OfficeID != nil {
		if metric == "events" {
			q += fmt.Sprintf(" AND sr.office_id = $%d", ix)
		} else {
			q += fmt.Sprintf(" AND p.office_id = $%d", ix)
		}
		args = append(args, *filter.OfficeID)
		ix++
	}
	if filter.CategoryID != nil {
		q += fmt.Sprintf(" AND p.id IN (SELECT person_id FROM person_categories WHERE category_id = $%d)", ix)
		args = append(args, *filter.CategoryID)
		ix++
	}
	if filter.CaseStatus != nil {
		q += fmt.Sprintf(" AND p.case_status = $%d", ix)
		args = append(args, *filter.CaseStatus)
		ix++
	}
	_ = ix

	q += " GROUP BY " + strings.Join(groupCols, ", ")
	q += " ORDER BY count DESC"

	type row struct {
		Dim0  *string `db:"dim_0"`
		Dim1  *string `db:"dim_1"`
		Count int     `db:"count"`
	}

	var rows []row
	if err := r.db.SelectContext(ctx, &rows, q, args...); err != nil {
		return nil, 0, fmt.Errorf("custom query: %w", err)
	}

	total := 0
	results := make([]report.CustomResult, len(rows))
	for i, rw := range rows {
		dims := make(map[string]string, len(groupBy))
		if len(groupBy) > 0 && rw.Dim0 != nil {
			dims[groupBy[0]] = *rw.Dim0
		}
		if len(groupBy) > 1 && rw.Dim1 != nil {
			dims[groupBy[1]] = *rw.Dim1
		}
		results[i] = report.CustomResult{Dimensions: dims, Count: rw.Count}
		total += rw.Count
	}

	return results, total, nil
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
