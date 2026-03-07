package report

import "time"

// ReportFilter contains common filter parameters for all reports.
type ReportFilter struct {
	ProjectID    string
	DateFrom     *time.Time
	DateTo       *time.Time
	OfficeID     *string
	CategoryID   *string
	ConsultantID *string
	CaseStatus   *string
	Sex          *string
	AgeGroup     *string
	SupportType  *string
}

// CountResult represents a single count in a report breakdown.
type CountResult struct {
	Label string `json:"label" db:"label"`
	Count int    `json:"count" db:"count"`
}

// PetReportFilter contains filter parameters for pet reports.
type PetReportFilter struct {
	ProjectID string
	DateFrom  *time.Time
	DateTo    *time.Time
	Status    *string
}

// MonthlyStatusCount represents a monthly count for a specific status.
type MonthlyStatusCount struct {
	Month  string `json:"month" db:"month"`
	Status string `json:"status" db:"status"`
	Count  int    `json:"count" db:"count"`
}

// StatusFlow represents a transition between two case statuses.
type StatusFlow struct {
	FromStatus string  `json:"from_status" db:"from_status"`
	ToStatus   string  `json:"to_status" db:"to_status"`
	Count      int     `json:"count" db:"count"`
	AvgDays    float64 `json:"avg_days" db:"avg_days"`
}
