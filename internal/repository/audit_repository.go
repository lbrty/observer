package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/lbrty/observer/internal/domain/audit"
	"github.com/lbrty/observer/internal/ulid"
)

type auditLogRepo struct {
	db *sqlx.DB
}

// NewAuditLogRepository creates an AuditLogRepository.
func NewAuditLogRepository(db *sqlx.DB) AuditLogRepository {
	return &auditLogRepo{db: db}
}

func (r *auditLogRepo) Log(ctx context.Context, entry audit.Entry) error {
	entry.ID = ulid.NewString()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO audit_logs (id, project_id, user_id, action, entity_type, entity_id, summary, ip, user_agent)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		entry.ID, entry.ProjectID, entry.UserID, entry.Action, entry.EntityType, entry.EntityID,
		entry.Summary, entry.IP, entry.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	return nil
}

func (r *auditLogRepo) List(ctx context.Context, filter audit.Filter) ([]audit.Entry, int, error) {
	var (
		args    []any
		filters string
		ix      = 1
	)

	filters, args, ix = appendIf(filters, args, ix, " AND a.project_id = $%d", filter.ProjectID)
	filters, args, ix = appendIf(filters, args, ix, " AND a.user_id = $%d", filter.UserID)
	filters, args, ix = appendIf(filters, args, ix, " AND a.action = $%d", filter.Action)
	filters, args, ix = appendIf(filters, args, ix, " AND a.entity_type = $%d", filter.EntityType)
	filters, args, ix = appendIf(filters, args, ix, " AND a.created_at >= $%d", filter.DateFrom)
	filters, args, ix = appendIf(filters, args, ix, " AND a.created_at <= $%d", filter.DateTo)

	var total int
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM audit_logs a WHERE 1=1`+filters, args...); err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}

	offset := (filter.Page - 1) * filter.PerPage
	args = append(args, filter.PerPage, offset)

	q := `SELECT a.id, a.project_id, a.user_id, a.action, a.entity_type, a.entity_id, a.summary, a.ip, a.user_agent, a.created_at,
	             COALESCE(u.first_name, '') AS user_first_name, COALESCE(u.last_name, '') AS user_last_name, COALESCE(u.email, '') AS user_email
	      FROM audit_logs a LEFT JOIN users u ON u.id = a.user_id WHERE 1=1` +
		filters + fmt.Sprintf(" ORDER BY a.created_at DESC LIMIT $%d OFFSET $%d", ix, ix+1)

	var entries []audit.Entry
	if err := r.db.SelectContext(ctx, &entries, q, args...); err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}
	return entries, total, nil
}
