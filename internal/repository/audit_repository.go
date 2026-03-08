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
	q := `SELECT id, project_id, user_id, action, entity_type, entity_id, summary, ip, user_agent, created_at
	      FROM audit_logs WHERE 1=1`
	countQ := `SELECT COUNT(*) FROM audit_logs WHERE 1=1`
	args := []any{}
	ix := 1

	if filter.ProjectID != nil {
		clause := fmt.Sprintf(" AND project_id = $%d", ix)
		q += clause
		countQ += clause
		args = append(args, *filter.ProjectID)
		ix++
	}
	if filter.UserID != nil {
		clause := fmt.Sprintf(" AND user_id = $%d", ix)
		q += clause
		countQ += clause
		args = append(args, *filter.UserID)
		ix++
	}
	if filter.Action != nil {
		clause := fmt.Sprintf(" AND action = $%d", ix)
		q += clause
		countQ += clause
		args = append(args, *filter.Action)
		ix++
	}
	if filter.EntityType != nil {
		clause := fmt.Sprintf(" AND entity_type = $%d", ix)
		q += clause
		countQ += clause
		args = append(args, *filter.EntityType)
		ix++
	}
	if filter.DateFrom != nil {
		clause := fmt.Sprintf(" AND created_at >= $%d", ix)
		q += clause
		countQ += clause
		args = append(args, *filter.DateFrom)
		ix++
	}
	if filter.DateTo != nil {
		clause := fmt.Sprintf(" AND created_at <= $%d", ix)
		q += clause
		countQ += clause
		args = append(args, *filter.DateTo)
		ix++
	}

	var total int
	if err := r.db.GetContext(ctx, &total, countQ, args...); err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}

	q += " ORDER BY created_at DESC"
	offset := (filter.Page - 1) * filter.PerPage
	q += fmt.Sprintf(" LIMIT $%d", ix)
	args = append(args, filter.PerPage)
	ix++
	q += fmt.Sprintf(" OFFSET $%d", ix)
	args = append(args, offset)

	var entries []audit.Entry
	if err := r.db.SelectContext(ctx, &entries, q, args...); err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}
	return entries, total, nil
}
