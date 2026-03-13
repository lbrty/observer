# ADR-002: Users and Auth

| Field      | Value      |
| ---------- | ---------- |
| Status     | Accepted   |
| Date       | 2026-02-21 |
| Supersedes | —          |
| Components | auth, mfa  |

## Decision

Authentication uses RSA-signed JWTs with a two-token scheme (short-lived access + long-lived refresh). Passwords are hashed with Argon2id. MFA uses TOTP with single-use backup recovery codes.

## User Entity

```go
type User struct {
    ID                  ulid.ULID
    FirstName           string
    LastName            string
    Email               string
    Phone               string
    OfficeID            *string
    Role                Role       // platform role
    IsVerified          bool
    IsActive            bool
    DeactivatedAt       *time.Time // soft deactivation (reversible)
    LockedPermanentlyAt *time.Time // permanent lock (irreversible via UI)
    CreatedAt           time.Time
    UpdatedAt           time.Time
}
```

### Platform roles

| Role         | Access                                          |
| ------------ | ----------------------------------------------- |
| `admin`      | Full platform access, user management           |
| `staff`      | Create/manage projects, view all cases          |
| `consultant` | Assigned to projects, works with people records |
| `guest`      | Read-only on explicitly assigned projects       |

## Credentials

Passwords are stored in a separate `credentials` table (not on the user row):

| Column          | Type        | Notes                             |
| --------------- | ----------- | --------------------------------- |
| `user_id`       | TEXT        | FK → users                        |
| `password_hash` | TEXT        | Argon2id output (includes params) |
| `salt`          | TEXT        | Per-user random salt              |
| `updated_at`    | TIMESTAMPTZ |                                   |

Argon2id parameters follow OWASP recommendations. The `crypto.ArgonHasher` implements the `crypto.PasswordHasher` interface.

## JWT Tokens

Tokens are RSA-signed (minimum 4096-bit key). Keys are loaded from files at startup via `JWT_PRIVATE_KEY_PATH` and `JWT_PUBLIC_KEY_PATH`.

| Token    | TTL    | Env variable       | Claims                          |
| -------- | ------ | ------------------ | ------------------------------- |
| Access   | 15 min | `JWT_ACCESS_TTL`   | `user_id`, `role`, `exp`, `iss` |
| Refresh  | 7 days | `JWT_REFRESH_TTL`  | `user_id`, `exp`, `iss`         |
| MFA temp | 5 min  | `JWT_MFA_TEMP_TTL` | `user_id`, `mfa_pending`, `exp` |

**Refresh token rotation:** every `/api/auth/refresh` call issues a new refresh token and invalidates the old one. The token value is stored hashed in the `sessions` table.

**Token delivery:** access and refresh tokens are returned in both the JSON response body and as `HttpOnly` cookies (`access_token`, `refresh_token`). Cookie config is driven by `COOKIE_*` env vars.

**Token extraction order:** `Authenticate()` middleware checks the `Authorization: Bearer <token>` header first, then falls back to the `access_token` cookie.

## Sessions

```sql
CREATE TABLE sessions (
    id            TEXT PRIMARY KEY,
    user_id       TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL,
    user_agent    TEXT,
    ip            TEXT,
    expires_at    TIMESTAMPTZ NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Logout does **not** require authentication — the client sends the refresh token in the request body and the matching session row is deleted.

## MFA

TOTP-based MFA (RFC 6238). Flow:

1. User enables MFA → server generates TOTP secret, returns it with a `otpauth://` URI for QR code rendering.
2. User verifies with first TOTP code → MFA activated (`is_enabled = true` in `mfa_configs`).
3. On login: server issues a short-lived MFA temp token; user submits TOTP code; server exchanges it for full access + refresh tokens.
4. `mfa_configs` stores method, secret, and enabled flag (one row per user).

### MFA Recovery Codes

Ten single-use backup codes are generated at MFA activation. Each code is stored **hashed** in `mfa_recovery_codes`:

```sql
CREATE TABLE mfa_recovery_codes (
    id         TEXT PRIMARY KEY,
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash  TEXT NOT NULL,
    used_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

A recovery code can substitute for a TOTP code during login. After use, `used_at` is stamped and the code cannot be reused.

## Account States

| State              | Column                  | Reversible  | Effect                     |
| ------------------ | ----------------------- | ----------- | -------------------------- |
| Active             | —                       | —           | Normal access              |
| Soft-deactivated   | `deactivated_at`        | Yes (admin) | Login rejected             |
| Permanently locked | `locked_permanently_at` | No (via UI) | Login rejected permanently |

`Authenticate()` middleware checks both `deactivated_at` and `locked_permanently_at` on every request and returns `403` if either is set.

## Verification Tokens

Email verification and password reset use single-use, time-limited tokens:

```sql
CREATE TABLE verification_tokens (
    id         TEXT PRIMARY KEY,
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT NOT NULL,
    type       TEXT NOT NULL,   -- 'email_verification' | 'password_reset'
    expires_at TIMESTAMPTZ NOT NULL,
    used_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```
