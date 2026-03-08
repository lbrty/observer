package project

// ProjectRole represents a user's role within a project.
type ProjectRole string

const (
	ProjectRoleOwner      ProjectRole = "owner"
	ProjectRoleManager    ProjectRole = "manager"
	ProjectRoleConsultant ProjectRole = "consultant"
	ProjectRoleViewer     ProjectRole = "viewer"
)

// Rank returns numeric rank for hierarchy comparison.
func (r ProjectRole) Rank() int {
	switch r {
	case ProjectRoleOwner:
		return 4
	case ProjectRoleManager:
		return 3
	case ProjectRoleConsultant:
		return 2
	case ProjectRoleViewer:
		return 1
	default:
		return 0
	}
}

// Action represents a project-scoped operation.
type Action string

const (
	ActionRead          Action = "read"
	ActionCreate        Action = "create"
	ActionUpdate        Action = "update"
	ActionDelete        Action = "delete"
	ActionManageMembers Action = "manage_members"
	ActionExport        Action = "export"
)

// MinRoleForAction maps each action to its minimum required project role.
var MinRoleForAction = map[Action]ProjectRole{
	ActionRead:          ProjectRoleViewer,
	ActionCreate:        ProjectRoleConsultant,
	ActionUpdate:        ProjectRoleConsultant,
	ActionDelete:        ProjectRoleManager,
	ActionManageMembers: ProjectRoleManager,
	ActionExport:        ProjectRoleConsultant,
}

// Permission holds a user's project-level role and sensitivity flags.
type Permission struct {
	Role             ProjectRole
	CanViewContact   bool
	CanViewPersonal  bool
	CanViewDocuments bool
}
