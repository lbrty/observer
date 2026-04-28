package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	domainaudit "github.com/lbrty/observer/internal/domain/audit"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

type auditRow struct {
	ID            string          `db:"id"`
	ProjectID     *string         `db:"project_id"`
	UserID        *string         `db:"user_id"`
	Action        string          `db:"action"`
	EntityType    string          `db:"entity_type"`
	EntityID      *string         `db:"entity_id"`
	Summary       string          `db:"summary"`
	Details       json.RawMessage `db:"details"`
	IP            *string         `db:"ip"`
	UserAgent     *string         `db:"user_agent"`
	CreatedAt     time.Time       `db:"created_at"`
	UserFirstName *string         `db:"user_first_name"`
	UserLastName  *string         `db:"user_last_name"`
	UserEmail     *string         `db:"user_email"`
}

func scanAuditRow(r auditRow) domainaudit.Entry {
	entry := domainaudit.Entry{
		ID:            r.ID,
		ProjectID:     r.ProjectID,
		UserID:        r.UserID,
		Action:        r.Action,
		EntityType:    r.EntityType,
		EntityID:      r.EntityID,
		Summary:       r.Summary,
		IP:            r.IP,
		UserAgent:     r.UserAgent,
		CreatedAt:     r.CreatedAt,
		UserFirstName: r.UserFirstName,
		UserLastName:  r.UserLastName,
		UserEmail:     r.UserEmail,
	}
	if r.Details != nil {
		_ = json.Unmarshal(r.Details, &entry.Details)
	}
	return entry
}

type auditLogRepo struct {
	db *sqlx.DB
}

// New creates an AuditLogRepository.
func New(db *sqlx.DB) repository.AuditLogRepository {
	return &auditLogRepo{db: db}
}

func (r *auditLogRepo) Log(ctx context.Context, entry domainaudit.Entry) error {
	entry.ID = ulid.NewString()
	var detailsJSON []byte
	if entry.Details != nil {
		b, err := json.Marshal(entry.Details)
		if err != nil {
			slog.Warn("audit: failed to marshal details", slog.Any("err", err))
		} else {
			detailsJSON = b
		}
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO audit_logs (id, project_id, user_id, action, entity_type, entity_id, summary, ip, user_agent, details)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		entry.ID, entry.ProjectID, entry.UserID, entry.Action, entry.EntityType, entry.EntityID,
		entry.Summary, entry.IP, entry.UserAgent, detailsJSON,
	)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	return nil
}

func (r *auditLogRepo) List(ctx context.Context, filter domainaudit.Filter) ([]domainaudit.Entry, int, error) {
	cond := sq.And{}

	if filter.ProjectID != nil {
		cond = append(cond, sq.Eq{"a.project_id": *filter.ProjectID})
	}

	if filter.UserID != nil {
		cond = append(cond, sq.Eq{"a.user_id": *filter.UserID})
	}

	if filter.Action != nil {
		cond = append(cond, sq.Eq{"a.action": *filter.Action})
	}

	if filter.EntityType != nil {
		cond = append(cond, sq.Eq{"a.entity_type": *filter.EntityType})
	}

	if filter.DateFrom != nil {
		cond = append(cond, sq.GtOrEq{"a.created_at": *filter.DateFrom})
	}

	if filter.DateTo != nil {
		cond = append(cond, sq.LtOrEq{"a.created_at": *filter.DateTo})
	}

	countSQL, countArgs, err := repository.PSQL.Select("COUNT(*)").From("audit_logs a").Where(cond).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build count query: %w", err)
	}

	var total int
	if err := r.db.GetContext(ctx, &total, countSQL, countArgs...); err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}

	offset := (filter.Page - 1) * filter.PerPage

	listSQL, listArgs, err := repository.PSQL.
		Select(
			"a.id", "a.project_id", "a.user_id", "a.action", "a.entity_type", "a.entity_id",
			"a.summary", "a.details", "a.ip", "a.user_agent", "a.created_at",
			"u.first_name AS user_first_name", "u.last_name AS user_last_name", "u.email AS user_email",
		).
		From("audit_logs a").
		LeftJoin("users u ON u.id = a.user_id").
		Where(cond).
		OrderBy("a.created_at DESC").
		Limit(uint64(filter.PerPage)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		return nil, 0, fmt.Errorf("build list query: %w", err)
	}

	var rows []auditRow
	if err := r.db.SelectContext(ctx, &rows, listSQL, listArgs...); err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}

	entries := make([]domainaudit.Entry, len(rows))
	for i, row := range rows {
		entries[i] = scanAuditRow(row)
	}

	return entries, total, nil
}
