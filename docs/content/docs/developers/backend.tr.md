---
title: Backend
weight: 3
---

Bu sayfa, Go backend'ine katkıda bulunmaya yönelik kuralları kapsar. Henüz okumadıysanız önce [Mimari](/docs/guide/architecture/) sayfasını okuyun.

## Yeni Bir Domain Varlığı Ekleme

Yeni bir varlık (örn. `document`) eklerken şu sırayı izleyin:

1. **Domain** — `internal/domain/<name>/entity.go`: varlık yapısı, enum'lar, hatalar
2. **Repository arayüzü** — `internal/repository/interfaces.go` dosyasına ekleyin
3. **Repository implementasyonu** — `internal/repository/<name>_repository.go`
4. **Use case** — `internal/usecase/<group>/<name>_usecase.go` + tipler dosyası
5. **Handler** — `internal/handler/<name>_handler.go`
6. **Route'lar** — `internal/server/server.go` dosyasına bağlayın
7. **DI** — `internal/app/container.go` dosyasında repository + use case bağlantısını yapın
8. **Migration** — `migrations/<seq>_create_<name>s_table.up.sql`
9. **Mock'lar** — `internal/repository/interfaces.go` dosyasındaki `go:generate` direktifine arayüz ekleyin, `just generate-mocks` çalıştırın
10. **Testler** — use case'i mock'larla unit test edin; repository'yi testcontainers ile entegrasyon testi yapın

## İsimlendirme Kuralları

- Paket adları: kısa, küçük harf, tekil — `user`, `project`, `support`
- Dosya adları: `<entity>_entity.go`, `<entity>_repository.go`, `<entity>_usecase.go`, `<entity>_handler.go`
- Varlık ID'leri: struct'larda `ulid.ULID`, DTO'larda `string` (`.String()` üzerinden)
- Dışa aktarılmamış repo struct'ları: `type userRepository struct { db *sqlx.DB }`
- Constructor: `func NewUserRepository(db *sqlx.DB) repository.UserRepository`

## Handler Deseni

Handler'lar incedir. Bağlar, çağırır, yanıt verir — başka hiçbir şey yapmaz.

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
        // ... req'den alanlar
    })
    if err != nil {
        handleError(c, err)
        return
    }

    c.JSON(http.StatusCreated, out)
}
```

## Use Case Deseni

Use case'ler repository'leri koordine eder. HTTP veya SQL kodu içermez.

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

## Hata Yönetimi

Domain hataları `internal/domain/<name>/errors.go` dosyasında tanımlanır:

```go
var (
    ErrPersonNotFound = errors.New("person not found")
    ErrPersonExists   = errors.New("person already exists")
)
```

Handler, `internal/handler/errors.go` dosyasında domain hatalarını HTTP durum kodlarına eşler. Gerektiğinde yeni hataları oraya ekleyin. Use case'lerden ham veritabanı hataları döndürmeyin.

## Migration'lar

Yalnızca ileri yönlü — `.down.sql` dosyası yok. Dosya adı deseni:

```
<seq>_<description>.up.sql
```

`<seq>`, sıfırla doldurulmuş 6 basamaklı bir sayıdır. Şu komutla oluşturun:

```bash
observer migrate create <description>
# veya
just migrate-create <description>
```

Uygulanmış bir migration'ı asla değiştirmeyin. Bunun yerine yeni bir tane oluşturun.

## Bağımlılık Enjeksiyonu

Tüm bağlantı `internal/app/container.go` dosyasında gerçekleşir. Desen:

```go
// 1. Repository oluştur
personRepo := repository.NewPersonRepository(c.db.GetDB())

// 2. Use case oluştur
personUC := projectUC.NewPersonUseCase(personRepo, ...)

// 3. Container'da sakla
c.PersonUC = personUC
```

Ardından `internal/server/server.go` dosyasında handler'a enjekte et:

```go
personHandler := handler.NewPersonHandler(container.PersonUC)
```

## Kod Stili

- Dekoratif yorum ayırıcıları yok (`//-----`, `//=====`)
- Docstring'ler yalnızca dışa aktarılan semboller için
- Karmaşık mantık: satır içi yorumlar yerine modül README'sindeki Mermaid diyagramını tercih edin
- `gofmt` formatlaması zorunlu — commit etmeden önce `just fmt` çalıştırın
- Linting: `just lint` (golangci-lint)
