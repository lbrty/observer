---
title: Backend
weight: 3
---

This page covers conventions for contributing to the Go backend. Read [Architecture](/docs/guide/architecture/) first if you haven't already.

## Adding a New Domain Entity

Follow this sequence when adding a new entity (e.g. `document`):

1. **Domain** — `internal/domain/<name>/entity.go`: entity struct, enums, errors
2. **Repository interface** — add to `internal/repository/interfaces.go`
3. **Repository implementation** — `internal/repository/<name>_repository.go`
4. **Use case** — `internal/usecase/<group>/<name>_usecase.go` + types file
5. **Handler** — `internal/handler/<name>_handler.go`
6. **Routes** — wire into `internal/server/server.go`
7. **DI** — wire repository + use case in `internal/app/container.go`
8. **Migration** — `migrations/<seq>_create_<name>s_table.up.sql`
9. **Mocks** — add interface to `go:generate` in `internal/repository/interfaces.go`, run `just generate-mocks`
10. **Tests** — unit test the use case with mocks; integration test the repository with testcontainers

## Naming Conventions

- Package names: short, lowercase, singular — `user`, `project`, `support`
- File names: `<entity>_entity.go`, `<entity>_repository.go`, `<entity>_usecase.go`, `<entity>_handler.go`
- Entity IDs: `ulid.ULID` in structs, `string` in DTOs (via `.String()`)
- Unexported repo structs: `type userRepository struct { db *sqlx.DB }`
- Constructor: `func NewUserRepository(db *sqlx.DB) repository.UserRepository`

## Handler Pattern

Handlers are thin. They bind, call, respond — nothing else.

```go
func (h *PersonHandler) Create(c *gin.Context) {
    var req CreatePersonRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    projectID := c.Param("project_id")
    userID, _ := middleware.UserIDFrom(c)

    out, err := h.personUC.Create(c.Request.Context(), usecase.CreatePersonInput{
        ProjectID: projectID,
        CreatedBy: userID.String(),
        // ... fields from req
    })
    if err != nil {
        handleError(c, err)
        return
    }

    c.JSON(http.StatusCreated, out)
}
```

## Use Case Pattern

Use cases coordinate repositories. They contain no HTTP or SQL code.

```go
type CreatePersonInput struct {
    ProjectID string
    FirstName string
    LastName  string
    // ...
}

type CreatePersonOutput struct {
    PersonID string `json:"person_id"`
}

func (uc *PersonUseCase) Create(ctx context.Context, in CreatePersonInput) (*CreatePersonOutput, error) {
    person := &domain.Person{
        ID:        ulid.New(),
        ProjectID: in.ProjectID,
        FirstName: in.FirstName,
        // ...
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
    }

    if err := uc.repo.Create(ctx, person); err != nil {
        return nil, err
    }

    return &CreatePersonOutput{PersonID: person.ID.String()}, nil
}
```

## Error Handling

Domain errors are defined in `internal/domain/<name>/errors.go`:

```go
var (
    ErrPersonNotFound = errors.New("person not found")
    ErrPersonExists   = errors.New("person already exists")
)
```

The handler maps domain errors to HTTP status codes in `internal/handler/errors.go`. Add new errors there when needed. Do not return raw database errors from use cases.

## Migrations

Forward-only only — no `.down.sql` files. Filename pattern:

```
<seq>_<description>.up.sql
```

Where `<seq>` is a zero-padded 6-digit number. Create with:

```bash
observer migrate create <description>
# or
just migrate-create <description>
```

Never modify an applied migration. Create a new one instead.

## Dependency Injection

All wiring happens in `internal/app/container.go`. The pattern:

```go
// 1. Create repo
personRepo := repository.NewPersonRepository(c.db.GetDB())

// 2. Create use case
personUC := projectUC.NewPersonUseCase(personRepo, ...)

// 3. Store in container
c.PersonUC = personUC
```

Then in `internal/server/server.go`, inject into the handler:

```go
personHandler := handler.NewPersonHandler(container.PersonUC)
```

## Code Style

- No decorative comment separators (`//-----`, `//=====`)
- Docstrings only on exported symbols
- Complex logic: prefer a Mermaid diagram in a module README over inline comments
- `gofmt` formatting enforced — run `just fmt` before committing
- Linting: `just lint` (golangci-lint)
