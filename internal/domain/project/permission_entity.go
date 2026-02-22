package project

import "time"

// ProjectPermission is the full persisted project permission record.
type ProjectPermission struct {
	ID               string
	ProjectID        string
	UserID           string
	Role             ProjectRole
	CanViewContact   bool
	CanViewPersonal  bool
	CanViewDocuments bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ValidateProjectRole checks that a role string is a valid project role.
func ValidateProjectRole(role string) (ProjectRole, error) {
	switch ProjectRole(role) {
	case ProjectRoleOwner, ProjectRoleManager, ProjectRoleConsultant, ProjectRoleViewer:
		return ProjectRole(role), nil
	default:
		return "", ErrInvalidProjectRole
	}
}
