---
title: Backend
weight: 3
---

Diese Seite beschreibt die Konventionen für Beiträge zum Go-Backend. Lesen Sie zuerst [Architektur](/docs/guide/architecture/), falls noch nicht geschehen.

## Hinzufügen einer neuen Domain-Entität

Folgen Sie dieser Reihenfolge beim Hinzufügen einer neuen Entität (z.B. `document`):

1. **Domain** — `internal/domain/<name>/entity.go`: Entitätsstruktur, Enumerationen, Fehler
2. **Repository-Schnittstelle** — zu `internal/repository/interfaces.go` hinzufügen
3. **Repository-Implementierung** — `internal/repository/<name>_repository.go`
4. **Use Case** — `internal/usecase/<group>/<name>_usecase.go` + Typen-Datei
5. **Handler** — `internal/handler/<name>_handler.go`
6. **Routen** — in `internal/server/server.go` verdrahten
7. **DI** — Repository + Use Case in `internal/app/container.go` verdrahten
8. **Migration** — `migrations/<seq>_create_<name>s_table.up.sql`
9. **Mocks** — Schnittstelle zur `go:generate`-Direktive in `internal/repository/interfaces.go` hinzufügen, `just generate-mocks` ausführen
10. **Tests** — Use Case mit Mocks unit-testen; Repository mit testcontainers integrationstest durchführen

## Namenskonventionen

- Paketnamen: kurz, kleingeschrieben, Singular — `user`, `project`, `support`
- Dateinamen: `<entity>_entity.go`, `<entity>_repository.go`, `<entity>_usecase.go`, `<entity>_handler.go`
- Entitäts-IDs: `ulid.ULID` in Strukturen, `string` in DTOs (via `.String()`)
- Nicht-exportierte Repo-Strukturen: `type userRepository struct { db *sqlx.DB }`
- Konstruktor: `func NewUserRepository(db *sqlx.DB) repository.UserRepository`

## Handler-Muster

Handler sind schlank. Sie binden, rufen auf, antworten — nichts weiter.

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
        // ... Felder aus req
    })
    if err != nil {
        handleError(c, err)
        return
    }

    c.JSON(http.StatusCreated, out)
}
```

## Use-Case-Muster

Use Cases koordinieren Repositories. Sie enthalten keinen HTTP- oder SQL-Code.

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

## Fehlerbehandlung

Domänenfehler werden in `internal/domain/<name>/errors.go` definiert:

```go
var (
    ErrPersonNotFound = errors.New("person not found")
    ErrPersonExists   = errors.New("person already exists")
)
```

Der Handler ordnet Domänenfehler HTTP-Statuscodes in `internal/handler/errors.go` zu. Fügen Sie dort bei Bedarf neue Fehler hinzu. Geben Sie keine rohen Datenbankfehler aus Use Cases zurück.

## Migrationen

Nur vorwärts — keine `.down.sql`-Dateien. Dateinamenmuster:

```
<seq>_<description>.up.sql
```

Wobei `<seq>` eine sechsstellige Zahl mit führenden Nullen ist. Erstellen mit:

```bash
observer migrate create <description>
# oder
just migrate-create <description>
```

Ändern Sie niemals eine bereits angewendete Migration. Erstellen Sie stattdessen eine neue.

## Dependency Injection

Alle Verdrahtung erfolgt in `internal/app/container.go`. Das Muster:

```go
// 1. Repository erstellen
personRepo := repository.NewPersonRepository(c.db.GetDB())

// 2. Use Case erstellen
personUC := projectUC.NewPersonUseCase(personRepo, ...)

// 3. Im Container speichern
c.PersonUC = personUC
```

Dann in `internal/server/server.go` in den Handler injizieren:

```go
personHandler := handler.NewPersonHandler(container.PersonUC)
```

## Code-Stil

- Keine dekorativen Kommentartrennzeichen (`//-----`, `//=====`)
- Docstrings nur für exportierte Symbole
- Komplexe Logik: bevorzugen Sie ein Mermaid-Diagramm in einem Modul-README gegenüber Inline-Kommentaren
- `gofmt`-Formatierung erzwungen — `just fmt` vor dem Commit ausführen
- Linting: `just lint` (golangci-lint)
