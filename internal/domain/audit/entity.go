package audit

import "time"

type Entry struct {
	ID         string
	ProjectID  *string
	UserID     *string // nil when the user has been deleted
	Action     string
	EntityType string
	EntityID   *string
	Summary    string
	IP         *string
	UserAgent  *string
	CreatedAt  time.Time

	// Populated by repository reads via LEFT JOIN with users; zero-valued on Log.
	UserFirstName string
	UserLastName  string
	UserEmail     string
}

type Filter struct {
	ProjectID  *string
	UserID     *string
	Action     *string
	EntityType *string
	DateFrom   *time.Time
	DateTo     *time.Time
	Page       int
	PerPage    int
}
