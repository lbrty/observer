package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

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
	var (
		where []string
		args  []any
		ix    int
	)

	if filter.OwnerID != nil {
		ix++
		where = append(where, "owner_id = $"+strconv.Itoa(ix))
		args = append(args, *filter.OwnerID)
	}
	if filter.Status != nil {
		ix++
		where = append(where, "status = $"+strconv.Itoa(ix))
		args = append(args, string(*filter.Status))
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	countQ := "SELECT COUNT(*) FROM projects " + whereClause
	var total int
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
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

	ix++
	args = append(args, perPage)
	limitParam := "$" + strconv.Itoa(ix)
	ix++
	args = append(args, offset)
	offsetParam := "$" + strconv.Itoa(ix)

	q := "SELECT id, name, description, owner_id, status, created_at, updated_at FROM projects " +
		whereClause + " ORDER BY created_at DESC LIMIT " + limitParam + " OFFSET " + offsetParam

	rows, err := r.db.QueryContext(ctx, q, args...)
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
		p.CreatedAt = p.CreatedAt.UTC()
		p.UpdatedAt = p.UpdatedAt.UTC()
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
	p.CreatedAt = p.CreatedAt.UTC()
	p.UpdatedAt = p.UpdatedAt.UTC()
	return &p, nil
}

func (r *projectRepo) Create(ctx context.Context, p *project.Project) error {
	const q = `INSERT INTO projects (id, name, description, owner_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q, p.ID, p.Name, p.Description, p.OwnerID, p.Status, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
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
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return project.ErrProjectNameExists
		}
		return fmt.Errorf("update project: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return project.ErrProjectNotFound
	}
	return nil
}
