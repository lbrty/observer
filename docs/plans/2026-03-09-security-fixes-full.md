# Security Fixes — Full Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Resolve all HIGH/MEDIUM/LOW/architecture findings from the 2026-03-09 security review across the full codebase.

**Architecture:** Tasks are ordered by severity. Tasks 3–4 share a `SessionRepository` interface change (done together to avoid two mock-regen rounds). Task 5 (IDOR) is the broadest — commit per entity sub-task. All other tasks are narrow, independent one-file changes.

**Tech Stack:** Go 1.25, Gin, sqlx/postgres, Redis (go-redis), `crypto/rand`, `mime`, `net/http.DetectContentType`, React/TypeScript (Vite, ky).

---

## Task 1: Registration security — hardcode guest role and prevent email enumeration

**Findings fixed:** Role self-assignment (HIGH), email enumeration (HIGH)

**Files:**

- Modify: `internal/usecase/auth/types.go`
- Modify: `internal/usecase/auth/auth_usecase.go`

Both fixes touch the same two files. Do them in one pass.

**Step 1: Write the failing tests**

In `internal/usecase/auth/` (create `register_test.go` if no test file exists):

```go
func TestRegister_IgnoresRoleField(t *testing.T) {
    // mock: GetByEmail → not found; Create → nil; credCreate → nil
    out, err := uc.Register(ctx, RegisterInput{
        Email:    "attacker@example.com",
        Password: "Password1!",
    })
    require.NoError(t, err)
    // assert userRepo.Create was called with Role == user.RoleGuest
}

func TestRegister_DuplicateEmailReturnsSuccess(t *testing.T) {
    // mock: GetByEmail → returns existing user (no error)
    out, err := uc.Register(ctx, RegisterInput{Email: "existing@example.com", Password: "Password1!"})
    require.NoError(t, err, "duplicate email must NOT return an error to the caller")
    assert.NotEmpty(t, out.Message)
}
```

**Step 2: Run to confirm they fail**

```bash
just test ./internal/usecase/auth/...
```

Expected: compile error (no `Role` field yet removed) or FAIL.

**Step 3: Remove `Role` from `RegisterInput`**

In `internal/usecase/auth/types.go`, delete the `Role` field:

```go
// RegisterInput holds data for user registration.
type RegisterInput struct {
    Email    string `json:"email"     binding:"required,email"`
    Password string `json:"password"  binding:"required,min=8"`
}
```

**Step 4: Fix `Register` in `auth_usecase.go`**

Remove the `ValidateRole` call. Hardcode `RoleGuest`. Silently return success on duplicate email:

```go
func (uc *AuthUseCase) Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
    // Email already registered — return success silently to prevent enumeration.
    if _, err := uc.userRepo.GetByEmail(ctx, input.Email); err == nil {
        return &RegisterOutput{
            Message: "Registration successful. Your account is pending admin approval.",
        }, nil
    }

    hash, salt, err := uc.hasher.Hash(input.Password)
    if err != nil {
        return nil, fmt.Errorf("hash password: %w", err)
    }

    userID := iulid.New()
    now := time.Now().UTC()

    newUser := &user.User{
        ID:         userID,
        Email:      input.Email,
        Role:       user.RoleGuest, // always guest; admin promotes via /admin/users/:id
        IsVerified: false,
        IsActive:   false,
        CreatedAt:  now,
        UpdatedAt:  now,
    }
    // ... rest unchanged (Create user, Create credentials, return RegisterOutput)
```

**Step 5: Run tests**

```bash
just test
```

Expected: all pass.

**Step 6: Commit**

```bash
git add internal/usecase/auth/types.go internal/usecase/auth/auth_usecase.go
git commit -m "Force registration role to guest and prevent email enumeration"
```

---

## Task 2: Enforce account status in VerifyMFA and RefreshToken

**Findings fixed:** Deactivated user can complete MFA flow (HIGH), deactivated user can rotate tokens (HIGH)

**Files:**

- Modify: `internal/usecase/auth/auth_usecase.go`

**Step 1: Write failing tests**

```go
func TestVerifyMFA_BlocksDeactivatedUser(t *testing.T) {
    // mock: ValidateMFAToken → valid claims
    // mock: mfaRepo.GetByUserID → enabled MFA
    // mock: totp.Validate → true
    // mock: userRepo.GetByID → user with DeactivatedAt != nil
    _, err := uc.VerifyMFA(ctx, VerifyMFAInput{MFAToken: "tok", TOTPCode: "123456"}, "ua", "127.0.0.1")
    require.Error(t, err)
}

func TestRefreshToken_BlocksDeactivatedUser(t *testing.T) {
    // mock: sessionRepo.GetByRefreshToken → valid non-expired session
    // mock: userRepo.GetByID → deactivated user
    _, err := uc.RefreshToken(ctx, RefreshTokenInput{RefreshToken: "tok"})
    require.Error(t, err)
}
```

**Step 2: Run to confirm they fail**

```bash
just test ./internal/usecase/auth/...
```

**Step 3: Fix `VerifyMFA`**

After `uc.userRepo.GetByID(ctx, userID)` (look for the user fetch inside `VerifyMFA`), add:

```go
u, err := uc.userRepo.GetByID(ctx, userID)
if err != nil {
    return nil, err
}
if err := u.CanLogin(); err != nil {
    return nil, err
}
```

**Step 4: Fix `RefreshToken`**

After `uc.userRepo.GetByID(ctx, session.UserID)` inside `RefreshToken`, add:

```go
u, err := uc.userRepo.GetByID(ctx, session.UserID)
if err != nil {
    return nil, fmt.Errorf("get user: %w", err)
}
if err := u.CanLogin(); err != nil {
    return nil, err
}
```

**Step 5: Run tests and commit**

```bash
just test
git add internal/usecase/auth/auth_usecase.go
git commit -m "Enforce CanLogin in VerifyMFA and RefreshToken to block deactivated accounts"
```

---

## Task 3: Extend SessionRepository — DeleteByUserID and DeleteExpired

**Findings fixed:** Sessions survive password change/reset (HIGH), expired sessions never vacuumed (LOW)

**Files:**

- Modify: `internal/repository/interfaces.go`
- Modify: `internal/repository/session_repository.go` (or wherever sessions are implemented)
- Run: `just generate-mocks`
- Create: `internal/app/session_vacuum.go`

Doing both interface additions here to avoid running mock-regen twice.

**Step 1: Add both methods to `SessionRepository`**

In `internal/repository/interfaces.go`:

```go
type SessionRepository interface {
    Create(ctx context.Context, session *auth.Session) error
    GetByRefreshToken(ctx context.Context, token string) (*auth.Session, error)
    Delete(ctx context.Context, id ulid.ULID) error
    DeleteByRefreshToken(ctx context.Context, token string) error
    // DeleteByUserID removes all sessions for a user (call after password change/reset).
    DeleteByUserID(ctx context.Context, userID ulid.ULID) error
    // DeleteExpired removes sessions whose expires_at is in the past.
    DeleteExpired(ctx context.Context) (int64, error)
}
```

**Step 2: Implement both in session_repository.go**

