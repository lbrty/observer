---
title: Backend
weight: 3
---

Эта страница описывает соглашения для внесения вклада в Go-бэкенд. Если вы ещё не читали, сначала прочитайте страницу [Архитектура](/docs/guide/architecture/).

## Добавление нового доменного объекта

Следуйте этой последовательности при добавлении нового объекта (например, `document`):

1. **Домен** — `internal/domain/<name>/entity.go`: структура объекта, перечисления, ошибки
2. **Интерфейс репозитория** — добавить в `internal/repository/interfaces.go`
3. **Реализация репозитория** — `internal/repository/<name>_repository.go`
4. **Сценарий использования** — `internal/usecase/<group>/<name>_usecase.go` + файл типов
5. **Обработчик** — `internal/handler/<name>_handler.go`
6. **Маршруты** — подключить в `internal/server/server.go`
7. **DI** — подключить репозиторий + сценарий использования в `internal/app/container.go`
8. **Миграция** — `migrations/<seq>_create_<name>s_table.up.sql`
9. **Моки** — добавить интерфейс в директиву `go:generate` в `internal/repository/interfaces.go`, запустить `just generate-mocks`
10. **Тесты** — юнит-тест сценария использования с моками; интеграционный тест репозитория с testcontainers

## Соглашения по именованию

- Имена пакетов: короткие, строчные, в единственном числе — `user`, `project`, `support`
- Имена файлов: `<entity>_entity.go`, `<entity>_repository.go`, `<entity>_usecase.go`, `<entity>_handler.go`
- ID объектов: `ulid.ULID` в структурах, `string` в DTO (через `.String()`)
- Неэкспортируемые структуры репозитория: `type userRepository struct { db *sqlx.DB }`
- Конструктор: `func NewUserRepository(db *sqlx.DB) repository.UserRepository`

## Шаблон обработчика

Обработчики тонкие. Они привязывают, вызывают, отвечают — и больше ничего.

```go
func (h *PersonHandler) Create(c *gin.Context) {
    req, ok := bindJSON[CreatePersonRequest](c)
    if !ok {
        return
    }

    projectID := c.Param("project_id")
    userID, _ := middleware.UserIDFrom(c)

    out, err := h.personUC.Create(c.Request.Context(), usecase.CreatePersonInput{
        ProjectID: projectID,
        CreatedBy: userID.String(),
        // ... поля из req
    })
    if err != nil {
        HandleError(c, err)
        return
    }

    c.JSON(http.StatusCreated, out)
}
```

`bindJSON[T]` определена в `internal/handler/errors.go`. Она декодирует тело запроса и записывает ответ `400` при ошибке, возвращая `false`, чтобы обработчик мог немедленно завершиться.

## Шаблон сценария использования

Сценарии использования координируют репозитории. Они не содержат HTTP или SQL кода.

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

## Обработка ошибок

Доменные ошибки определяются в `internal/domain/<name>/errors.go`:

```go
var (
    ErrPersonNotFound = errors.New("person not found")
    ErrPersonExists   = errors.New("person already exists")
)
```

Обработчик сопоставляет доменные ошибки с HTTP-кодами статуса в `internal/handler/errors.go`. Добавляйте туда новые ошибки по мере необходимости. Не возвращайте сырые ошибки базы данных из сценариев использования.

## Миграции

Только вперёд — без файлов `.down.sql`. Шаблон имени файла:

```
<seq>_<description>.up.sql
```

Где `<seq>` — шестизначное число с ведущими нулями. Создать с помощью:

```bash
observer migrate create <description>
# или
just migrate-create <description>
```

Никогда не изменяйте применённую миграцию. Вместо этого создайте новую.

## Внедрение зависимостей

Всё подключение происходит в `internal/app/container.go`. Шаблон:

```go
// 1. Создать репозиторий
personRepo := repository.NewPersonRepository(c.db.GetDB())

// 2. Создать сценарий использования
personUC := projectUC.NewPersonUseCase(personRepo, ...)

// 3. Сохранить в контейнере
c.PersonUC = personUC
```

Затем в `internal/server/server.go` внедрить в обработчик:

```go
personHandler := handler.NewPersonHandler(container.PersonUC)
```

## Стиль кода

- Никаких декоративных разделителей комментариев (`//-----`, `//=====`)
- Docstring только для экспортируемых символов
- Сложная логика: предпочитайте диаграмму Mermaid в README модуля вместо встроенных комментариев
- Форматирование `gofmt` обязательно — запускайте `just fmt` перед коммитом
- Линтинг: `just lint` (golangci-lint)
