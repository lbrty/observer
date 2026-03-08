package audit

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	domainaudit "github.com/lbrty/observer/internal/domain/audit"
	"github.com/lbrty/observer/internal/middleware"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/usecase"
)

type AuditUseCase struct {
	repo repository.AuditLogRepository
}

func NewAuditUseCase(repo repository.AuditLogRepository) *AuditUseCase {
	return &AuditUseCase{repo: repo}
}

// Log records an audit event. Called from other use cases.
func (uc *AuditUseCase) Log(ctx context.Context, input LogInput) error {
	entry := domainaudit.Entry{
		ProjectID:  input.ProjectID,
		UserID:     input.UserID,
		Action:     input.Action,
		EntityType: input.EntityType,
		EntityID:   input.EntityID,
		Summary:    input.Summary,
		IP:         input.IP,
		UserAgent:  input.UserAgent,
	}
	if err := uc.repo.Log(ctx, entry); err != nil {
		return fmt.Errorf("audit log: %w", err)
	}
	return nil
}

// List retrieves paginated audit log entries with filters.
func (uc *AuditUseCase) List(ctx context.Context, input ListInput) (*ListOutput, error) {
	page, perPage := usecase.ClampPagination(input.Page, input.PerPage)
	filter := domainaudit.Filter{
		ProjectID:  input.ProjectID,
		UserID:     input.UserID,
		Action:     input.Action,
		EntityType: input.EntityType,
		Page:       page,
		PerPage:    perPage,
	}

	if input.DateFrom != nil {
		t, err := time.Parse(time.DateOnly, *input.DateFrom)
		if err != nil {
			return nil, fmt.Errorf("parse date_from: %w", err)
		}
		filter.DateFrom = &t
	}
	if input.DateTo != nil {
		t, err := time.Parse(time.DateOnly, *input.DateTo)
		if err != nil {
			return nil, fmt.Errorf("parse date_to: %w", err)
		}
		filter.DateTo = &t
	}

	entries, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}

	dtos := make([]EntryDTO, len(entries))
	for i, e := range entries {
		dtos[i] = EntryDTO{
			ID:         e.ID,
			ProjectID:  e.ProjectID,
			UserID:     e.UserID,
			Action:     e.Action,
			EntityType: e.EntityType,
			EntityID:   e.EntityID,
			Summary:    e.Summary,
			IP:         e.IP,
			UserAgent:  e.UserAgent,
			CreatedAt:  e.CreatedAt.Format(time.RFC3339),
		}
	}
	return &ListOutput{Entries: dtos, Total: total, Page: page, PerPage: perPage}, nil
}

// Record logs an audit event using metadata from context. Failures are logged but not returned.
func (uc *AuditUseCase) Record(ctx context.Context, projectID *string, action, entityType string, entityID *string, summary string) {
	entry := domainaudit.Entry{
		ProjectID:  projectID,
		UserID:     middleware.AuditUserID(ctx),
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Summary:    summary,
		IP:         middleware.AuditIP(ctx),
		UserAgent:  middleware.AuditUserAgent(ctx),
	}
	if err := uc.repo.Log(ctx, entry); err != nil {
		slog.Error("audit log failed", slog.String("action", action), slog.Any("err", err))
	}
}