```go
func (r *sessionRepo) DeleteByUserID(ctx context.Context, userID ulid.ULID) error {
    const q = `DELETE FROM sessions WHERE user_id = $1`
    _, err := r.db.ExecContext(ctx, q, userID.String())
    if err != nil {
        return fmt.Errorf("delete sessions by user: %w", err)
    }
    return nil
}

func (r *sessionRepo) DeleteExpired(ctx context.Context) (int64, error) {
    result, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE expires_at < NOW()`)
    if err != nil {
        return 0, err
    }
    n, _ := result.RowsAffected()
    return n, nil
}
```

**Step 3: Regenerate mocks**

```bash
just generate-mocks
```

**Step 4: Create the vacuum goroutine**

Create `internal/app/session_vacuum.go`:

```go
package app

import (
    "context"
    "log/slog"
    "time"

    "github.com/lbrty/observer/internal/repository"
)

// StartSessionVacuum deletes expired sessions every 24 hours.
// Call with a cancellable context; cancel it to stop.
func StartSessionVacuum(ctx context.Context, repo repository.SessionRepository) {
    go func() {
        ticker := time.NewTicker(24 * time.Hour)
        defer ticker.Stop()
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                n, err := repo.DeleteExpired(ctx)
                if err != nil {
                    slog.Error("session vacuum", slog.Any("err", err))
                } else {
                    slog.Info("session vacuum", slog.Int64("deleted", n))
                }
            }
        }
    }()
}
```

**Step 5: Run tests**

```bash
just test
```

**Step 6: Commit**

```bash
git add internal/repository/interfaces.go \
        internal/repository/session_repository.go \
        internal/repository/mock/repository.go \
        internal/app/session_vacuum.go
git commit -m "Add DeleteByUserID and DeleteExpired to SessionRepository; add 24h session vacuum"
```

---

## Task 4: Invalidate sessions on password change and admin reset

**Findings fixed:** Active sessions survive credential rotation (HIGH)

**Files:**

- Modify: `internal/usecase/auth/auth_usecase.go`
- Modify: `internal/usecase/admin/user_usecase.go`
- Modify: `internal/app/container.go`

Depends on Task 3 (`DeleteByUserID` now exists on the interface and mock).

**Step 1: Write failing tests**

```go
func TestChangePassword_InvalidatesSessions(t *testing.T) {
    // mock credRepo.GetByUserID, hasher.Verify, hasher.Hash, credRepo.Update — all succeed
    // assert sessionRepo.DeleteByUserID called once with testUserID
    err := uc.ChangePassword(ctx, testUserID, ChangePasswordInput{
        CurrentPassword: "old",
        NewPassword:     "newPassword1!",
    })
    require.NoError(t, err)
}

func TestResetPassword_InvalidatesSessions(t *testing.T) {
    // admin use case mock: userRepo.GetByID succeeds, credRepo.GetByUserID + Update succeed
    // assert sessionRepo.DeleteByUserID called once
    err := uc.ResetPassword(ctx, targetUserID, ResetPasswordInput{NewPassword: "newPassword1!"})
    require.NoError(t, err)
}
```

**Step 2: Run to confirm they fail**

```bash
just test ./internal/usecase/...
```

**Step 3: Fix `ChangePassword` in `auth_usecase.go`**

After `uc.credRepo.Update(ctx, cred)` succeeds:

```go
if err := uc.credRepo.Update(ctx, cred); err != nil {
    return fmt.Errorf("update credentials: %w", err)
}
if err := uc.sessionRepo.DeleteByUserID(ctx, userID); err != nil {
    return fmt.Errorf("invalidate sessions: %w", err)
}
return nil
```

**Step 4: Fix `ResetPassword` in `user_usecase.go`**

`UserUseCase` likely does not hold `sessionRepo`. Add it:

```go
type UserUseCase struct {
    // ... existing fields
    sessionRepo repository.SessionRepository
}

// Update NewUserUseCase to accept and store sessionRepo.
```

Then after `uc.credRepo.Update(ctx, cred)` in `ResetPassword`:

```go
if err := uc.sessionRepo.DeleteByUserID(ctx, userID); err != nil {
    return fmt.Errorf("invalidate sessions: %w", err)
}
```

**Step 5: Update `container.go`**

Find `NewUserUseCase(...)` call in `internal/app/container.go` and add `container.SessionRepo`.

**Step 6: Run tests and commit**

```bash
just test
git add internal/usecase/auth/auth_usecase.go \
        internal/usecase/admin/user_usecase.go \
        internal/app/container.go
git commit -m "Invalidate all user sessions on password change and admin reset"
```

---

## Task 5: Fix cross-project IDOR on all project-scoped entities

**Findings fixed:** Pervasive IDOR — any project member can read/write records from other projects (HIGH)

**Strategy:** After every `GetByID`, compare the returned entity's `ProjectID` (or `PersonID` for notes and migration records) against the URL param. Return the domain "not found" error on mismatch — avoids leaking the existence of foreign records. Work one entity at a time; commit after each.

---

**5a — Person** (`internal/usecase/project/person_usecase.go`, `internal/handler/person_handler.go`)

Add `projectID string` as first parameter to `Get`, `Update`, and `Delete`. Check `p.ProjectID != projectID`.

```go
func (uc *PersonUseCase) Get(ctx context.Context, projectID, id string, canViewContact, canViewPersonal bool) (*PersonDTO, error) {
    p, err := uc.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get person: %w", err)
    }
    if p.ProjectID != projectID {
        return nil, person.ErrPersonNotFound
    }
    dto := personToDTO(p, canViewContact, canViewPersonal)
    return &dto, nil
}
// Apply same pattern to Update and Delete
```

In `person_handler.go`, pass `c.Param("project_id")` as first argument to `Get` and `Update`.

Write test:

```go
func TestPersonUseCase_Get_WrongProject_ReturnsNotFound(t *testing.T) {
    p := &person.Person{ID: "01ABC", ProjectID: "project-B"}
    mockRepo.EXPECT().GetByID(ctx, "01ABC").Return(p, nil)
    _, err := uc.Get(ctx, "project-A", "01ABC", true, true)
    assert.ErrorIs(t, err, person.ErrPersonNotFound)
}
```

```bash
just test
git add internal/usecase/project/person_usecase.go internal/handler/person_handler.go
git commit -m "Enforce project ownership check in PersonUseCase Get/Update/Delete"
```

---

**5b — SupportRecord** (`support_record_usecase.go`, `support_record_handler.go`)

Same pattern: add `projectID string` to `Get`, `Update`, `Delete`. Check `r.ProjectID != projectID`. Return `support.ErrRecordNotFound`.

```bash
git add internal/usecase/project/support_record_usecase.go internal/handler/support_record_handler.go
git commit -m "Enforce project ownership check in SupportRecordUseCase Get/Update/Delete"
```

---

**5c — Household** (`household_usecase.go`, `household_handler.go`)

Same pattern. Check `h.ProjectID != projectID`. Return `household.ErrHouseholdNotFound`.

```bash
git add internal/usecase/project/household_usecase.go internal/handler/household_handler.go
git commit -m "Enforce project ownership check in HouseholdUseCase Get/Update/Delete"
```

---

**5d — Pet** (`pet_usecase.go`, `pet_handler.go`)

Same pattern. Check `p.ProjectID != projectID`. Return `pet.ErrPetNotFound`.

```bash
git add internal/usecase/project/pet_usecase.go internal/handler/pet_handler.go
git commit -m "Enforce project ownership check in PetUseCase Get/Update/Delete"
```

---

**5e — Document** (`document_usecase.go`, `document_handler.go`)

Apply to `Get`, `Download`, `Thumbnail`, `Update`, `Delete`. Check `d.ProjectID != projectID`. Return `document.ErrDocumentNotFound`.

```bash
git add internal/usecase/project/document_usecase.go internal/handler/document_handler.go
git commit -m "Enforce project ownership check in DocumentUseCase all read/write operations"
```

---

**5f — Note** (`note_usecase.go`, `note_handler.go`)

Notes have `PersonID`, not `ProjectID`. Add `personID string` param to `Get`, `Update`, `Delete`. Check `n.PersonID != personID`. Return `note.ErrNoteNotFound`. Handler passes `c.Param("person_id")`.

```go
func (uc *NoteUseCase) Get(ctx context.Context, personID, id string) (*NoteDTO, error) {
    n, err := uc.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if n.PersonID != personID {
        return nil, note.ErrNoteNotFound
    }
    dto := noteToDTO(n)
    return &dto, nil
}
```

```bash
git add internal/usecase/project/note_usecase.go internal/handler/note_handler.go
git commit -m "Verify note belongs to expected person in Get/Update/Delete"
```

---

**5g — MigrationRecord** (`migration_record_usecase.go`, `migration_record_handler.go`)

Same as notes: add `personID string`, check `r.PersonID != personID`, return `migration.ErrRecordNotFound`.

```bash
git add internal/usecase/project/migration_record_usecase.go internal/handler/migration_record_handler.go
git commit -m "Verify migration record belongs to expected person in Get/Update"
```

---

**5h — Permission** (`internal/usecase/admin/permission_usecase.go`, `internal/handler/permission_handler.go`)

Add `projectID string` to `Update`. Check `perm.ProjectID != projectID`. Return `project.ErrPermissionNotFound`. Handler passes `c.Param("project_id")`.

```go
func (uc *PermissionUseCase) Update(ctx context.Context, projectID, id string, input UpdatePermissionInput) (*PermissionDTO, error) {
    perm, err := uc.permRepo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get permission for update: %w", err)
    }
    if perm.ProjectID != projectID {
        return nil, project.ErrPermissionNotFound
    }
    // ... rest unchanged
