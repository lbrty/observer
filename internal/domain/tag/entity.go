package tag

import "time"

// Tag represents a project-scoped label.
type Tag struct {
	ID        string
	ProjectID string
	Name      string
	Color     string
	CreatedAt time.Time
}
