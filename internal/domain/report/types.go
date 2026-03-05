package report

import "time"

// ReportFilter contains common filter parameters for all reports.
type ReportFilter struct {
	ProjectID string
	DateFrom  *time.Time
	DateTo    *time.Time
}

// CountResult represents a single count in a report breakdown.
type CountResult struct {
	Label string `json:"label" db:"label"`
	Count int    `json:"count" db:"count"`
}