```

```bash
git add internal/usecase/admin/permission_usecase.go internal/handler/permission_handler.go
git commit -m "Scope permission update to the project in URL; reject cross-project mutations"
```

---

## Task 6: Document upload and serving security

**Findings fixed:** Path traversal in filename (MEDIUM), stored XSS via declared MIME (MEDIUM), Content-Disposition header injection (MEDIUM), stream endpoint deletes X-Frame-Options (MEDIUM)

**Files:**

- Modify: `internal/handler/document_handler.go`

All four issues are in the same handler; fix them together.

**Step 1: Write failing tests**

```go
func TestUpload_RejectsForbiddenMIME(t *testing.T) {
    // Send multipart with file content that DetectContentType returns "text/html; charset=utf-8"
    // Expect 400
}

func TestUpload_SanitizesFilename(t *testing.T) {
    // header.Filename = "../../etc/passwd"
    // filepath.Base should reduce it to "passwd" — use case called with "passwd", not the traversal path
}
```

**Step 2: Fix path traversal — strip to base filename**

In `Upload`, immediately after `file, header, err := c.Request.FormFile("file")`:

```go
import "path/filepath"

filename := filepath.Base(header.Filename)
if filename == "." || filename == "" {
    c.JSON(http.StatusBadRequest, errJSON("errors.validation", "invalid filename"))
    return
}
// Use `filename` everywhere below, NOT `header.Filename`
```

**Step 3: Detect MIME server-side and reject dangerous types**

```go
import (
    "bytes"
    "io"
)

buf := make([]byte, 512)
n, _ := file.Read(buf)
detectedMIME := http.DetectContentType(buf[:n])

// Reject types that execute in browser
switch detectedMIME {
case "text/html; charset=utf-8", "application/xhtml+xml":
    c.JSON(http.StatusBadRequest, errJSON("errors.validation", "file type not permitted"))
    return
}

body := io.MultiReader(bytes.NewReader(buf[:n]), file)

out, err := h.uc.Upload(ctx, projectID, personID, userID.String(), filename, detectedMIME, header.Size, body)
```

**Step 4: Fix Content-Disposition header injection**

In `Download` and `Stream`, replace the string concat:

```go
import "mime"

// Download
c.Header("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": doc.Name}))

// Stream
c.Header("Content-Disposition", mime.FormatMediaType("inline", map[string]string{"filename": doc.Name}))
```

`mime.FormatMediaType` handles quoting and RFC 5987 encoding, so filenames with `"` or `\r\n` cannot inject headers.

**Step 5: Fix stream endpoint — restore framing protection and restrict inline MIME**

The stream handler currently deletes `X-Frame-Options` entirely. Replace with:

```go
// Allow framing only from same origin (needed for in-app viewer)
c.Header("Content-Security-Policy", "frame-ancestors 'self'")
c.Header("X-Frame-Options", "SAMEORIGIN")

// Only stream known-safe MIME types inline; force download for everything else
safeMIME := doc.MimeType
switch {
case strings.HasPrefix(safeMIME, "image/"),
     strings.HasPrefix(safeMIME, "video/"),
     safeMIME == "application/pdf":
    // keep as-is
default:
    safeMIME = "application/octet-stream"
}
c.DataFromReader(http.StatusOK, doc.Size, safeMIME, rc, nil)
```

Add `"strings"` import if not already present.

**Step 6: Run tests and commit**

```bash
just test
git add internal/handler/document_handler.go
git commit -m "Fix document upload path traversal, MIME detection, Content-Disposition injection, stream XSS"
```

---

## Task 7: Global JSON body size limit (exempt multipart)

**Findings fixed:** No server-wide body cap — attackers can POST multi-GB JSON (MEDIUM)

**Files:**

- Modify: `internal/server/server.go:74–89` (setupMiddleware)

**Why we exempt multipart:** `http.MaxBytesReader` wraps the underlying reader. If the global 1 MB wrapper is applied first, the upload handler's subsequent 50 MB wrap still reads from the 1 MB-limited inner reader and silently breaks file uploads. The upload handler (`document_handler.go:55`) already enforces its own 50 MB cap — no change needed there.

**Step 1: Write the failing tests**

In `internal/server/server_test.go`:

```go
func TestBodySizeLimit_JSONRejectedOver1MB(t *testing.T) {
    w := httptest.NewRecorder()
    bigBody := strings.Repeat("x", 2<<20) // 2 MB
    req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(bigBody))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
}

func TestBodySizeLimit_MultipartNotLimited(t *testing.T) {
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader("body"))
    req.Header.Set("Content-Type", "multipart/form-data; boundary=xxx")
    router.ServeHTTP(w, req)
    assert.NotEqual(t, http.StatusRequestEntityTooLarge, w.Code)
}
```

**Step 2: Run to confirm they fail**

```bash
just test -run TestBodySizeLimit ./internal/server/...
```

**Step 3: Add middleware at the top of setupMiddleware**

```go
func (s *Server) setupMiddleware(cfg *config.Config, log *slog.Logger) {
    // Limit non-multipart request bodies to 1 MB to prevent memory exhaustion.
    // Multipart (file uploads) are excluded — document_handler enforces 50 MB itself.
    s.router.Use(func(c *gin.Context) {
        if !strings.HasPrefix(c.GetHeader("Content-Type"), "multipart/form-data") {
            c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
        }
        c.Next()
    })
    s.router.Use(requestIDMiddleware())
    // ... rest unchanged
```

