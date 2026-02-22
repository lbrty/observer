package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/lbrty/observer/internal/domain/project"
)

// projectPermissionRepo is a PostgreSQL-backed project permission repository.
type projectPermissionRepo struct {
	db *sqlx.DB
}

// NewProjectPermissionRepository creates a ProjectPermissionRepository.
func NewProjectPermissionRepository(db *sqlx.DB) PermissionRepository {
	return &projectPermissionRepo{db: db}
}

func (r *projectPermissionRepo) List(ctx context.Context, projectID string) ([]*project.ProjectPermission, error) {
	const q = `
		SELECT id, project_id, user_id, role, can_view_contact, can_view_personal, can_view_documents, created_at, updated_at
		FROM project_permissions
		WHERE project_id = $1
		ORDER BY created_at
	`
	rows, err := r.db.QueryContext(ctx, q, projectID)
	if err != nil {
		return nil, fmt.Errorf("list project permissions: %w", err)
	}
	defer rows.Close()

	var perms []*project.ProjectPermission
	for rows.Next() {
		p, err := r.scanPermissionRow(rows)
		if err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, rows.Err()
}

func (r *projectPermissionRepo) GetByID(ctx context.Context, id string) (*project.ProjectPermission, error) {
	const q = `
		SELECT id, project_id, user_id, role, can_view_contact, can_view_personal, can_view_documents, created_at, updated_at
		FROM project_permissions
		WHERE id = $1
	`
	var p project.ProjectPermission
	var role string
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&p.ID, &p.ProjectID, &p.UserID, &role,
		&p.CanViewContact, &p.CanViewPersonal, &p.CanViewDocuments,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, project.ErrPermissionNotFound
		}
		return nil, fmt.Errorf("get project permission: %w", err)
	}
	p.Role = project.ProjectRole(role)
	p.CreatedAt = p.CreatedAt.UTC()
	p.UpdatedAt = p.UpdatedAt.UTC()
	return &p, nil
}

func (r *projectPermissionRepo) Create(ctx context.Context, p *project.ProjectPermission) error {
	const q = `
		INSERT INTO project_permissions (id, project_id, user_id, role, can_view_contact, can_view_personal, can_view_documents, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, q,
		p.ID, p.ProjectID, p.UserID, string(p.Role),
		p.CanViewContact, p.CanViewPersonal, p.CanViewDocuments,
		p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return project.ErrPermissionExists
		}
		return fmt.Errorf("create project permission: %w", err)
	}
	return nil
}

func (r *projectPermissionRepo) Update(ctx context.Context, p *project.ProjectPermission) error {
	const q = `
		UPDATE project_permissions
		SET role=$2, can_view_contact=$3, can_view_personal=$4, can_view_documents=$5, updated_at=$6
		WHERE id=$1
	`
	p.UpdatedAt = time.Now().UTC()
	res, err := r.db.ExecContext(ctx, q,
		p.ID, string(p.Role),
		p.CanViewContact, p.CanViewPersonal, p.CanViewDocuments,
		p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update project permission: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return project.ErrPermissionNotFound
	}
	return nil
}

func (r *projectPermissionRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM project_permissions WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete project permission: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return project.ErrPermissionNotFound
	}
	return nil
}

func (r *projectPermissionRepo) scanPermissionRow(rows *sql.Rows) (*project.ProjectPermission, error) {
	var p project.ProjectPermission
	var role string
	err := rows.Scan(
		&p.ID, &p.ProjectID, &p.UserID, &role,
		&p.CanViewContact, &p.CanViewPersonal, &p.CanViewDocuments,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan project permission: %w", err)
	}
	p.Role = project.ProjectRole(role)
	p.CreatedAt = p.CreatedAt.UTC()
	p.UpdatedAt = p.UpdatedAt.UTC()
	return &p, nil
}
