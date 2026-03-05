package report

// ReportInput is the query input for all report endpoints.
type ReportInput struct {
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
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

// FullReportOutput returns all 10 groups at once.
type FullReportOutput struct {
	Consultations ReportOutput `json:"consultations"`
	BySex         ReportOutput `json:"by_sex"`
	ByIDPStatus   ReportOutput `json:"by_idp_status"`
	ByCategory    ReportOutput `json:"by_category"`
	ByRegion      ReportOutput `json:"by_region"`
	BySphere      ReportOutput `json:"by_sphere"`
	ByOffice      ReportOutput `json:"by_office"`
	ByAgeGroup    ReportOutput `json:"by_age_group"`
	ByTag         ReportOutput `json:"by_tag"`
	FamilyUnits   ReportOutput `json:"family_units"`
}