Add `"strings"` and `"net/http"` to imports if not already present.

**Step 4: Run tests and commit**

```bash
just test
git add internal/server/server.go internal/server/server_test.go
git commit -m "Add global 1 MB body size limit, exempt multipart uploads"
```

---

## Task 8: CSRF protection via double-submit cookie

**Findings fixed:** CSRF gap for cookie-based flows (MEDIUM)

**Files:**

- Create: `internal/middleware/csrf.go`
- Create: `internal/middleware/csrf_test.go`
- Modify: `internal/handler/auth_handler.go` (setTokenCookies, clearTokenCookies)
- Modify: `internal/server/server.go` (wire middleware)

**Step 1: Create `internal/middleware/csrf.go`**

```go
package middleware

import (
    "crypto/subtle"
    "net/http"

    "github.com/gin-gonic/gin"
)

const (
    CSRFTokenHeader = "X-CSRF-Token"
    CSRFTokenCookie = "csrf_token"
)

var csrfSafeMethods = map[string]bool{
    http.MethodGet:     true,
    http.MethodHead:    true,
    http.MethodOptions: true,
}

// CSRFProtection validates state-changing requests carry an X-CSRF-Token header
// matching the csrf_token cookie (double-submit cookie pattern).
// The csrf_token cookie is set at login (HttpOnly: false) so JS can read it.
// A cross-site attacker can force cookie sending but cannot read cookies or set custom headers.
func CSRFProtection() gin.HandlerFunc {
    return func(c *gin.Context) {
        if csrfSafeMethods[c.Request.Method] {
            c.Next()
            return
        }

        cookie, err := c.Cookie(CSRFTokenCookie)
        header := c.GetHeader(CSRFTokenHeader)

        if err != nil || cookie == "" || header == "" ||
            subtle.ConstantTimeCompare([]byte(cookie), []byte(header)) != 1 {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": "invalid or missing CSRF token",
                "code":  "errors.auth.invalidCSRFToken",
            })
            return
        }

        c.Next()
    }
}
```

**Step 2: Create `internal/middleware/csrf_test.go`**

```go
package middleware_test

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"

    "github.com/lbrty/observer/internal/middleware"
)

func setupCSRFRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.Use(middleware.CSRFProtection())
    r.POST("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
    return r
}

func TestCSRF_MissingCookieRejected(t *testing.T) {
    r := setupCSRFRouter()
    req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}"))
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCSRF_MissingHeaderRejected(t *testing.T) {
    r := setupCSRFRouter()
    req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}"))
    req.AddCookie(&http.Cookie{Name: "csrf_token", Value: "abc123"})
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCSRF_MismatchedTokenRejected(t *testing.T) {
    r := setupCSRFRouter()
    req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}"))
    req.AddCookie(&http.Cookie{Name: "csrf_token", Value: "abc123"})
    req.Header.Set("X-CSRF-Token", "wrong")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCSRF_MatchingTokenAllowed(t *testing.T) {
    r := setupCSRFRouter()
    req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}"))
    req.AddCookie(&http.Cookie{Name: "csrf_token", Value: "abc123"})
    req.Header.Set("X-CSRF-Token", "abc123")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
}

func TestCSRF_GETAllowedWithoutToken(t *testing.T) {
    r := setupCSRFRouter()
    r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
    req := httptest.NewRequest(http.MethodGet, "/test", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
}
```

**Step 3: Run to confirm tests fail (middleware not wired yet)**

```bash
just test -run TestCSRF ./internal/middleware/...
```

**Step 4: Issue CSRF cookie alongside auth cookies**

In `internal/handler/auth_handler.go`, update `setTokenCookies` to issue the CSRF cookie, and `clearTokenCookies` to clear it:

```go
import (
    "crypto/rand"
    "encoding/hex"
    // existing imports...
)

func (h *AuthHandler) setTokenCookies(c *gin.Context, accessToken, refreshToken string) {
    sameSite := h.cookie.HTTPSameSite()

    http.SetCookie(c.Writer, &http.Cookie{
        Name:     accessTokenCookie,
        Value:    accessToken,
        Path:     "/",
        Domain:   h.cookie.Domain,
        MaxAge:   int(h.cookie.MaxAge.Seconds()),
        HttpOnly: true,
        Secure:   h.cookie.Secure,
        SameSite: sameSite,
    })

    http.SetCookie(c.Writer, &http.Cookie{
        Name:     refreshTokenCookie,
        Value:    refreshToken,
        Path:     "/auth",
        Domain:   h.cookie.Domain,
        MaxAge:   int(h.cookie.MaxAge.Seconds()),
        HttpOnly: true,
        Secure:   h.cookie.Secure,
        SameSite: sameSite,
    })

    // CSRF double-submit cookie — HttpOnly: false so JS can read it.
    http.SetCookie(c.Writer, &http.Cookie{
        Name:     middleware.CSRFTokenCookie,
        Value:    generateCSRFToken(),
        Path:     "/",
        Domain:   h.cookie.Domain,
        MaxAge:   int(h.cookie.MaxAge.Seconds()),
        HttpOnly: false,
        Secure:   h.cookie.Secure,
        SameSite: sameSite,
    })
}

func generateCSRFToken() string {
    b := make([]byte, 16)
    if _, err := rand.Read(b); err != nil {
        return "fallback-csrf"
    }
    return hex.EncodeToString(b)
}

func (h *AuthHandler) clearTokenCookies(c *gin.Context) {
    sameSite := h.cookie.HTTPSameSite()
    // ... existing access_token and refresh_token clears unchanged ...

    http.SetCookie(c.Writer, &http.Cookie{
        Name:     middleware.CSRFTokenCookie,
        Value:    "",
        Path:     "/",
        Domain:   h.cookie.Domain,
        MaxAge:   -1,
        HttpOnly: false,
        Secure:   h.cookie.Secure,
        SameSite: sameSite,
    })
}
```

**Step 5: Wire CSRF middleware in `server.go`**

Add `middleware.CSRFProtection()` at the end of `setupMiddleware`, after `SecurityHeaders()`:

```go
s.router.Use(middleware.SecurityHeaders())
s.router.Use(middleware.CSRFProtection())
```

**Step 6: Frontend integration note**

In your ky instance (`src/lib/api.ts`):

```ts
hooks: {
  beforeRequest: [
    (request) => {
      const csrf = document.cookie.match(/csrf_token=([^;]+)/)?.[1];
      if (csrf) request.headers.set("X-CSRF-Token", csrf);
    },
  ];
}
```

This works identically in local dev (localhost:5173 → localhost:9000) and production — no environment branching needed.

**Step 7: Run tests and commit**

```bash
just test
git add internal/middleware/csrf.go internal/middleware/csrf_test.go \
        internal/handler/auth_handler.go internal/server/server.go
git commit -m "Add CSRF protection via double-submit cookie pattern"
```

---

## Task 9: Clamp pagination before repository call

**Findings fixed:** Unbounded `per_page` materializes full table (MEDIUM)

**Files:**

- Modify: `internal/usecase/project/person_usecase.go`
- Modify: `internal/usecase/project/support_record_usecase.go`
- Check all other list use cases for the same pattern

Current wrong order: filter built with raw input → repo queried → clamp called (too late).
Correct order: clamp first → use clamped values in filter.

**Step 1: Write test**

