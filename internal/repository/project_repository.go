package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/domain/project"
)

type projectRepo struct {
	db *sqlx.DB
}

// NewProjectRepository creates a ProjectRepository.
func NewProjectRepository(db *sqlx.DB) ProjectRepository {
	return &projectRepo{db: db}
}

func (r *projectRepo) List(ctx context.Context, filter project.ProjectListFilter) ([]*project.Project, int, error) {
	cond := sq.And{}

	if filter.OwnerID != nil {
		cond = append(cond, sq.Eq{"owner_id": *filter.OwnerID})
	}

	if filter.Status != nil {
		cond = append(cond, sq.Eq{"status": string(*filter.Status)})
	}

	countSQL, countArgs, err := psql.Select("COUNT(*)").From("projects").Where(cond).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build count query: %w", err)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count projects: %w", err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}

	perPage := filter.PerPage
	if perPage < 1 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	listSQL, listArgs, err := psql.Select("id, name, description, owner_id, status, created_at, updated_at").
		From("projects").Where(cond).
		OrderBy("created_at DESC").
		Limit(uint64(perPage)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		return nil, 0, fmt.Errorf("build list query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var out []*project.Project
	for rows.Next() {
		var p project.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan project: %w", err)
		}
		TimesToUTC(&p.CreatedAt, &p.UpdatedAt)
		out = append(out, &p)
	}

	return out, total, rows.Err()
}

func (r *projectRepo) GetByID(ctx context.Context, id string) (*project.Project, error) {
	const q = `SELECT id, name, description, owner_id, status, created_at, updated_at FROM projects WHERE id = $1`
	var p project.Project
	err := r.db.QueryRowContext(ctx, q, id).Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, project.ErrProjectNotFound
		}
		return nil, fmt.Errorf("get project: %w", err)
	}

	TimesToUTC(&p.CreatedAt, &p.UpdatedAt)
	return &p, nil
}

func (r *projectRepo) Create(ctx context.Context, p *project.Project) error {
	const q = `INSERT INTO projects (id, name, description, owner_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q, p.ID, p.Name, p.Description, p.OwnerID, p.Status, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		if IsUniqueViolation(err) {
			return project.ErrProjectNameExists
		}
		return fmt.Errorf("create project: %w", err)
	}

	return nil
}

func (r *projectRepo) Update(ctx context.Context, p *project.Project) error {
	const q = `UPDATE projects SET name=$2, description=$3, status=$4, updated_at=$5 WHERE id=$1`
	p.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q, p.ID, p.Name, p.Description, p.Status, p.UpdatedAt)
	if err != nil {
		if IsUniqueViolation(err) {
			return project.ErrProjectNameExists
		}
		return fmt.Errorf("update project: %w", err)
	}

	return CheckRowsAffected(res, project.ErrProjectNotFound)
}
