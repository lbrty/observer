---
title: Audit Log
weight: 9
---

Observer records an audit log entry for every create, update, and delete operation on case data. Audit logs are append-only — they cannot be edited or deleted.

## What Gets Logged

| Category          | Operations logged                                   |
| ----------------- | --------------------------------------------------- |
| People            | Create, update, delete person records               |
| Support records   | Create, update, delete consultations                |
| Migration records | Create, update, delete movement records             |
| Households        | Create, update, delete household and member records |
| Notes             | Create, update, delete case notes                   |
| Documents         | Upload, update metadata, delete documents           |
| Pets              | Create, update, delete pet records                  |
| Permissions       | Assign, update, revoke project permissions          |

Authentication events (login, logout, token refresh) are not in the project audit log — they appear in server logs.

## Viewing the Audit Log

Only project managers and owners can access the audit log.

```http
GET /projects/:project_id/audit?page=1&per_page=50
```

| Parameter  | Description                        |
| ---------- | ---------------------------------- |
| `page`     | Page number (default 1)            |
| `per_page` | Results per page (default 50)      |
| `actor_id` | Filter by user who made the change |
| `start`    | Filter by date (YYYY-MM-DD)        |
| `end`      | Filter by date (YYYY-MM-DD)        |

### Response format

Each audit entry contains:

```json
{
  "id": "01J...",
  "project_id": "01J...",
  "actor_id": "01J...",
  "actor_ip": "192.168.1.1",
  "action": "create",
  "entity_type": "person",
  "entity_id": "01J...",
  "created_at": "2024-06-15T10:30:00Z"
}
```

| Field         | Description                                            |
| ------------- | ------------------------------------------------------ |
| `actor_id`    | User who performed the action                          |
| `actor_ip`    | IP address of the request                              |
| `action`      | `create`, `update`, or `delete`                        |
| `entity_type` | Record type (`person`, `support_record`, `note`, etc.) |
| `entity_id`   | ULID of the affected record                            |

## Audit in the Web Interface

1. Open a project
2. Click **Audit Log** in the navigation
3. Filter by date range or user

{{% callout type="note" %}}
Audit log access requires the `manager` or `owner` project role. Consultants and viewers cannot see the audit log.
{{% /callout %}}