```go
func TestPersonList_ClampsPerPage(t *testing.T) {
    mockRepo.EXPECT().List(ctx, gomock.MatchedBy(func(f person.PersonListFilter) bool {
        return f.PerPage <= 100
    })).Return([]*person.Person{}, 0, nil)

    _, err := uc.List(ctx, "proj-1", ListPeopleInput{PerPage: 999999}, true, true)
    require.NoError(t, err)
}
```

**Step 2: Run to confirm it fails**

```bash
just test -run TestPersonList_ClampsPerPage ./internal/usecase/project/...
```

**Step 3: Move `ClampPagination` before filter construction**

```go
func (uc *PersonUseCase) List(...) {
    page, perPage := usecase.ClampPagination(input.Page, input.PerPage) // FIRST

    filter := person.PersonListFilter{
        ProjectID: projectID,
        Page:      page,    // clamped
        PerPage:   perPage, // clamped
        // ... other fields
    }
    // ... rest unchanged; remove any ClampPagination call that came after the repo query
```

Apply the same fix to `SupportRecordUseCase.List` and any other list use case where `ClampPagination` comes after the repo call.

**Step 4: Run tests and commit**

```bash
just test
git add internal/usecase/project/person_usecase.go internal/usecase/project/support_record_usecase.go
git commit -m "Clamp pagination before repository call to prevent unbounded result sets"
```

---

## Task 10: Suppress raw error details from 500 responses

**Findings fixed:** Internal DB error messages leaked via API (MEDIUM)

**Files:**

- Modify: `internal/handler/errors.go`

Current: `c.JSON(status, errJSON(code, err.Error()))` for all errors including 500 — exposes raw postgres/driver error strings.

**Step 1: Fix the 500 branch in `HandleError`**

```go
func HandleError(c *gin.Context, err error) {
    status, code := MapDomainError(err)
    if status == http.StatusInternalServerError {
        slog.ErrorContext(c.Request.Context(), "unexpected error",
            "error", err,
            "path", c.Request.URL.Path,
        )
        c.JSON(status, errJSON(code, "internal server error"))
        return
    }
    c.JSON(status, errJSON(code, err.Error()))
}
```

Domain errors (not-found, conflict, validation) still return descriptive messages — those are intentionally user-facing. Only the unexpected 500 path changes.

**Step 2: Run tests and commit**

```bash
just test
git add internal/handler/errors.go
git commit -m "Return generic message for 500 errors; log actual error server-side"
```

---

## Task 11: Remove staff implicit can_export override

**Findings fixed:** `RoleStaff` silently overrides explicit `can_export=false` (MEDIUM)

**Files:**

- Modify: `internal/middleware/project_auth.go`

**Step 1: Find and remove the override**

Find the line:

```go
canExport := perm.CanExport || userRole == user.RoleStaff
```

Replace with:

```go
canExport := perm.CanExport
```

Staff users who need export access must have `can_export = true` set explicitly, like all other non-owner users.

**Step 2: Run tests and commit**

```bash
just test
git add internal/middleware/project_auth.go
git commit -m "Honour explicit can_export flag for staff; remove implicit export grant"
```

---

## Task 12: Persist permanent lockouts to PostgreSQL

**Findings fixed:** Redis restart clears all lockouts including permanent ones (MEDIUM)

**Files:**

- Create: `migrations/000034_add_user_permanent_lock.up.sql`
- Modify: `internal/domain/user/entity.go`
- Modify: `internal/repository/interfaces.go`
- Modify: `internal/repository/user_repository.go`
- Modify: `internal/repository/login_attempt_store.go`
- Modify: `internal/middleware/auth.go`
- Modify: `internal/app/container.go`
- Run: `just generate-mocks`

**Step 1: Create the migration**

```sql
-- migrations/000034_add_user_permanent_lock.up.sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS locked_permanently_at TIMESTAMPTZ;
```

```bash
just migrate
```

**Step 2: Add field to domain entity**

In `internal/domain/user/entity.go`, add:

```go
LockedPermanentlyAt *time.Time `db:"locked_permanently_at"`
```

**Step 3: Add `LockPermanently` to `UserRepository` interface**

In `internal/repository/interfaces.go`:

```go
// LockPermanently sets locked_permanently_at for the user identified by email.
LockPermanently(ctx context.Context, email string) error
```

**Step 4: Implement in user_repository.go**

```go
func (r *userRepository) LockPermanently(ctx context.Context, email string) error {
    _, err := r.db.ExecContext(ctx,
        `UPDATE users SET locked_permanently_at = NOW() WHERE email = $1`,
        email,
    )
    return err
}
```

**Step 5: Inject UserRepository into LoginAttemptStore**

In `internal/repository/login_attempt_store.go`:

```go
type redisLoginAttemptStore struct {
    client   *redis.Client
    userRepo UserRepository
}

func NewLoginAttemptStore(client *redis.Client, userRepo UserRepository) LoginAttemptStore {
    return &redisLoginAttemptStore{client: client, userRepo: userRepo}
}
```

In `RecordFailure`, after the permanent lock Redis write:

```go
if lockDuration == 0 {
    if err := s.client.Set(ctx, lockKey, "permanent", 0).Err(); err != nil {
        return 0, fmt.Errorf("set permanent lock: %w", err)
    }
    // Also persist to PostgreSQL so lockout survives Redis restarts.
    if err := s.userRepo.LockPermanently(ctx, email); err != nil {
        slog.Error("persist permanent lock to db", slog.Any("err", err))
    }
    return -1, nil
}
```

**Step 6: Check `LockedPermanentlyAt` in auth middleware**

In `internal/middleware/auth.go`, update the deactivation check:

```go
if err != nil || u.DeactivatedAt != nil || u.LockedPermanentlyAt != nil {
    c.JSON(http.StatusForbidden, gin.H{
        "error": "account locked or deactivated",
        "code":  "errors.auth.accountLocked",
    })
    c.Abort()
    return
}
```

**Step 7: Update container.go**

Pass `container.UserRepo` as second argument to `NewLoginAttemptStore`:

```go
container.LoginAttemptStore = repository.NewLoginAttemptStore(redisClient, container.UserRepo)
```

**Step 8: Regenerate mocks**

```bash
just generate-mocks
```

**Step 9: Run tests and commit**

```bash
just test
git add migrations/000034_add_user_permanent_lock.up.sql \
        internal/domain/user/entity.go \
        internal/repository/interfaces.go \
        internal/repository/user_repository.go \
        internal/repository/login_attempt_store.go \
        internal/middleware/auth.go \
        internal/app/container.go \
        internal/repository/mock/repository.go
git commit -m "Persist permanent account lockouts to PostgreSQL for Redis-restart resilience"
```

---

## Task 13: Remove lockout duration from error response

**Findings fixed:** Exact remaining lockout duration exposed to attacker (LOW)

**Files:**

- Modify: `internal/handler/auth_handler.go`

**Step 1: Write the failing test**

```go
func TestLogin_LockedAccount_DoesNotLeakDuration(t *testing.T) {
    // mock IsLocked returning 15 * time.Minute
    // call Login
    var resp map[string]string
    json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NotContains(t, resp["error"], "15m")
    assert.NotContains(t, resp["error"], "for ")
}
```

**Step 2: Simplify both lockout response blocks**

Replace both `fmt.Sprintf("account locked for %s", remaining.Truncate(time.Second))` usages:

