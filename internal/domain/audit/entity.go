package audit

import "time"

type Entry struct {
	ID         string    `db:"id"`
	ProjectID  *string   `db:"project_id"`
	UserID     string    `db:"user_id"`
	Action     string    `db:"action"`
	EntityType string    `db:"entity_type"`
	EntityID   *string   `db:"entity_id"`
	Summary    string    `db:"summary"`
	IP         string    `db:"ip"`
	UserAgent  string    `db:"user_agent"`
	CreatedAt  time.Time `db:"created_at"`
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
