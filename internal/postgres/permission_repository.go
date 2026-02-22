package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"

	"github.com/lbrty/observer/internal/domain/project"
)

// PermissionRepository is a PostgreSQL-backed project permission loader.
type PermissionRepository struct {
	db *sqlx.DB
}

// NewPermissionRepository creates a PermissionRepository backed by the given DB.
func NewPermissionRepository(db *sqlx.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) IsProjectOwner(ctx context.Context, userID ulid.ULID, projectID string) (bool, error) {
	const q = `SELECT owner_id FROM projects WHERE id = $1`
	var ownerID string
	err := r.db.QueryRowContext(ctx, q, projectID).Scan(&ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, project.ErrProjectNotFound
		}
		return false, fmt.Errorf("query project owner: %w", err)
	}
	return ownerID == userID.String(), nil
}

func (r *PermissionRepository) GetPermission(ctx context.Context, userID ulid.ULID, projectID string) (*project.Permission, error) {
	const q = `
		SELECT role, can_view_contact, can_view_personal, can_view_documents
		FROM project_permissions
		WHERE user_id = $1 AND project_id = $2
	`
	var role string
	var perm project.Permission
	err := r.db.QueryRowContext(ctx, q, userID.String(), projectID).Scan(
		&role, &perm.CanViewContact, &perm.CanViewPersonal, &perm.CanViewDocuments,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, project.ErrPermissionNotFound
		}
		return nil, fmt.Errorf("query project permission: %w", err)
	}
	perm.Role = project.ProjectRole(role)
	return &perm, nil
}