```go
// For both the IsLocked check and the RecordFailure check:
if remaining != 0 {
    msg := "account locked, contact administrator"
    if remaining > 0 {
        msg = "account temporarily locked, please try again later"
    }
    c.JSON(http.StatusTooManyRequests, errJSON("errors.auth.accountLocked", msg))
    return
}
```

Also remove the `"fmt"` and `"time"` imports if no longer used after this change.

**Step 3: Run tests and commit**

```bash
just test
git add internal/handler/auth_handler.go
git commit -m "Remove lockout duration from login error response"
```

---

## Task 14: Increase Argon2id time cost

**Findings fixed:** `time=1` at OWASP minimum (MEDIUM)

**Files:**

- Modify: `internal/crypto/password.go`

**⚠ Backward compatibility:** Existing stored hashes were produced with `time=1`. After this change, `Verify` uses `time=2` — all existing users will fail login. Handle this one of two ways:

- **Option A (clean):** Force all users to reset their password at next login (detect the old hash via a version field or re-hash on successful login with old params before updating).
- **Option B (simple):** Accept the breakage in a controlled migration: inform admins, use `admin.ResetPassword` to issue new passwords for active users.

Agree with the team on the approach before merging this task.

**Step 1: Update the time parameter**

```go
func NewArgonHasher() *ArgonHasher {
    return &ArgonHasher{
        time:    2,          // was 1; OWASP recommends t≥2 at 64 MB memory
        memory:  64 * 1024,
        threads: 4,
        keyLen:  32,
    }
}
```

**Step 2: Run tests and commit**

```bash
just test
git add internal/crypto/password.go
git commit -m "Increase Argon2id time cost to 2 per OWASP recommendation"
```

---

## Task 15: Add Content-Security-Policy header

**Findings fixed:** No CSP — stored XSS executes unrestricted (LOW)

**Files:**

- Modify: `internal/middleware/security_headers.go`

**Step 1: Write the failing test**

```go
func TestSecurityHeaders_CSP(t *testing.T) {
    // ... setup gin with SecurityHeaders middleware ...
    csp := w.Header().Get("Content-Security-Policy")
    assert.NotEmpty(t, csp)
    assert.Contains(t, csp, "default-src 'self'")
}
```

**Step 2: Add the header**

```go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-XSS-Protection", "0")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
        c.Header("Content-Security-Policy",
            "default-src 'self'; "+
                "script-src 'self'; "+
                "style-src 'self' 'unsafe-inline'; "+ // SPA may use inline styles
                "img-src 'self' data:; "+
                "connect-src 'self'; "+
                "font-src 'self'; "+
                "object-src 'none'; "+
                "frame-ancestors 'none'")
        c.Next()
    }
}
```

**Step 3: Run tests and commit**

```bash
just test
git add internal/middleware/security_headers.go
git commit -m "Add Content-Security-Policy header"
```

---

## Task 16: Wire session vacuum to server startup

**Findings fixed:** Companion to Task 3 — goroutine exists but is not started (LOW)

**Files:**

- Modify: `cmd/observer/cmd/serve.go` (or wherever `server.Start()` is called)

**Step 1: Locate the serve entrypoint**

```bash
grep -r "\.Start()" cmd/
```

**Step 2: Start the vacuum with a cancellable context**

```go
vacuumCtx, vacuumCancel := context.WithCancel(context.Background())
defer vacuumCancel()
app.StartSessionVacuum(vacuumCtx, container.SessionRepo)
```

Add this before `srv.Start()`.

**Step 3: Commit**

```bash
git add cmd/observer/cmd/serve.go
git commit -m "Start session vacuum goroutine on server startup"
```

---

## Task 17: Graceful shutdown on SIGTERM/SIGINT

**Findings fixed:** In-flight requests dropped on Kubernetes pod termination (LOW)

**Files:**

- Modify: `cmd/observer/cmd/serve.go`

**Step 1: Replace blocking `Start()` call**

```go
import (
    "os"
    "os/signal"
    "syscall"
)

go func() {
    if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
        slog.Error("server error", slog.Any("err", err))
        os.Exit(1)
    }
}()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

slog.Info("shutting down server...")
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil {
    slog.Error("server shutdown", slog.Any("err", err))
}
slog.Info("server stopped")
```

**Step 2: Commit**

```bash
git add cmd/observer/cmd/serve.go
git commit -m "Add graceful shutdown with SIGTERM/SIGINT handling"
```

---

## Task 18: Enforce HTTPS for API URL in frontend production builds

**Findings fixed:** `VITE_API_URL` defaults to `http://localhost` — could be used in production (MEDIUM)

**Files:**

- Modify: `packages/observer-web/vite.config.ts`
- Modify: `packages/observer-web/src/lib/api.ts`
- Modify: `packages/observer-web/src/hooks/use-documents.ts`

**Step 1: Add build-time guard in `vite.config.ts`**

```ts
import { defineConfig, loadEnv } from "vite";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");

  if (mode !== "development" && mode !== "test") {
    const apiUrl = env.VITE_API_URL ?? "";
    if (!apiUrl.startsWith("https://")) {
      throw new Error(`VITE_API_URL must start with https:// for ${mode} builds. Got: "${apiUrl}"`);
    }
  }

  return {
    // ... existing config unchanged
  };
});
```

**Step 2: Export `apiBase` from `src/lib/api.ts`**

```ts
export const apiBase = import.meta.env.VITE_API_URL ?? "http://localhost:9000";
```

**Step 3: Consolidate inline URL literals in `use-documents.ts`**

Replace all `${import.meta.env.VITE_API_URL ?? 'http://localhost:9000'}` occurrences:

```ts
import { apiBase } from "@/lib/api";

export const documentDownloadUrl = (projectId: string, docId: string) =>
  `${apiBase}/projects/${projectId}/documents/${docId}/download`;

export const documentStreamUrl = (projectId: string, docId: string) =>
  `${apiBase}/projects/${projectId}/documents/${docId}/stream`;

export const documentThumbnailUrl = (projectId: string, docId: string) =>
  `${apiBase}/projects/${projectId}/documents/${docId}/thumbnail`;
```

**Step 4: Verify the guard fires**

```bash
cd packages/observer-web
VITE_API_URL=http://insecure.example.com bunx vite build
# Expected: Error: VITE_API_URL must start with https://

VITE_API_URL=https://api.example.com bunx vite build
# Expected: build succeeds
```

**Step 5: Commit**

```bash
git add packages/observer-web/vite.config.ts \
        packages/observer-web/src/lib/api.ts \
        packages/observer-web/src/hooks/use-documents.ts
git commit -m "Enforce HTTPS for API URL in production builds; consolidate URL helpers"
```

---

## Task 19: Audit logging for admin mutations

**Findings fixed:** Admin mutations leave no audit trail (ARCH)

**Files:**

- Modify: `internal/usecase/admin/user_usecase.go`
- Modify: `internal/usecase/admin/permission_usecase.go`
- Modify: `internal/app/container.go`

**Step 1: Inject AuditUseCase into admin use cases**

```go
type UserUseCase struct {
    // existing fields...
    auditUC *audit.UseCase
}
```

**Step 2: Add `auditUC.Record` calls after each mutation**

Operations to cover:

```
admin.user.create    → CreateUser
admin.user.update    → UpdateUser
admin.user.deactivate → DeactivateUser
admin.user.reactivate → ReactivateUser
admin.user.reset_password → ResetPassword
admin.permission.assign   → AssignPermission
admin.permission.update   → UpdatePermission
admin.permission.revoke   → RevokePermission
```

Pattern (from any project use case):

```go
uc.auditUC.Record(ctx, nil, "admin.user.deactivate", "user", &targetUserIDStr,
    fmt.Sprintf("Deactivated user %s", targetUser.Email))
