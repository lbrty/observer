package project

import "time"

// ProjectStatus represents the lifecycle state of a project.
type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "active"
	ProjectStatusArchived ProjectStatus = "archived"
	ProjectStatusClosed   ProjectStatus = "closed"
)

// Project represents a top-level organizational unit.
type Project struct {
	ID          string
	Name        string
	Description *string
	OwnerID     string
	Status      ProjectStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ProjectListFilter holds optional filters for listing projects.
type ProjectListFilter struct {
	OwnerID *string
	Status  *ProjectStatus
	Page    int
	PerPage int
}
