# ADR-006: Authorization Middleware

| Field      | Value                               |
| ---------- | ----------------------------------- |
| Status     | Accepted                            |
| Date       | 2026-02-22                          |
| Supersedes | —                                   |
| Components | observer, middleware, authorization |

## Decision

Authorization is a two-layer stack: a **platform role** (per user account) and a **project role** (per user-project pair), enforced by two independent middleware structs. Sensitive data fields are gated by explicit boolean flags on the project permission row, not by role alone.

## Platform Layer — `AuthMiddleware`

`internal/middleware/auth.go`

`Authenticate()` validates the access JWT (Bearer header or `access_token` cookie), loads the user, rejects deactivated or permanently locked accounts, and sets two context keys:

| Context key | Type        | Value                     |
| ----------- | ----------- | ------------------------- |
| `user_id`   | `ulid.ULID` | Authenticated user's ULID |
| `user_role` | `string`    | Platform role (see below) |

It also enriches the `context.Context` with audit metadata (user ID, client IP, user agent) via `WithAuditContext()`.

`RequireRole(roles ...user.Role)` checks that the user's platform role is in the allowed set, returning `403` otherwise. Used on admin-only routes.

### Platform roles

| Role         | Implied capabilities                                         |
| ------------ | ------------------------------------------------------------ |
| `admin`      | Full platform access; implicit project owner on all projects |
| `staff`      | Create/manage projects; read all cases                       |
| `consultant` | Assigned to projects; works with people records              |
| `guest`      | Read-only on explicitly assigned projects                    |

## Project Layer — `ProjectAuthMiddleware`

`internal/middleware/project_auth.go`

`RequireProjectRole(action project.Action)` runs after `Authenticate()` on all `/api/projects/:project_id/*` routes. Resolution order:

1. **Platform admin bypass** — role `admin` receives owner-level project access without a `project_permissions` row.
2. **Project owner bypass** — `projects.owner_id` match grants owner-level access without a `project_permissions` row.
3. **Permission lookup** — loads the `project_permissions` row for the user + project. If none exists, returns `403`. If the role rank is below the required minimum for the action, returns `403`.

Sets the following context keys after a successful check:

| Context key          | Type     | Source                                   |
| -------------------- | -------- | ---------------------------------------- |
| `project_id`         | `string` | Route param                              |
| `project_role`       | `string` | Resolved project role                    |
| `can_view_contact`   | `bool`   | `project_permissions.can_view_contact`   |
| `can_view_personal`  | `bool`   | `project_permissions.can_view_personal`  |
| `can_view_documents` | `bool`   | `project_permissions.can_view_documents` |
| `can_export`         | `bool`   | `project_permissions.can_export`         |

Platform admins and project owners receive `true` for all sensitivity flags.

### Project roles and action permissions

| Role         | read | create | update | delete | manage members     |
| ------------ | ---- | ------ | ------ | ------ | ------------------ |
| `viewer`     | ✓    |        |        |        |                    |
| `consultant` | ✓    | ✓      | ✓      |        |                    |
| `manager`    | ✓    | ✓      | ✓      | ✓      | ✓                  |
| `owner`      | ✓    | ✓      | ✓      | ✓      | ✓ + delete project |

Role comparison uses numeric rank: owner=4, manager=3, consultant=2, viewer=1. A user satisfies an action if their rank ≥ the action's minimum rank.

### Sensitivity flags

Data sensitivity is controlled independently of role. A `consultant` may read, create, and update records, but cannot see contact information, personal identifiers, or documents unless the corresponding flag is explicitly set on their `project_permissions` row.

| Flag                 | Gates                                                    |
| -------------------- | -------------------------------------------------------- |
| `can_view_contact`   | Phone numbers, email addresses, contact details          |
| `can_view_personal`  | National ID, external_id, sensitive personal identifiers |
| `can_view_documents` | Document list and file download                          |
| `can_export`         | CSV/export endpoints for people, support records, pets   |

### Export gate

`RequireExport()` is a separate middleware that reads `can_export` from context and returns `403` if false. It must be placed after `RequireProjectRole()` on export routes.

## Context Helpers

All context values are accessed via typed helper functions to avoid string key collisions:

```go
middleware.UserIDFrom(c)            // (ulid.ULID, bool)
middleware.UserRoleFrom(c)          // (user.Role, bool)
middleware.ProjectRoleFrom(c)       // (project.ProjectRole, bool)
middleware.CanViewContactFrom(c)    // bool
middleware.CanViewPersonalFrom(c)   // bool
middleware.CanViewDocumentsFrom(c)  // bool
middleware.CanExportFrom(c)         // bool
```
