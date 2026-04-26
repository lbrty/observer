package report

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	domainreport "github.com/lbrty/observer/internal/domain/report"
	"github.com/lbrty/observer/internal/repository"
)

type reportRepo struct {
	db *sqlx.DB
}

// New creates a ReportRepository backed by the given DB.
func New(db *sqlx.DB) repository.ReportRepository {
	return &reportRepo{db: db}
}

// peopleFilterCond returns squirrel conditions for people-based queries.
// dateCol is the column for date range filtering (e.g. "p.registered_at").
func peopleFilterCond(f domainreport.ReportFilter, dateCol string) sq.And {
	cond := sq.And{}
	if f.DateFrom != nil {
		cond = append(cond, sq.GtOrEq{dateCol: *f.DateFrom})
	}
	if f.DateTo != nil {
		cond = append(cond, sq.LtOrEq{dateCol: *f.DateTo})
	}
	if f.CaseStatus != nil {
		cond = append(cond, sq.Eq{"p.case_status": *f.CaseStatus})
	}
	if f.Sex != nil {
		cond = append(cond, sq.Eq{"p.sex": *f.Sex})
	}
	if f.CategoryID != nil {
		cond = append(cond, sq.Expr("p.id IN (SELECT person_id FROM person_categories WHERE category_id = ?)", *f.CategoryID))
	}
	if f.OfficeID != nil {
		cond = append(cond, sq.Eq{"p.office_id": *f.OfficeID})
	}
	if f.ConsultantID != nil {
		cond = append(cond, sq.Eq{"p.consultant_id": *f.ConsultantID})
	}
	if f.SupportType != nil {
		cond = append(cond, sq.Expr("p.id IN (SELECT sr2.person_id FROM support_records sr2 WHERE sr2.project_id = ? AND sr2.type = ?)", f.ProjectID, *f.SupportType))
	}
	return cond
}

// supportFilterCond returns squirrel conditions for support_records-based queries.
func supportFilterCond(f domainreport.ReportFilter, dateCol string) sq.And {
	cond := sq.And{}
	if f.DateFrom != nil {
		cond = append(cond, sq.GtOrEq{dateCol: *f.DateFrom})
	}
	if f.DateTo != nil {
		cond = append(cond, sq.LtOrEq{dateCol: *f.DateTo})
	}
	if f.OfficeID != nil {
		cond = append(cond, sq.Eq{"sr.office_id": *f.OfficeID})
	}
	if f.ConsultantID != nil {
		cond = append(cond, sq.Eq{"sr.consultant_id": *f.ConsultantID})
	}
	if f.SupportType != nil {
		cond = append(cond, sq.Eq{"sr.type": *f.SupportType})
	}
	// Person-level subquery filter
	var personCond sq.And
	if f.CategoryID != nil {
		personCond = append(personCond, sq.Expr("id IN (SELECT person_id FROM person_categories WHERE category_id = ?)", *f.CategoryID))
	}
	if f.Sex != nil {
		personCond = append(personCond, sq.Eq{"sex": *f.Sex})
	}
	if f.CaseStatus != nil {
		personCond = append(personCond, sq.Eq{"case_status": *f.CaseStatus})
	}
	if len(personCond) > 0 {
		subSQL, subArgs, _ := repository.PSQL.Select("id").From("people").Where(personCond).ToSql()
		cond = append(cond, sq.Expr("sr.person_id IN ("+subSQL+")", subArgs...))
	}
	return cond
}

