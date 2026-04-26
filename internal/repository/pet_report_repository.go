package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
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

func buildPetCond(f report.PetReportFilter) sq.And {
	cond := sq.And{sq.Eq{"project_id": f.ProjectID}}
	if f.DateFrom != nil {
		cond = append(cond, sq.GtOrEq{"created_at": *f.DateFrom})
	}

	if f.DateTo != nil {
		cond = append(cond, sq.LtOrEq{"created_at": *f.DateTo})
	}

	if f.Status != nil {
		cond = append(cond, sq.Eq{"status": *f.Status})
	}

	return cond
}

func (r *petReportRepo) CountByStatus(ctx context.Context, f report.PetReportFilter) ([]report.CountResult, error) {
	sqlStr, args, err := psql.
		Select("status AS label", "COUNT(*) AS count").
		From("pets").
		Where(buildPetCond(f)).
		GroupBy("status").
		OrderBy("count DESC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build count pets by status query: %w", err)
	}

	var rows []report.CountResult
	if err := r.db.SelectContext(ctx, &rows, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count pets by status: %w", err)
	}

	return rows, nil
}

func (r *petReportRepo) CountByOwnership(ctx context.Context, f report.PetReportFilter) ([]report.CountResult, error) {
	sqlStr, args, err := psql.
		Select(
			"CASE WHEN owner_id IS NOT NULL THEN 'with_owner' ELSE 'without_owner' END AS label",
			"COUNT(*) AS count",
		).
		From("pets").
		Where(buildPetCond(f)).
		GroupBy("label").
		OrderBy("count DESC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build count pets by ownership query: %w", err)
	}

	var rows []report.CountResult
	if err := r.db.SelectContext(ctx, &rows, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count pets by ownership: %w", err)
	}

	return rows, nil
}

func (r *petReportRepo) CountByMonth(ctx context.Context, f report.PetReportFilter) ([]report.CountResult, error) {
	sqlStr, args, err := psql.
		Select("TO_CHAR(created_at, 'YYYY-MM') AS label", "COUNT(*) AS count").
		From("pets").
		Where(buildPetCond(f)).
		GroupBy("label").
		OrderBy("label").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build count pets by month query: %w", err)
	}

	var rows []report.CountResult
	if err := r.db.SelectContext(ctx, &rows, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count pets by month: %w", err)
	}

	return rows, nil
}

func (r *petReportRepo) CountByStatusByMonth(ctx context.Context, f report.PetReportFilter) ([]report.MonthlyStatusCount, error) {
	sqlStr, args, err := psql.
		Select("TO_CHAR(created_at, 'YYYY-MM') AS month", "status", "COUNT(*) AS count").
		From("pets").
		Where(buildPetCond(f)).
		GroupBy("month", "status").
		OrderBy("month", "status").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build count pets by status by month query: %w", err)
	}

	var rows []report.MonthlyStatusCount
	if err := r.db.SelectContext(ctx, &rows, sqlStr, args...); err != nil {
		return nil, fmt.Errorf("count pets by status by month: %w", err)
	}

	return rows, nil
}
