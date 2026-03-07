package report

import (
	"context"
	"fmt"
	"time"

	domainreport "github.com/lbrty/observer/internal/domain/report"
	"github.com/lbrty/observer/internal/repository"
)

// PetReportUseCase generates pet reports for a project.
type PetReportUseCase struct {
	repo repository.PetReportRepository
}

// NewPetReportUseCase creates a PetReportUseCase.
func NewPetReportUseCase(repo repository.PetReportRepository) *PetReportUseCase {
	return &PetReportUseCase{repo: repo}
}

// Generate runs all pet report queries and returns the combined result.
func (uc *PetReportUseCase) Generate(ctx context.Context, projectID string, input PetReportInput) (*PetReportOutput, error) {
	f := domainreport.PetReportFilter{ProjectID: projectID}

	if input.DateFrom != "" {
		t, err := time.Parse("2006-01-02", input.DateFrom)
		if err != nil {
			return nil, fmt.Errorf("parse date_from: %w", err)
		}
		f.DateFrom = &t
	}
	if input.DateTo != "" {
		t, err := time.Parse("2006-01-02", input.DateTo)
		if err != nil {
			return nil, fmt.Errorf("parse date_to: %w", err)
		}
		f.DateTo = &t
	}
	if input.Status != "" {
		f.Status = &input.Status
	}

	out := &PetReportOutput{}

	byStatus, err := uc.repo.CountByStatus(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("pet status report: %w", err)
	}
	out.ByStatus = toOutput("by_status", byStatus)

	byOwnership, err := uc.repo.CountByOwnership(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("pet ownership report: %w", err)
	}
	out.ByOwnership = toOutput("by_ownership", byOwnership)

	byMonth, err := uc.repo.CountByMonth(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("pet monthly report: %w", err)
	}
	out.ByMonth = toOutput("by_month", byMonth)

	byStatusMonth, err := uc.repo.CountByStatusByMonth(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("pet status by month report: %w", err)
	}
	dtos := make([]MonthlyStatusCountDTO, len(byStatusMonth))
	for i, r := range byStatusMonth {
		dtos[i] = MonthlyStatusCountDTO{
			Month:  r.Month,
			Status: r.Status,
			Count:  r.Count,
		}
	}
	out.ByStatusByMonth = dtos

	return out, nil
}
