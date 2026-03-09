---
title: Backend
weight: 3
---

Bul bet Go backend'ine salym qoşuu bojunça konvensijalardy qamtyjt. Eger azyrynça oquğan coq bolsoñuz, adeginde [Arhitektura](/docs/guide/architecture/) betin oquñuz.

## Cañy domen objektin qoşuu

Cañy objekt qoşuuda (mis. `document`) bul yraattuuluktu atqaryñyz:

1. **Domen** — `internal/domain/<name>/entity.go`: objekt strukturasy, enumdar, katalar
2. **Repozitorij interfejsi** — `internal/repository/interfaces.go` fajlyna qoşuu
3. **Repozitorij işke aşyruu** — `internal/repository/<name>_repository.go`
4. **Qoldonuu uçuru** — `internal/usecase/<group>/<name>_usecase.go` + tipter fajly
5. **Iştetkiç** — `internal/handler/<name>_handler.go`
6. **Marşruttar** — `internal/server/server.go` fajlyna tutaştyruu
7. **DI** — `internal/app/container.go` fajlynda repozitorij + qoldonuu uçurun tutaştyruu
8. **Migratsija** — `migrations/<seq>_create_<name>s_table.up.sql`
9. **Moqtor** — `internal/repository/interfaces.go` fajlyndağy `go:generate` direktivasyna interfejs qoşuu, `just generate-mocks` iştetüü
10. **Testter** — qoldonuu uçurun moqtor menen unit-testtöö; repozitorijdi testcontainers menen integratsiyalyq testtöö

## Atalyş konvensijalary

- Paket attary: qysqa, kiçine, cekelik — `user`, `project`, `support`
- Fajl attary: `<entity>_entity.go`, `<entity>_repository.go`, `<entity>_usecase.go`, `<entity>_handler.go`
- Objekt IDleri: strukturalarda `ulid.ULID`, DTOlarda `string` (`.String()` arqyluu)
- Caryjalanbaĝan repozitorij strukturalary: `type userRepository struct { db *sqlx.DB }`
- Konstruktor: `func NewUserRepository(db *sqlx.DB) repository.UserRepository`

## Iştetkiç ülgüsü

Iştetkiçter cuqa. Alar bajloo, çaqyruu, coop berüü ğana — başqa eç nerse emes.

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
        // ... req dan talalar
    })
    if err != nil {
        handleError(c, err)
        return
    }

    c.JSON(http.StatusCreated, out)
}
```

## Qoldonuu uçuru ülgüsü

Qoldonuu uçurlary repozitorijlerdi koordinasiyalayt. Alar HTTP ce SQL kodun qamtybayt.

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

## Qatany iştetüü

Domen qatalary `internal/domain/<name>/errors.go` fajlynda anyqtalat:

```go
var (
    ErrPersonNotFound = errors.New("person not found")
    ErrPersonExists   = errors.New("person already exists")
)
```

Iştetkiç domen qatalaryn `internal/handler/errors.go` fajlynda HTTP status koddoruna çağyldyrat. Kerek bolğondo cañy qatalardy oşol cerge qoşuñuz. Qoldonuu uçurlarynan özül maalymat bazasy qatalaryn qajtarbañyz.

## Migratsijalar

Calğyz aldyğa ğana — `.down.sql` fajldary coq. Fajl aty ülgüsü:

```
<seq>_<description>.up.sql
```

`<seq>` — nöl menen tolturulğan 6 orunduu san. Tüzüü:

```bash
observer migrate create <description>
# ce
just migrate-create <description>
```

Qoldonulğan migratsijany eç qaçan özgörtpöñüz. Anyn orduna cañysyn tüzüñüz.

## Köz qarandylykty inyeksijaloo

Bardyq tutaştyruu `internal/app/container.go` fajlynda bolot. Ülgü:

```go
// 1. Repozitorij tüzüü
personRepo := repository.NewPersonRepository(c.db.GetDB())

// 2. Qoldonuu uçurun tüzüü
personUC := projectUC.NewPersonUseCase(personRepo, ...)

// 3. Kontejnerde saqtoo
c.PersonUC = personUC
```

Andan kijin `internal/server/server.go` fajlynda iştetkiçke inyeksijaloo:

```go
personHandler := handler.NewPersonHandler(container.PersonUC)
```

## Kod stili

- Dekorativdüü kommentarij bölüüçülör coq (`//-----`, `//=====`)
- Docstring'der calğyz caryjalanğan simboldordo
- Tataal logika: syzyqtuu kommentarijlerdin orduna modul README fajlyndağy Mermaid diagrammasyn artyq körüñüz
- `gofmt` formattoo talap qylynat — commit casoodon murun `just fmt` iştetiñiz
- Linting: `just lint` (golangci-lint)
