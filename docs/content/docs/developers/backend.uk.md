---
title: Backend
weight: 3
---

Ця сторінка описує угоди для внесення вкладу у Go-бекенд. Якщо ви ще не читали, спочатку прочитайте сторінку [Архітектура](/docs/guide/architecture/).

## Додавання нової доменної сутності

Дотримуйтесь цієї послідовності при додаванні нової сутності (наприклад, `document`):

1. **Домен** — `internal/domain/<name>/entity.go`: структура сутності, перерахування, помилки
2. **Інтерфейс репозиторію** — додати до `internal/repository/interfaces.go`
3. **Реалізація репозиторію** — `internal/repository/<name>_repository.go`
4. **Варіант використання** — `internal/usecase/<group>/<name>_usecase.go` + файл типів
5. **Обробник** — `internal/handler/<name>_handler.go`
6. **Маршрути** — підключити до `internal/server/server.go`
7. **DI** — підключити репозиторій + варіант використання в `internal/app/container.go`
8. **Міграція** — `migrations/<seq>_create_<name>s_table.up.sql`
9. **Моки** — додати інтерфейс до директиви `go:generate` у `internal/repository/interfaces.go`, запустити `just generate-mocks`
10. **Тести** — unit-тест варіанту використання з моками; інтеграційний тест репозиторію з testcontainers

## Угоди щодо іменування

- Назви пакетів: короткі, малі літери, однина — `user`, `project`, `support`
- Назви файлів: `<entity>_entity.go`, `<entity>_repository.go`, `<entity>_usecase.go`, `<entity>_handler.go`
- ID сутностей: `ulid.ULID` у структурах, `string` у DTO (через `.String()`)
- Не-експортовані структури репозиторію: `type userRepository struct { db *sqlx.DB }`
- Конструктор: `func NewUserRepository(db *sqlx.DB) repository.UserRepository`

## Шаблон обробника

Обробники тонкі. Вони прив'язують, викликають, відповідають — і більше нічого.

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
        // ... поля з req
    })
    if err != nil {
        handleError(c, err)
        return
    }

    c.JSON(http.StatusCreated, out)
}
```

## Шаблон варіанту використання

Варіанти використання координують репозиторії. Вони не містять HTTP або SQL коду.

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

## Обробка помилок

Доменні помилки визначаються в `internal/domain/<name>/errors.go`:

```go
var (
    ErrPersonNotFound = errors.New("person not found")
    ErrPersonExists   = errors.New("person already exists")
)
```

Обробник зіставляє доменні помилки з HTTP-кодами статусу в `internal/handler/errors.go`. Додавайте туди нові помилки за потреби. Не повертайте сирі помилки бази даних з варіантів використання.

## Міграції

Лише вперед — без файлів `.down.sql`. Шаблон імені файлу:

```
<seq>_<description>.up.sql
```

Де `<seq>` — шестизначне число з провідними нулями. Створити за допомогою:

```bash
observer migrate create <description>
# або
just migrate-create <description>
```

Ніколи не змінюйте застосовану міграцію. Замість цього створіть нову.

## Впровадження залежностей

Все підключення відбувається в `internal/app/container.go`. Шаблон:

```go
// 1. Створити репозиторій
personRepo := repository.NewPersonRepository(c.db.GetDB())

// 2. Створити варіант використання
personUC := projectUC.NewPersonUseCase(personRepo, ...)

// 3. Зберегти в контейнері
c.PersonUC = personUC
```

Потім у `internal/server/server.go` впровадити в обробник:

```go
personHandler := handler.NewPersonHandler(container.PersonUC)
```

## Стиль коду

- Ніяких декоративних роздільників коментарів (`//-----`, `//=====`)
- Docstring лише для експортованих символів
- Складна логіка: надавайте перевагу діаграмі Mermaid у README модуля, а не рядковим коментарям
- Форматування `gofmt` обов'язкове — запускайте `just fmt` перед комітом
- Лінтинг: `just lint` (golangci-lint)
