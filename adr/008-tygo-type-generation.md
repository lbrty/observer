# ADR-008: Tygo Type Generation


| Field      | Value                                          |
| ---------- | ---------------------------------------------- |
| Status     | Proposed                                       |
| Date       | 2026-02-24                                     |
| Components | observer (backend DTOs), observer-web (types/) |


## Decision

Use [tygo](https://github.com/gzuidhof/tygo) to generate TypeScript interfaces from Go structs.

### Why tygo

- Reads Go source directly (no runtime reflection, no OpenAPI intermediate step).
- Respects `json:` tags for field names and `omitempty` for optionality.
- Handles pointer types (`*string` → `string | undefined`).
- Supports custom type mappings (`time.Time` → `string`).
- Preserves Go doc comments as TSDoc.
- Understands const blocks → TypeScript string literal unions.
- Single binary, no runtime dependency.

### Configuration

**`tygo.yaml`** at repository root:

```yaml
packages:
  - path: github.com/lbrty/observer/internal/usecase/auth
    output_path: packages/observer-web/src/types/generated/auth.ts
    type_mappings:
      time.Time: string

  - path: github.com/lbrty/observer/internal/usecase/admin
    output_path: packages/observer-web/src/types/generated/admin.ts
    type_mappings:
      time.Time: string

  - path: github.com/lbrty/observer/internal/usecase/project
    output_path: packages/observer-web/src/types/generated/project.ts
    type_mappings:
      time.Time: string
```

### Generated output

Given this Go source:

```go
// UserDTO is the admin-facing user representation.
type UserDTO struct {
    ID         string    `json:"id"`
    FirstName  string    `json:"first_name"`
    LastName   string    `json:"last_name"`
    Email      string    `json:"email"`
    Phone      string    `json:"phone"`
    OfficeID   *string   `json:"office_id,omitempty"`
    Role       string    `json:"role"`
    IsVerified bool      `json:"is_verified"`
    IsActive   bool      `json:"is_active"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

tygo produces:

```ts
/** UserDTO is the admin-facing user representation. */
export interface UserDTO {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  office_id?: string;
  role: string;
  is_verified: boolean;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
```

### Build integration

**Justfile target:**

```just
generate-types:
    tygo generate
```

**Workflow:**

1. Developer modifies a Go DTO struct.
2. Runs `just generate-types`.
3. Generated `.ts` files update in `packages/observer-web/src/types/generated/`.
4. Frontend code imports from `@/types/generated/admin`.
5. Generated files are committed to the repository (no CI generation step needed).

### Frontend re-export pattern

`packages/observer-web/src/types/admin.ts` re-exports with any frontend-only additions:

```ts
export type { UserDTO, ListUsersOutput, UpdateUserInput } from "./generated/admin";

// Frontend-only types (query params, form state, etc.)
export interface ListUsersParams {
  page?: number;
  per_page?: number;
  search?: string;
  role?: string;
  is_active?: boolean;
}
```

This keeps generated files untouched while allowing frontend-specific types alongside them.

### What stays hand-written

- Query parameter types (`ListUsersParams`) — Go uses `form:` tags, not `json:`, and tygo only reads `json:` tags.
- Wrapper response types that the frontend destructures differently.
- React-specific types (form state, component props).


## Alternatives Considered

### A. OpenAPI spec generation + openapi-typescript

Generate an OpenAPI spec from Go (via swag or similar), then generate TS types from the spec.

**Rejected because**: two-step pipeline, heavier tooling, swag annotations add noise to handler code, and we already have clean DTO structs that tygo reads directly.

### B. Hand-written TypeScript types

Continue writing TS interfaces manually to match Go DTOs.

**Rejected because**: drift is inevitable as the API surface grows. Already experienced with `auth.ts` types diverging from `usecase/auth/types.go` (e.g. `User` vs `UserDTO` field differences).

### C. Shared JSON Schema

Define types in JSON Schema, generate both Go and TS from it.

**Rejected because**: JSON Schema as source of truth is awkward for Go development. Developers think in structs, not schemas. Adds a third artifact to maintain.
