# Observer — Security Developer Guide

This document explains the security controls in Observer, how they work, and what to keep in mind when adding new features.

---

## Authentication

### Password policy

Passwords are validated at the binding layer using two rules:

```go
binding:"required,min=12,strongpassword"
```

`strongpassword` is a custom Gin validator (registered in `server.go`) that requires at least one digit **and** one special character. This applies to:

- `RegisterInput.Password`
- `ChangePasswordInput.NewPassword`

The hasher is Argon2id with `time=2, memory=64MB, threads=4`. Parameters are in `internal/crypto/password.go`. If you tune them, check [OWASP Argon2id recommendations](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html).

### Login lockout

Failed logins are tracked per email in Redis (`internal/repository/login_attempt_store.go`).

| Threshold            | Effect                                 |
| -------------------- | -------------------------------------- |
| 5 failures in 15 min | Temporary lock (15 min)                |
| 10 failures total    | Permanent lock — requires admin unlock |

Permanent locks are **also written to `users.locked_permanently_at`** (migration `000034`) so they survive Redis restarts. The `Authenticate()` middleware checks this column on every request.

The login endpoint fails **closed**: if Redis is unavailable, access is denied, not granted.

Error messages intentionally omit lock duration to slow down enumeration:

```
"account temporarily locked, please try again later"
"account locked, contact administrator"
```

### Email enumeration prevention

`Register` returns the same success message regardless of whether the email is already taken. The duplicate-email path returns silently without creating a new user.

### Role enforcement

New users are always created with `RoleGuest` regardless of what the request body contains. Role promotion only happens through `PATCH /admin/users/:id` (admin-only).

---

## Sessions and tokens

### Access token

JWT signed with RS256 (RSA private key at `keys/jwt_rsa`). Default TTL: 15 minutes. Contains `user_id` and `role`.

### Refresh token

64 random bytes (512 bits) stored in the `sessions` table. **Rotated on every use** — the old token is deleted and a new one issued. Reuse of a consumed token returns 401.

### MFA token

Short-lived JWT (default 5 minutes) issued when a user with MFA enabled completes password verification. Does not create a session. Must be exchanged for real tokens via `POST /auth/mfa`.

### Session vacuum

A background goroutine (`app.StartSessionVacuum`) periodically deletes expired sessions from the database. Wired in `cmd/observer/cmd/serve.go`.

### Password change / reset

Changing or resetting a password calls `sessionRepo.DeleteByUserID` to invalidate all existing sessions. The user must log in again.

---

## Multi-factor authentication (TOTP)

MFA uses RFC 6238 TOTP (30-second window). The TOTP secret is stored in `mfa_configs.secret` but is never transmitted after initial setup.

### Recovery codes

On `EnableMFA`, 8 single-use recovery codes are generated:

```
plain text → Argon2id hash → stored in mfa_recovery_codes
plain text → returned in EnableMFAOutput.RecoveryCodes (shown once)
```

The plain-text codes are **never stored** — show them to the user exactly once. If lost, the admin must disable MFA via the admin API or directly in the database.

On `VerifyMFA`, if the TOTP code fails, it is tried as a recovery code. On match, the code is marked `used_at = NOW()` and cannot be reused.

On `DisableMFA`, all recovery codes for the user are deleted.

---

## CSRF protection

Observer uses the **double-submit cookie** pattern (`internal/middleware/csrf.go`).

At login/refresh/MFA-verify, three cookies are set:

| Cookie          | HttpOnly  | Purpose                             |
| --------------- | --------- | ----------------------------------- |
| `access_token`  | true      | Bearer credential                   |
| `refresh_token` | true      | Rotation credential (path: `/auth`) |
| `csrf_token`    | **false** | Readable by JS                      |

For state-mutating requests (POST/PUT/PATCH/DELETE), the middleware requires:

1. `csrf_token` cookie is present
2. `X-CSRF-Token` header matches the cookie value (compared with `crypto/subtle.ConstantTimeCompare`)

Safe methods (GET/HEAD/OPTIONS) bypass this check.

The frontend reads `csrf_token` from document.cookie and sends it as `X-CSRF-Token`. This is wired in `packages/observer-web/src/lib/api.ts` via the `ky` `beforeRequest` hook.

**Why this works:** An attacker on a different origin cannot read cookies (SameSite + HttpOnly protects auth tokens) and cannot set custom headers in a cross-origin form or fetch. Reading the cookie value requires same-origin access, which CSRF attacks do not have.

---

## CORS

Origins are configured via `CORS_ORIGINS` (comma-separated). At startup, the server logs a warning if any origin contains `localhost` or `127.0.0.1` — these should not be present in production.

`AllowCredentials: true` is required for cookie-based auth. Be precise with `AllowOrigins` — never use `"*"` with credentials.

---

## Authorization

### Platform roles

| Role         | Description                                                                       |
| ------------ | --------------------------------------------------------------------------------- |
| `admin`      | Full access to all endpoints and projects                                         |
| `staff`      | Read/write for users and reference data; project access controlled by permissions |
| `consultant` | Project access controlled by permissions                                          |
| `guest`      | Project access controlled by permissions (default for new registrations)          |

Role is verified on every request by `Authenticate()` → JWT claim.

### Project roles

Project-scoped access is controlled by the `permissions` table. The `RequireProjectRole(action)` middleware:

