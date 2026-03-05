# ADR-009: Deferred Features (Phase 2)

| Field      | Value                                  |
| ---------- | -------------------------------------- |
| Status     | Accepted                               |
| Date       | 2026-02-27                             |
| Supersedes | —                                      |
| Components | observer, crypto, auth, reports, audit |

---

## Context

The initial release prioritises core case-management workflows: person registration, support records, migration tracking, households, documents, pets, project-scoped RBAC, and ADR-005 aggregate reports. Several features were discussed during design but deferred to keep the MVP focused and shippable.

This ADR records what was deferred and why, so the decisions are visible and can be revisited.

## Deferred Features

### 1. MEK/DEK Envelope Encryption

**What:** Encrypt PII columns (names, phone, external_id) at rest using a Master Encryption Key / Data Encryption Key scheme. Each project would have its own DEK, wrapped by a MEK stored in a KMS or HSM.

**Why deferred:** Adds significant complexity to the repository layer (transparent encrypt/decrypt), key rotation procedures, and backup/restore workflows. The current deployment model (self-hosted Postgres with disk-level encryption) provides adequate protection for the MVP target audience.

**Prerequisite for Phase 2:** Decide on KMS backend (AWS KMS, HashiCorp Vault, or local file-based).

### 2. Detailed Audit Logs

**What:** Record every data mutation (create, update, delete) with actor, timestamp, old/new values, and IP address in a dedicated `audit_log` table.

**Why deferred:** Requires a cross-cutting middleware or repository decorator pattern. The current `created_at` / `updated_at` timestamps and request logging provide basic traceability. Full audit logging should be designed alongside the encryption layer to avoid logging plaintext PII.

### 3. MFA Enforcement

**What:** Require TOTP-based multi-factor authentication for admin and staff roles. The `mfa_configs` table and `MFARepository` already exist; the login flow detects MFA and returns a `mfa_required` response. The missing piece is the frontend TOTP setup and verification flow.

**Why deferred:** The backend plumbing is in place. Frontend implementation requires a QR code generator, recovery codes display, and a verification input step during login. Deferred to avoid scope creep in the initial release.

### 4. Email Verification

**What:** Send a verification email on registration with a one-time token link. The `verification_tokens` table exists but is unused.

**Why deferred:** The target deployment (small NGO, internal users) uses admin-approval registration (`is_active = false` on signup). Email verification adds an SMTP dependency and is unnecessary when an admin manually activates accounts.

### 5. CSV / PDF Report Export

**What:** Allow downloading ADR-005 reports as CSV or PDF files from the reports page.

**Why deferred:** The current JSON API + D3 visualisation covers the primary use case (on-screen review). Export requires a server-side rendering library (PDF) or CSV serialisation endpoint. Can be added as a thin handler on top of the existing `ReportRepository` without architectural changes.

### 6. File Upload / Storage

**What:** Store actual document files (scans, photos) referenced by the `documents` table. Currently documents store metadata only.

**Why deferred:** Requires choosing a storage backend (local filesystem, S3-compatible, or database BLOBs), adding upload/download endpoints, virus scanning considerations, and access-control integration with the existing `can_view_documents` permission flag.

## Consequences

- These features are **not blocked** — the schema and interfaces are designed to accommodate them.
- Phase 2 work should begin with MEK/DEK and audit logs together, since they share concerns around PII handling.
- MFA frontend can be implemented independently at any time.
- Export and file storage are isolated features with no cross-cutting impact.
