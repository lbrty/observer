package report

// PetReportInput is the query input for pet report endpoints.
type PetReportInput struct {
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
	Status   string `form:"status"`
}

// MonthlyStatusCountDTO is the response row for monthly status counts.
type MonthlyStatusCountDTO struct {
	Month  string `json:"month"`
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// PetReportOutput returns all pet report groups at once.
type PetReportOutput struct {
	ByStatus        ReportOutput             `json:"by_status"`
	ByOwnership     ReportOutput             `json:"by_ownership"`
	ByMonth         ReportOutput             `json:"by_month"`
	ByStatusByMonth []MonthlyStatusCountDTO  `json:"by_status_by_month"`
}