func (r *reportRepo) CountConsultations(ctx context.Context, f domainreport.ReportFilter) ([]domainreport.CountResult, error) {
	sqlStr, args, err := repository.PSQL.
		Select("sr.type AS label, COUNT(*) AS count").
		From("support_records sr").
		Where(sq.Eq{"sr.project_id": f.ProjectID}).
		Where(supportFilterCond(f, "sr.provided_at")).
		GroupBy("sr.type").
		OrderBy("count DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var results []domainreport.CountResult
	if err := r.db.SelectContext(ctx, &results, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count consultations: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountBySex(ctx context.Context, f domainreport.ReportFilter) ([]domainreport.CountResult, error) {
	sqlStr, args, err := repository.PSQL.
		Select("COALESCE(p.sex, 'unknown') AS label, COUNT(*) AS count").
		From("people p").
		Where(sq.Eq{"p.project_id": f.ProjectID}).
		Where(peopleFilterCond(f, "p.registered_at")).
		GroupBy("label").
		OrderBy("count DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var results []domainreport.CountResult
	if err := r.db.SelectContext(ctx, &results, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count by sex: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByIDPStatus(ctx context.Context, f domainreport.ReportFilter) ([]domainreport.CountResult, error) {
	sqlStr, args, err := repository.PSQL.
		Select("COALESCE(s.conflict_zone, 'unknown') AS label, COUNT(DISTINCT p.id) AS count").
		From("people p").
		LeftJoin("places pl ON p.origin_place_id = pl.id").
		LeftJoin("states s ON pl.state_id = s.id").
		Where(sq.Eq{"p.project_id": f.ProjectID}).
		Where(peopleFilterCond(f, "p.registered_at")).
		GroupBy("label").
		OrderBy("count DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var results []domainreport.CountResult
	if err := r.db.SelectContext(ctx, &results, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count by idp status: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByCategory(ctx context.Context, f domainreport.ReportFilter) ([]domainreport.CountResult, error) {
	sqlStr, args, err := repository.PSQL.
		Select("c.name AS label, COUNT(DISTINCT pc.person_id) AS count").
		From("person_categories pc").
		Join("categories c ON pc.category_id = c.id").
		Join("people p ON pc.person_id = p.id").
		Where(sq.Eq{"p.project_id": f.ProjectID}).
		Where(peopleFilterCond(f, "p.registered_at")).
		GroupBy("c.name").
		OrderBy("count DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var results []domainreport.CountResult
	if err := r.db.SelectContext(ctx, &results, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count by category: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByCurrentRegion(ctx context.Context, f domainreport.ReportFilter) ([]domainreport.CountResult, error) {
	sqlStr, args, err := repository.PSQL.
		Select("COALESCE(s.name, 'unknown') AS label, COUNT(DISTINCT p.id) AS count").
		From("people p").
		LeftJoin("places pl ON p.current_place_id = pl.id").
		LeftJoin("states s ON pl.state_id = s.id").
		Where(sq.Eq{"p.project_id": f.ProjectID}).
		Where(peopleFilterCond(f, "p.registered_at")).
		GroupBy("label").
		OrderBy("count DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var results []domainreport.CountResult
	if err := r.db.SelectContext(ctx, &results, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count by region: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountBySphere(ctx context.Context, f domainreport.ReportFilter) ([]domainreport.CountResult, error) {
	sqlStr, args, err := repository.PSQL.
		Select("COALESCE(sr.sphere::text, 'unspecified') AS label, COUNT(*) AS count").
		From("support_records sr").
		Where(sq.Eq{"sr.project_id": f.ProjectID}).
		Where(supportFilterCond(f, "sr.provided_at")).
		GroupBy("sr.sphere").
		OrderBy("count DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var results []domainreport.CountResult
	if err := r.db.SelectContext(ctx, &results, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count by sphere: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountPeopleBySphere(ctx context.Context, f domainreport.ReportFilter) ([]domainreport.CountResult, error) {
	sqlStr, args, err := repository.PSQL.
		Select("COALESCE(sr.sphere::text, 'unspecified') AS label, COUNT(DISTINCT sr.person_id) AS count").
		From("support_records sr").
		Where(sq.Eq{"sr.project_id": f.ProjectID}).
		Where(supportFilterCond(f, "sr.provided_at")).
		GroupBy("sr.sphere").
		OrderBy("count DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var results []domainreport.CountResult
	if err := r.db.SelectContext(ctx, &results, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count people by sphere: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByOffice(ctx context.Context, f domainreport.ReportFilter) ([]domainreport.CountResult, error) {
	sqlStr, args, err := repository.PSQL.
		Select("COALESCE(o.name, 'unassigned') AS label, COUNT(*) AS count").
		From("support_records sr").
		LeftJoin("offices o ON sr.office_id = o.id").
		Where(sq.Eq{"sr.project_id": f.ProjectID}).
		Where(supportFilterCond(f, "sr.provided_at")).
		GroupBy("label").
		OrderBy("count DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var results []domainreport.CountResult
	if err := r.db.SelectContext(ctx, &results, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count by office: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByAgeGroup(ctx context.Context, f domainreport.ReportFilter) ([]domainreport.CountResult, error) {
	sqlStr, args, err := repository.PSQL.
		Select(`COALESCE(p.age_group::text,
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
		) AS label, COUNT(DISTINCT p.id) AS count`).
		From("people p").
		Where(sq.Eq{"p.project_id": f.ProjectID}).
		Where(peopleFilterCond(f, "p.registered_at")).
		GroupBy("label").
		OrderBy("label").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var results []domainreport.CountResult
	if err := r.db.SelectContext(ctx, &results, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count by age group: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountConsultationsByAgeGroup(ctx context.Context, f domainreport.ReportFilter) ([]domainreport.CountResult, error) {
	sqlStr, args, err := repository.PSQL.
		Select(`COALESCE(p.age_group::text,
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
		) AS label, COUNT(*) AS count`).
		From("support_records sr").
		Join("people p ON sr.person_id = p.id").
		Where(sq.Eq{"sr.project_id": f.ProjectID}).
		Where(supportFilterCond(f, "sr.provided_at")).
		GroupBy("label").
		OrderBy("label").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var results []domainreport.CountResult
	if err := r.db.SelectContext(ctx, &results, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count consultations by age group: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByTag(ctx context.Context, f domainreport.ReportFilter) ([]domainreport.CountResult, error) {
	sqlStr, args, err := repository.PSQL.
		Select("t.name AS label, COUNT(DISTINCT pt.person_id) AS count").
		From("person_tags pt").
		Join("tags t ON pt.tag_id = t.id").
		Join("people p ON pt.person_id = p.id").
		Where(sq.Eq{"p.project_id": f.ProjectID}).
		Where(peopleFilterCond(f, "p.registered_at")).
		GroupBy("t.name").
		OrderBy("count DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var results []domainreport.CountResult
	if err := r.db.SelectContext(ctx, &results, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count by tag: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountFamilyUnits(ctx context.Context, f domainreport.ReportFilter) ([]domainreport.CountResult, error) {
	cond := sq.And{sq.Eq{"h.project_id": f.ProjectID}}
	if f.SupportType != nil {
		cond = append(cond, sq.Expr(`h.id IN (
			SELECT hm.household_id FROM household_members hm
			JOIN support_records sr ON sr.person_id = hm.person_id
			WHERE sr.project_id = ? AND sr.type = ?)`, f.ProjectID, *f.SupportType))
	}

	sqlStr, args, err := repository.PSQL.
		Select("'households' AS label, COUNT(DISTINCT h.id) AS count").
		From("households h").
		Where(cond).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var results []domainreport.CountResult
	if err := r.db.SelectContext(ctx, &results, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count family units: %w", err)
	}
	return results, nil
}

func (r *reportRepo) CountByCaseStatus(ctx context.Context, f domainreport.ReportFilter) ([]domainreport.CountResult, error) {
	sqlStr, args, err := repository.PSQL.
		Select("p.case_status AS label, COUNT(*) AS count").
		From("people p").
		Where(sq.Eq{"p.project_id": f.ProjectID}).
		Where(peopleFilterCond(f, "p.registered_at")).
		GroupBy("p.case_status").
		OrderBy("count DESC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var results []domainreport.CountResult
	if err := r.db.SelectContext(ctx, &results, sqlStr, args...); err != nil {
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
		return " LEFT JOIN pets ON pets.owner_id = p.id"
	case "pet_tag":
		if metric == "pets" {
			return " LEFT JOIN pet_tags ptg ON ptg.pet_id = pets.id LEFT JOIN tags pt_tag ON pt_tag.id = ptg.tag_id"
		}
		return " LEFT JOIN pets ON pets.owner_id = p.id LEFT JOIN pet_tags ptg ON ptg.pet_id = pets.id LEFT JOIN tags pt_tag ON pt_tag.id = ptg.tag_id"
	default:
		return ""
	}
}

func (r *reportRepo) CustomQuery(ctx context.Context, projectID string, metric string, groupBy []string, filter domainreport.ReportFilter) ([]domainreport.CustomResult, int, error) {
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
		from = "pets JOIN people p ON p.id = pets.owner_id"
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

	qb := repository.PSQL.
		Select(strings.Join(selectCols, ", ") + ", " + countExpr + " AS count").
		From(from + joins.String()).
		Where(sq.Eq{projectCol: projectID})

	if filter.DateFrom != nil {
		qb = qb.Where(sq.GtOrEq{dateCol: *filter.DateFrom})
	}
	if filter.DateTo != nil {
		qb = qb.Where(sq.LtOrEq{dateCol: *filter.DateTo})
	}
	if filter.SupportType != nil {
		if metric == "events" {
			qb = qb.Where(sq.Eq{"sr.type": *filter.SupportType})
		} else {
			qb = qb.Where(sq.Expr("p.id IN (SELECT sr2.person_id FROM support_records sr2 WHERE sr2.project_id = ? AND sr2.type = ?)", projectID, *filter.SupportType))
		}
	}
	if filter.Sex != nil {
		qb = qb.Where(sq.Eq{"p.sex": *filter.Sex})
	}
	if filter.OfficeID != nil {
		if metric == "events" {
			qb = qb.Where(sq.Eq{"sr.office_id": *filter.OfficeID})
		} else {
			qb = qb.Where(sq.Eq{"p.office_id": *filter.OfficeID})
		}
	}
	if filter.CategoryID != nil {
		qb = qb.Where(sq.Expr("p.id IN (SELECT person_id FROM person_categories WHERE category_id = ?)", *filter.CategoryID))
	}
	if filter.CaseStatus != nil {
		qb = qb.Where(sq.Eq{"p.case_status": *filter.CaseStatus})
	}

	qb = qb.GroupBy(strings.Join(groupCols, ", ")).OrderBy("count DESC")

	sqlStr, args, err := qb.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build custom query: %w", err)
	}

	type row struct {
		Dim0  *string `db:"dim_0"`
		Dim1  *string `db:"dim_1"`
		Count int     `db:"count"`
	}

	var rows []row
	if err := r.db.SelectContext(ctx, &rows, sqlStr, args...); err != nil {
		return nil, 0, fmt.Errorf("custom query: %w", err)
	}

	total := 0
	results := make([]domainreport.CustomResult, len(rows))
	for i, rw := range rows {
		dims := make(map[string]string, len(groupBy))
		if len(groupBy) > 0 && rw.Dim0 != nil {
			dims[groupBy[0]] = *rw.Dim0
		}
		if len(groupBy) > 1 && rw.Dim1 != nil {
			dims[groupBy[1]] = *rw.Dim1
		}
		results[i] = domainreport.CustomResult{Dimensions: dims, Count: rw.Count}
		total += rw.Count
	}

	return results, total, nil
}

func (r *reportRepo) StatusFlowReport(ctx context.Context, f domainreport.ReportFilter) ([]domainreport.StatusFlow, error) {
	sqlStr, args, err := repository.PSQL.
		Select(`from_status, to_status, COUNT(*) AS count,
			COALESCE(AVG(EXTRACT(EPOCH FROM (h.changed_at - COALESCE(prev.changed_at, p.registered_at, p.created_at))) / 86400)::numeric(10,1), 0) AS avg_days`).
		From("person_status_history h").
		Join("people p ON h.person_id = p.id").
		LeftJoin(`LATERAL (
			SELECT changed_at FROM person_status_history h2
			WHERE h2.person_id = h.person_id AND h2.changed_at < h.changed_at
			ORDER BY h2.changed_at DESC LIMIT 1
		) prev ON true`).
		Where(sq.Eq{"p.project_id": f.ProjectID}).
		Where(peopleFilterCond(f, "h.changed_at")).
		GroupBy("from_status, to_status").
		OrderBy("count DESC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var results []domainreport.StatusFlow
	if err := r.db.SelectContext(ctx, &results, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("status flow report: %w", err)
	}

	return results, nil
}