```

**Step 3: Update container.go**

Pass `container.AuditUC` when constructing admin use cases.

**Step 4: Run tests and commit**

```bash
just test
git add internal/usecase/admin/ internal/app/container.go
git commit -m "Add audit logging to all admin user and permission mutations"
```

---

## Task 20: Stronger password policy

**Findings fixed:** `min=8` accepts trivially weak passwords (ARCH)

**Files:**

- Modify: `internal/usecase/auth/types.go`
- Modify: `internal/app/container.go` (register custom validator)

**Step 1: Write the failing tests**

```go
func TestRegisterInput_WeakPasswordRejected(t *testing.T) {
    input := RegisterInput{Email: "a@b.com", Password: "aaaaaaaa"}
    err := validator.New().Struct(input)
    require.Error(t, err)
}

func TestRegisterInput_StrongPasswordAccepted(t *testing.T) {
    input := RegisterInput{Email: "a@b.com", Password: "CorrectHorse1!"}
    err := validator.New().Struct(input)
    require.NoError(t, err)
}
```

**Step 2: Register custom validator**

In `internal/app/container.go` (or server init), register with Gin's validator engine:

```go
import (
    "strings"
    "github.com/gin-gonic/gin/binding"
    "github.com/go-playground/validator/v10"
)

func registerCustomValidators() {
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("strongpassword", func(fl validator.FieldLevel) bool {
            p := fl.Field().String()
            hasDigit := strings.ContainsAny(p, "0123456789")
            hasSpecial := strings.ContainsAny(p, "!@#$%^&*()-_=+[]{}|;:',.<>?/`~")
            return hasDigit && hasSpecial
        })
    }
}
```

Call `registerCustomValidators()` once at startup.

**Step 3: Update binding tags**

In `internal/usecase/auth/types.go`:

```go
type RegisterInput struct {
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required,min=12,strongpassword"`
}
```

Apply the same `binding:"required,min=12,strongpassword"` to `ChangePasswordInput.NewPassword`.

**Step 4: Run tests and commit**

```bash
just test
git add internal/usecase/auth/types.go internal/app/container.go
git commit -m "Enforce stronger password policy: min 12 chars, digit, special character"
```

---

## Task 21: Log session IP/UA mismatch on token refresh

**Findings fixed:** Stolen refresh tokens used from new IP/UA go undetected (ARCH)

**Files:**

- Modify: `internal/usecase/auth/types.go`
- Modify: `internal/usecase/auth/auth_usecase.go`
- Modify: `internal/handler/auth_handler.go`

**Step 1: Add IP and UA to `RefreshTokenInput`**

```go
type RefreshTokenInput struct {
    RefreshToken string `json:"refresh_token"`
    UserAgent    string `json:"-"` // from request, not body
    IP           string `json:"-"`
}
```

**Step 2: Pass them from the handler**

In `internal/handler/auth_handler.go` `RefreshToken` handler:

```go
tokens, err := h.authUC.RefreshToken(c.Request.Context(), ucauth.RefreshTokenInput{
    RefreshToken: refreshToken,
    UserAgent:    c.GetHeader("User-Agent"),
    IP:           c.ClientIP(),
})
```

**Step 3: Log mismatches in the use case**

After loading the session in `RefreshToken`:

```go
if input.IP != "" && session.IP != input.IP {
    slog.Warn("token refresh from new IP",
        slog.String("session_ip", session.IP),
        slog.String("request_ip", input.IP),
        slog.String("user_id", session.UserID.String()),
    )
}
if input.UserAgent != "" && session.UserAgent != input.UserAgent {
    slog.Warn("token refresh from new user-agent",
        slog.String("user_id", session.UserID.String()),
    )
}
```

**Step 4: Run tests and commit**

```bash
just test
git add internal/usecase/auth/types.go \
        internal/usecase/auth/auth_usecase.go \
        internal/handler/auth_handler.go
git commit -m "Log session IP and User-Agent mismatches on token refresh"
```

---

## Task 22: CORS localhost warning on startup

**Findings fixed:** L-6 — CORS defaults to `localhost:5173` with no production warning

**Files:**

- Modify: `cmd/observer/cmd/serve.go` (or wherever config is loaded and server is started)

**Step 1: Add startup warning**

After `config.Load()`, add:

```go
import "strings"

for _, origin := range cfg.CORS.Origins {
    if strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
        slog.Warn("CORS_ORIGINS contains a localhost entry — ensure this is intentional in production",
            slog.String("origin", origin))
    }
}
```

**Step 2: Run and verify**

```bash
just build && ./observer serve
# Expected: WARN log line "CORS_ORIGINS contains a localhost entry..." when CORS_ORIGINS is unset
```

**Step 3: Commit**

```bash
git add cmd/observer/cmd/serve.go
git commit -m "Warn at startup when CORS_ORIGINS contains localhost"
```

---

## Task 23: Consolidate can_export into route middleware

**Findings fixed:** AC-2 — `can_export` manually checked in 4 handler spots; easy to miss on new endpoints

**Files:**

- Modify: `internal/middleware/project_auth.go`
- Modify: `internal/server/server.go`
- Modify: `internal/handler/export_handler.go`

**Step 1: Add `RequireExport` middleware**

In `internal/middleware/project_auth.go`, add alongside the existing `CanViewDocumentsFrom` helper:

```go
const ctxCanExport = "can_export"

// RequireExport aborts with 403 if the request context does not carry the can_export flag.
func RequireExport() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !CanExportFrom(c) {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": "export permission required",
                "code":  "errors.export.insufficientPermissions",
            })
            return
        }
        c.Next()
    }
}

