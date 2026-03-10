---
title: Permissions
weight: 6
---

Observer uses a two-level permission model. Every user has a **platform role** that controls access to admin functions. Within each project, users also have a **project role** that controls what they can do with case data. On top of this, three sensitivity flags restrict which fields appear in responses.

## Platform Roles

Platform roles are set when a user account is created and can only be changed by an admin.

| Role         | What they can do                                                       |
| ------------ | ---------------------------------------------------------------------- |
| `admin`      | Full access to everything — users, projects, reference data, all cases |
| `staff`      | Manage reference data and projects; cannot manage other users          |
| `consultant` | Work within projects they are assigned to; no admin panel access       |
| `guest`      | Read-only access to projects they are assigned to                      |

Platform admins automatically act as **project owner** on all projects — they bypass project-level permission checks entirely.

## Project Roles

Project roles apply within a specific project. A user must be explicitly assigned to a project to access it.

| Role         | Rank | Can do                                         |
| ------------ | ---- | ---------------------------------------------- |
| `owner`      | 4    | Everything, including deleting the project     |
| `manager`    | 3    | All case data operations + manage team members |
| `consultant` | 2    | Create, read, update cases and export data     |
| `viewer`     | 1    | Read-only access to case data                  |

### Action Requirements

| Action         | Minimum project role |
| -------------- | -------------------- |
| Read data      | `viewer`             |
| Create records | `consultant`         |
| Update records | `consultant`         |
| Delete records | `manager`            |
| Manage members | `manager`            |
| Export data    | `consultant`         |

## Sensitivity Flags

Each project permission has three independent boolean flags. These are set per-user, per-project by a manager or owner.

| Flag                 | Controls                                                         |
| -------------------- | ---------------------------------------------------------------- |
| `can_view_contact`   | Phone numbers and email addresses in person records              |
| `can_view_personal`  | Full name, birth date, national ID (`external_id`), consent info |
| `can_view_documents` | Document file access and document metadata                       |
| `can_export`         | Access to all CSV export endpoints for this project              |

When a flag is off, the corresponding fields are omitted from API responses and CSV exports — the data stays in the database but is not sent to that user. A consultant with `can_view_personal: false` cannot recover national IDs or birth dates through an export even though they can create records.

`can_export` is an all-or-nothing gate: without it, the export endpoints return 403 regardless of project role. Platform admins and project owners always have export access. Staff platform role receives export access automatically without needing the flag set.

## Assigning Permissions

Only platform admins and project owners/managers can assign project permissions.

### Assign a user to a project

```http
POST /admin/projects/:project_id/permissions
```

```json
{
  "user_id": "01J...",
  "role": "consultant",
  "can_view_contact": true,
  "can_view_personal": false,
  "can_view_documents": false
}
```

### Update an existing permission

```http
PUT /admin/projects/:project_id/permissions/:permission_id
```

### Remove a user from a project

```http
DELETE /admin/projects/:project_id/permissions/:permission_id
```

## Common Setups

### Field worker (consultant, restricted personal data)

```json
{
  "role": "consultant",
  "can_view_contact": true,
  "can_view_personal": false,
  "can_view_documents": false
}
```

Phone numbers visible for coordination; no access to national IDs or birth dates.

### Supervisor (manager, full access)

```json
{
  "role": "manager",
  "can_view_contact": true,
  "can_view_personal": true,
  "can_view_documents": true
}
```

### External auditor (viewer, no sensitive data)

```json
{
  "role": "viewer",
  "can_view_contact": false,
  "can_view_personal": false,
  "can_view_documents": false
}
```