1. Bypasses checks for `admin` role
2. Checks `IsProjectOwner` (owners have full access)
3. Loads the user's `ProjectPermission` record
4. Compares `perm.Role.Rank()` against `MinRoleForAction[action]`

Permission flags are stored in context:

| Context key          | Source                           |
| -------------------- | -------------------------------- |
| `can_view_contact`   | `permissions.can_view_contact`   |
| `can_view_personal`  | `permissions.can_view_personal`  |
| `can_view_documents` | `permissions.can_view_documents` |
| `can_export`         | `permissions.can_export`         |

**Export access** is enforced at the route level by `RequireExport()` middleware on the `/export/*` route group. Do not add manual `CanExportFrom(c)` checks in handlers — use the middleware.

### IDOR prevention

Every use-case method that fetches an entity by ID must verify that the entity belongs to the project (or person) in the URL. Pattern:

```go
record, err := repo.GetByID(ctx, id)
if record.ProjectID != projectID {
    return nil, project.ErrPermissionNotFound  // returns 404, not 403
}
```

Returning 404 (not 403) on ownership mismatch prevents attackers from confirming that a record exists in another project.

**When adding a new entity**: add an IDOR check in the `Get`, `Update`, and `Delete` use-case methods. See `internal/usecase/project/support_record_usecase.go` for reference.

---

## Input validation

### Request body size

Non-multipart request bodies are limited to **1 MB** via `http.MaxBytesReader` (wired in `server.go`). The limit is intentionally not applied to multipart/file-upload requests — document uploads enforce a **50 MB** cap in `document_handler.go`.

Exceeding the limit returns HTTP 413.

### Pagination clamping

`usecase.ClampPagination` must be called **before** constructing the filter struct. Calling it after silently ignores the clamped values. Max `per_page` is 100 to prevent unbounded queries.

```go
// correct
page, perPage := usecase.ClampPagination(input.Page, input.PerPage)
filter := PersonListFilter{Page: page, PerPage: perPage, ...}
```

---

## File uploads and serving

### Upload (`POST /projects/:id/documents`)

1. Filename is sanitized with `filepath.Base` to strip any path traversal (`../../etc/passwd` → `passwd`).
2. MIME type is detected server-side using `http.DetectContentType` on the first 512 bytes (not trusted from `Content-Type` header).
3. `text/html` and `application/xhtml+xml` are rejected with 400 to prevent stored XSS via document serving.
4. The sanitized name and server-detected MIME are stored; the original filename from the client is discarded.

### Download

`Content-Disposition: attachment; filename="..."` is set using `mime.FormatMediaType` to ensure proper encoding of special characters in filenames.

### Stream (inline)

Only a restricted set of MIME types are served inline:

```
image/*, video/*, audio/*, application/pdf, text/plain
```

Everything else is served as attachment. CSP is set to `frame-ancestors 'self'` for inline serving.

---

## Security headers

The `SecurityHeaders()` middleware (`internal/middleware/security_headers.go`) sets:

| Header                    | Value                                                                |
| ------------------------- | -------------------------------------------------------------------- |
| `X-Content-Type-Options`  | `nosniff`                                                            |
| `X-Frame-Options`         | `DENY`                                                               |
| `Referrer-Policy`         | `strict-origin-when-cross-origin`                                    |
| `Content-Security-Policy` | `default-src 'self'; object-src 'none'; frame-ancestors 'none'; ...` |

**CSP note:** `style-src 'self' 'unsafe-inline'` is present to support Tailwind utility classes injected at runtime. If you add a strict CSP in future, use nonces.

---

## Error handling

`HandleError` in `internal/handler/errors.go` maps domain errors to HTTP status codes. For unrecognized errors (5xx), it:

1. Logs the real error with `slog.Error`
2. Returns `"internal server error"` to the client (never the original error message)

This prevents leaking database errors, stack traces, or internal paths to clients. Never do `c.JSON(500, err.Error())` directly in handlers.

---

## Audit logging

Sensitive operations are recorded via `auditUC.Record(ctx, projectID, action, entity, entityID, note)`. Currently instrumented:

- User create, password reset (admin)
- Permission updates
- Tag/person/support-record/document/pet create, update, delete (project scope)
- CSV exports

When adding features that touch PII or access control, add an audit record.

---

## Frontend

### API base URL enforcement

In non-dev/test builds, `vite.config.ts` asserts that `VITE_API_URL` starts with `https://`. This prevents accidental mixed-content deployments.

### CSRF token injection

The `ky` instance in `src/lib/api.ts` reads `csrf_token` from the cookie on every state-mutating request and sends it as `X-CSRF-Token`. This is transparent to all callers — no manual header injection needed.

---

## Adding a new endpoint — security checklist

When adding a new route, go through this checklist:

- [ ] Route is placed under the correct middleware group (unauthenticated / authenticated / role-gated / project-scoped)
- [ ] If project-scoped: uses `RequireProjectRole(action)` with the minimum required action
- [ ] If export: uses `RequireExport()` middleware, not a manual `CanExportFrom` check
- [ ] Use-case methods check ownership before returning or modifying entities (IDOR)
- [ ] Pagination calls `ClampPagination` before building the filter
- [ ] New passwords use `binding:"required,min=12,strongpassword"`
- [ ] File uploads sanitize filename and detect MIME server-side
- [ ] 5xx errors are logged and return `"internal server error"` to the client
- [ ] Sensitive operations emit an audit record
- [ ] New repository interfaces are added to `go:generate` and mocks regenerated