// CanExportFrom reads the can_export flag set by RequireProjectRole.
func CanExportFrom(c *gin.Context) bool {
    v, _ := c.Get(ctxCanExport)
    b, _ := v.(bool)
    return b
}
```

Make sure `RequireProjectRole` stores the flag under `ctxCanExport` — check how `can_view_contact` etc. are stored and follow the same pattern.

**Step 2: Apply middleware to the export route group in server.go**

```go
// Export-level access (consultant+)
export := proj.Group("", projectAuthMW.RequireProjectRole(project.ActionExport), middleware.RequireExport())
{
    export.GET("/export/people", exportHandler.ExportPeople)
    export.GET("/export/support-records", exportHandler.ExportSupportRecords)
    export.GET("/export/pets", exportHandler.ExportPets)
    export.GET("/export/households", exportHandler.ExportHouseholds)
}
```

**Step 3: Remove the four manual checks in export_handler.go**

Delete the `if !middleware.CanExportFrom(c) { ... return }` blocks from `ExportPeople`, `ExportSupportRecords`, `ExportPets`, `ExportHouseholds`. The middleware enforces it before they are reached.

**Step 4: Write test**

```go
func TestExport_WithoutCanExport_Returns403(t *testing.T) {
    // Set up context with can_export = false
    // Call RequireExport() middleware
    // Assert 403
}
```

**Step 5: Run tests and commit**

```bash
just test
git add internal/middleware/project_auth.go internal/server/server.go internal/handler/export_handler.go
git commit -m "Move can_export check to route middleware; remove manual checks from export handlers"
```

---

## Task 24: MFA recovery codes

**Findings fixed:** AC-4 — device loss causes permanent lockout; no self-service recovery

**Files:**

- Create: `migrations/000035_add_mfa_recovery_codes.up.sql`
- Modify: `internal/domain/user/entity.go` (add `MFARecoveryCode` type)
- Modify: `internal/repository/interfaces.go` (add `MFARecoveryCodeRepository`)
- Create: `internal/repository/mfa_recovery_code_repository.go`
- Modify: `internal/usecase/auth/auth_usecase.go` (generate on EnableMFA, consume on VerifyMFA)
- Modify: `internal/usecase/auth/types.go` (add `EnableMFAOutput` with codes)
- Run: `just generate-mocks`

**Step 1: Create migration**

```sql
-- migrations/000035_add_mfa_recovery_codes.up.sql
CREATE TABLE mfa_recovery_codes (
    id          TEXT        PRIMARY KEY,
    user_id     TEXT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash   TEXT        NOT NULL,
    used_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX ON mfa_recovery_codes (user_id);
```

```bash
just migrate
```

**Step 2: Add domain type**

In `internal/domain/user/entity.go`:

```go
// MFARecoveryCode is a single-use backup code stored hashed.
type MFARecoveryCode struct {
    ID        string
    UserID    ulid.ULID
    CodeHash  string
    UsedAt    *time.Time
    CreatedAt time.Time
}
```

**Step 3: Add repository interface**

In `internal/repository/interfaces.go`:

```go
// MFARecoveryCodeRepository persists and validates MFA recovery codes.
type MFARecoveryCodeRepository interface {
    // CreateBatch stores a set of hashed recovery codes for a user.
    CreateBatch(ctx context.Context, codes []*user.MFARecoveryCode) error
    // FindUnused returns the first unused code matching the given hash.
    FindUnused(ctx context.Context, userID ulid.ULID, codeHash string) (*user.MFARecoveryCode, error)
    // MarkUsed marks a code as used (sets used_at = now).
    MarkUsed(ctx context.Context, id string) error
    // DeleteByUserID removes all recovery codes for a user (called on MFA disable).
    DeleteByUserID(ctx context.Context, userID ulid.ULID) error
}
```

**Step 4: Generate codes in `EnableMFA`**

In `auth_usecase.go`, after TOTP is verified and MFA config is saved, generate 8 recovery codes:

```go
import (
    "crypto/rand"
    "encoding/hex"
    "fmt"
)

// generateRecoveryCodes returns N plain-text codes and their Argon2id hashes.
func generateRecoveryCodes(n int, hasher crypto.PasswordHasher) ([]string, []*user.MFARecoveryCode, error) {
    codes := make([]string, n)
    records := make([]*user.MFARecoveryCode, n)
    for i := range codes {
        b := make([]byte, 8) // 16 hex chars
        if _, err := rand.Read(b); err != nil {
            return nil, nil, err
        }
        plain := hex.EncodeToString(b)
        hash, _, err := hasher.Hash(plain)
        if err != nil {
            return nil, nil, err
        }
        codes[i] = plain
        records[i] = &user.MFARecoveryCode{
            ID:       iulid.NewString(),
            UserID:   userID,
            CodeHash: hash,
        }
    }
    return codes, records, nil
}
```

Return the plain-text codes in `EnableMFAOutput` (shown once, not stored):

```go
type EnableMFAOutput struct {
    RecoveryCodes []string `json:"recovery_codes"`
}
```

**Step 5: Consume a recovery code in `VerifyMFA`**

When `TOTPCode` fails TOTP validation, try it as a recovery code:

```go
// If TOTP validation fails, try as a recovery code
if !totp.Validate(input.TOTPCode, mfaCfg.Secret) {
    // Try recovery code path
    hash, _, _ := uc.hasher.Hash(input.TOTPCode)
    rc, err := uc.recoveryRepo.FindUnused(ctx, userID, hash)
    if err != nil {
        return nil, user.ErrInvalidMFACode
    }
    _ = uc.recoveryRepo.MarkUsed(ctx, rc.ID)
    // fall through to session creation
} else {
    // normal TOTP success path
}
```

**Step 6: Delete codes on MFA disable**

In `DisableMFA`, after deleting the MFA config:

```go
_ = uc.recoveryRepo.DeleteByUserID(ctx, userID)
```

**Step 7: Regenerate mocks**

```bash
just generate-mocks
```

**Step 8: Run tests and commit**

```bash
just test
git add migrations/000035_add_mfa_recovery_codes.up.sql \
        internal/domain/user/entity.go \
        internal/repository/interfaces.go \
        internal/repository/mfa_recovery_code_repository.go \
        internal/usecase/auth/auth_usecase.go \
        internal/usecase/auth/types.go \
        internal/repository/mock/repository.go
git commit -m "Add MFA recovery codes: generate on enable, consume on verify, delete on disable"
```

---

## Deferred to Phase 2 (no code change now)

**L-4 — RSA Key Rotation**

JWT signing keys are loaded once at startup from `keys/jwt_rsa`. Key rotation requires a restart. A zero-downtime rotation strategy (support `JWT_PUBLIC_KEY_PATH_OLD` for verification during the rollover window) is a Phase 2 item. Access tokens expire in 15 minutes, so the blast radius of a key compromise is bounded.

---

## Final verification

After all tasks:

```bash
just test-all   # includes integration tests (Docker required)
```

---

## Execution order summary

| Task  | Finding(s)                                      | Severity    | Depends on |
| ----- | ----------------------------------------------- | ----------- | ---------- |
| 1     | Role self-assignment + email enumeration        | HIGH        | —          |
| 2     | CanLogin in VerifyMFA + RefreshToken            | HIGH        | —          |
| 3     | SessionRepository extensions + vacuum goroutine | HIGH/LOW    | —          |
| 4     | Session invalidation on password change/reset   | HIGH        | Task 3     |
| 5a–5h | Cross-project IDOR (8 entities)                 | HIGH        | —          |
| 6     | Document upload/serving security                | HIGH/MEDIUM | —          |
| 7     | Global body size limit (exempt multipart)       | MEDIUM      | —          |
| 8     | CSRF double-submit cookie                       | MEDIUM      | —          |
| 9     | Clamp pagination before repo call               | MEDIUM      | —          |
| 10    | Suppress 500 error details                      | MEDIUM      | —          |
| 11    | Staff can_export implicit override              | MEDIUM      | —          |
| 12    | Persist permanent lockouts to PostgreSQL        | MEDIUM      | —          |
| 13    | Remove lockout duration from response           | LOW         | —          |
| 14    | Argon2id time cost (see compat note)            | MEDIUM      | —          |
| 15    | Content-Security-Policy header                  | LOW         | —          |
| 16    | Wire session vacuum to startup                  | LOW         | Task 3     |
| 17    | Graceful shutdown                               | LOW         | —          |
| 18    | HTTPS API URL enforcement (frontend)            | MEDIUM      | —          |
| 19    | Audit logging for admin mutations               | ARCH        | —          |
| 20    | Stronger password policy                        | ARCH        | —          |
| 21    | Log IP/UA mismatch on token refresh             | ARCH        | —          |
| 22    | CORS localhost startup warning                  | LOW         | —          |
| 23    | can_export consolidated into route middleware   | ARCH        | —          |
| 24    | MFA recovery codes                              | ARCH        | —          |
| —     | RSA key rotation                                | LOW         | Phase 2    |
