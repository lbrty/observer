package report

import (
	"context"
	"fmt"
	"time"

	domainreport "github.com/lbrty/observer/internal/domain/report"
	"github.com/lbrty/observer/internal/repository"
)

// ReportUseCase generates ADR-005 reports for a project.
type ReportUseCase struct {
	repo repository.ReportRepository
}

// NewReportUseCase creates a ReportUseCase.
func NewReportUseCase(repo repository.ReportRepository) *ReportUseCase {
	return &ReportUseCase{repo: repo}
}

// Generate runs all 10 report groups and returns the combined result.
func (uc *ReportUseCase) Generate(ctx context.Context, projectID string, input ReportInput) (*FullReportOutput, error) {
	f := domainreport.ReportFilter{ProjectID: projectID}

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
	if input.OfficeID != "" {
		f.OfficeID = &input.OfficeID
	}
	if input.CategoryID != "" {
		f.CategoryID = &input.CategoryID
	}
	if input.ConsultantID != "" {
		f.ConsultantID = &input.ConsultantID
	}
	if input.CaseStatus != "" {
		f.CaseStatus = &input.CaseStatus
	}
	if input.Sex != "" {
		f.Sex = &input.Sex
	}
	if input.AgeGroup != "" {
		f.AgeGroup = &input.AgeGroup
	}
	if input.SupportType != "" {
		f.SupportType = &input.SupportType
	}

	out := &FullReportOutput{}

	consultations, err := uc.repo.CountConsultations(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("consultations report: %w", err)
	}
	out.Consultations = toOutput("consultations", consultations)

	bySex, err := uc.repo.CountBySex(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("sex report: %w", err)
	}
	out.BySex = toOutput("by_sex", bySex)

	byIDP, err := uc.repo.CountByIDPStatus(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("idp report: %w", err)
	}
	out.ByIDPStatus = toOutput("by_idp_status", byIDP)

	byCat, err := uc.repo.CountByCategory(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("category report: %w", err)
	}
	out.ByCategory = toOutput("by_category", byCat)

	byRegion, err := uc.repo.CountByCurrentRegion(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("region report: %w", err)
	}
	out.ByRegion = toOutput("by_region", byRegion)

	bySphere, err := uc.repo.CountBySphere(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("sphere report: %w", err)
	}
	out.BySphere = toOutput("by_sphere", bySphere)

	peopleBySphere, err := uc.repo.CountPeopleBySphere(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("people by sphere report: %w", err)
	}
	out.PeopleBySphere = toOutput("people_by_sphere", peopleBySphere)

	byOffice, err := uc.repo.CountByOffice(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("office report: %w", err)
	}
	out.ByOffice = toOutput("by_office", byOffice)

	byAge, err := uc.repo.CountByAgeGroup(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("age report: %w", err)
	}
	out.ByAgeGroup = toOutput("by_age_group", byAge)

	consultsByAge, err := uc.repo.CountConsultationsByAgeGroup(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("consultations by age report: %w", err)
	}
	out.ConsultationsByAgeGroup = toOutput("consultations_by_age_group", consultsByAge)

	byTag, err := uc.repo.CountByTag(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("tag report: %w", err)
	}
	out.ByTag = toOutput("by_tag", byTag)

	families, err := uc.repo.CountFamilyUnits(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("family report: %w", err)
	}
	out.FamilyUnits = toOutput("family_units", families)

	byCaseStatus, err := uc.repo.CountByCaseStatus(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("case status report: %w", err)
	}
	out.ByCaseStatus = toOutput("by_case_status", byCaseStatus)

	flows, err := uc.repo.StatusFlowReport(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("status flow report: %w", err)
	}
	dtos := make([]StatusFlowDTO, len(flows))
	for i, fl := range flows {
		dtos[i] = StatusFlowDTO{
			FromStatus: fl.FromStatus,
			ToStatus:   fl.ToStatus,
			Count:      fl.Count,
			AvgDays:    fl.AvgDays,
		}
	}
	out.StatusFlow = dtos

	return out, nil
}

// validDimensions is the set of allowed grouping dimensions.
var validDimensions = map[string]bool{
	"sex":           true,
	"age_group":     true,
	"region":        true,
	"conflict_zone": true,
	"office":        true,
	"sphere":        true,
	"category":      true,
	"person_tag":    true,
	"pet_status":    true,
	"pet_tag":       true,
}

// GenerateCustom runs the custom report builder.
func (uc *ReportUseCase) GenerateCustom(ctx context.Context, projectID string, input CustomReportInput) (*CustomReportOutput, error) {
	if len(input.GroupBy) == 0 {
		return nil, fmt.Errorf("at least one dimension required in group_by")
	}
	if len(input.GroupBy) > 2 {
		return nil, fmt.Errorf("at most 2 dimensions allowed in group_by")
	}
	for _, dim := range input.GroupBy {
		if !validDimensions[dim] {
			return nil, fmt.Errorf("unknown dimension: %s", dim)
		}
	}

	f := domainreport.ReportFilter{ProjectID: projectID}

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
	if input.SupportType != "" {
		f.SupportType = &input.SupportType
	}
	if input.OfficeID != "" {
		f.OfficeID = &input.OfficeID
	}
	if input.CategoryID != "" {
		f.CategoryID = &input.CategoryID
	}
	if input.CaseStatus != "" {
		f.CaseStatus = &input.CaseStatus
	}
	if input.Sex != "" {
		f.Sex = &input.Sex
	}

	results, total, err := uc.repo.CustomQuery(ctx, projectID, input.Metric, input.GroupBy, f)
	if err != nil {
		return nil, fmt.Errorf("custom report: %w", err)
	}

	rows := make([]CustomRow, len(results))
	for i, r := range results {
		rows[i] = CustomRow{Dimensions: r.Dimensions, Count: r.Count}
	}

	return &CustomReportOutput{
		Metric:  input.Metric,
		GroupBy: input.GroupBy,
		Rows:    rows,
		Total:   total,
	}, nil
}

func toOutput(group string, rows []domainreport.CountResult) ReportOutput {
	total := 0
	dtos := make([]CountResultDTO, len(rows))
	for i, r := range rows {
		dtos[i] = CountResultDTO{Label: r.Label, Count: r.Count}
		total += r.Count
	}
	return ReportOutput{Group: group, Rows: dtos, Total: total}
}
