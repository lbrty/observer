package report

// ReportInput is the query input for all report endpoints.
type ReportInput struct {
	DateFrom     string `form:"date_from"`
	DateTo       string `form:"date_to"`
	OfficeID     string `form:"office_id"`
	CategoryID   string `form:"category_id"`
	ConsultantID string `form:"consultant_id"`
	CaseStatus   string `form:"case_status"`
	Sex          string `form:"sex"`
	AgeGroup     string `form:"age_group"`
	SupportType  string `form:"support_type"`
}

// CountResultDTO is the response row.
type CountResultDTO struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

// ReportOutput wraps a single report group result.
type ReportOutput struct {
	Group string           `json:"group"`
	Rows  []CountResultDTO `json:"rows"`
	Total int              `json:"total"`
}

// StatusFlowDTO is the response row for status flow transitions.
type StatusFlowDTO struct {
	FromStatus string  `json:"from_status"`
	ToStatus   string  `json:"to_status"`
	Count      int     `json:"count"`
	AvgDays    float64 `json:"avg_days"`
}

// FullReportOutput returns all report groups at once.
type FullReportOutput struct {
	Consultations ReportOutput    `json:"consultations"`
	BySex         ReportOutput    `json:"by_sex"`
	ByIDPStatus   ReportOutput    `json:"by_idp_status"`
	ByCategory    ReportOutput    `json:"by_category"`
	ByRegion      ReportOutput    `json:"by_region"`
	BySphere      ReportOutput    `json:"by_sphere"`
	ByOffice      ReportOutput    `json:"by_office"`
	ByAgeGroup    ReportOutput    `json:"by_age_group"`
	ByTag         ReportOutput    `json:"by_tag"`
	FamilyUnits   ReportOutput    `json:"family_units"`
	ByCaseStatus  ReportOutput    `json:"by_case_status"`
	StatusFlow    []StatusFlowDTO `json:"status_flow"`
}
