---
title: Тестування
weight: 4
---

[package_name] використовує три рівні тестування: модульні тести, інтеграційні тести з testcontainers та E2E-тести через HTTP.

```mermaid
graph TD
    A[just test] -->|"-short flag"| B[Unit Tests]
    C[just test-all] --> B
    C --> D[Integration Tests]
    C --> E[E2E Tests]

    B --> F[Mocks + Testify]
    D --> G[Testcontainers: Postgres/Redis]
    E --> H[httptest + Real Router]
```

## Запуск тестів

```bash
just test              # лише модульні тести (швидко, без Docker)
just test-all          # всі тести, включно з інтеграційними
just test-coverage     # генерація HTML-звіту покриття
just test-race         # виявлення станів гонитви
just bench             # бенчмарки
just generate-mocks    # перегенерація файлів моків
```

## Конвенції тестових файлів

| Шаблон            | Призначення                                   |
| ----------------- | --------------------------------------------- |
| `*_test.go`       | Тестовий файл у тому ж пакеті                 |
| `testing.Short()` | Пропуск інтеграційних тестів у `just test`     |
| `t.Helper()`      | Позначення функцій як допоміжних для тестів   |
| `testutil.Setup*` | Налаштування контейнерів для інтеграційних тестів |

## 1. Модульні тести з Testify

Використовуйте `assert` для нефатальних перевірок (тест продовжується) та `require` для фатальних перевірок (тест зупиняється).

```go
package auth

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestValidateRole(t *testing.T) {
    role, err := user.ValidateRole("user")
    require.NoError(t, err)        // fatal: stop if error
    assert.Equal(t, user.RoleUser, role)  // non-fatal: continue if wrong
}

func TestValidateRole_Invalid(t *testing.T) {
    _, err := user.ValidateRole("superadmin")
    assert.ErrorIs(t, err, user.ErrInvalidRole)
}
```

### Поширені перевірки Testify

```go
require.NoError(t, err)                    // fatal if err != nil
require.NotNil(t, obj)                     // fatal if nil
assert.Equal(t, expected, actual)          // compare values
assert.NotEqual(t, a, b)                  // values differ
assert.Contains(t, str, "substring")      // substring check
assert.Len(t, slice, 3)                   // length check
assert.True(t, condition)                 // boolean check
assert.Error(t, err)                      // expect error
assert.ErrorIs(t, err, ErrSpecific)       // error type check
assert.NotEmpty(t, val)                   // non-empty check
```

### Табличні тести

```go
func TestUser_CanLogin(t *testing.T) {
    tests := []struct {
        name    string
        user    user.User
        wantErr error
    }{
        {
            name:    "active user can login",
            user:    user.User{IsActive: true},
            wantErr: nil,
        },
        {
            name:    "inactive user cannot login",
            user:    user.User{IsActive: false},
            wantErr: user.ErrUserNotActive,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.CanLogin()
            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## 2. Мокування з gomock

Моки генеруються з інтерфейсів за допомогою директив `go:generate`.

### Генерація моків

Директива в `internal/user/repository.go`:

```go
//go:generate mockgen -destination=mock/repository.go -package=mock [package_name]/internal/user UserRepository,CredentialsRepository,SessionRepository,MFARepository,VerificationTokenRepository
```

Виконайте `just generate-mocks` для перегенерації всіх моків після зміни інтерфейсів.

### Використання моків у тестах

```go
package auth_test

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"

    "[package_name]/internal/auth"
    "[package_name]/internal/user"
    mock_user "[package_name]/internal/user/mock"
)

func TestRegisterUseCase_Execute(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUserRepo := mock_user.NewMockUserRepository(ctrl)
    mockCredRepo := mock_user.NewMockCredentialsRepository(ctrl)
    hasher := auth.NewArgonHasher()

    uc := auth.NewRegisterUseCase(mockUserRepo, mockCredRepo, hasher)

    ctx := context.Background()
    input := auth.RegisterInput{
        Email:    "test@example.com",
        Phone:    "+49555000111",
        Password: "securepassword",
        Role:     "user",
    }

    // Setup expectations: email and phone don't exist yet
    mockUserRepo.EXPECT().
        GetByEmail(ctx, input.Email).
        Return(nil, errors.New("not found"))

    mockUserRepo.EXPECT().
        GetByPhone(ctx, input.Phone).
        Return(nil, errors.New("not found"))

    mockUserRepo.EXPECT().
        Create(ctx, gomock.Any()).
        Return(nil)

    mockCredRepo.EXPECT().
        Create(ctx, gomock.Any()).
        Return(nil)

    out, err := uc.Execute(ctx, input)
    require.NoError(t, err)
    assert.NotEmpty(t, out.UserID)
    assert.Contains(t, out.Message, "Registration successful")
}

func TestRegisterUseCase_EmailExists(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUserRepo := mock_user.NewMockUserRepository(ctrl)
    mockCredRepo := mock_user.NewMockCredentialsRepository(ctrl)
    hasher := auth.NewArgonHasher()

    uc := auth.NewRegisterUseCase(mockUserRepo, mockCredRepo, hasher)

    // Email already exists - return a user (no error)
    mockUserRepo.EXPECT().
        GetByEmail(gomock.Any(), "taken@example.com").
        Return(&user.User{ID: "existing"}, nil)

    _, err := uc.Execute(context.Background(), auth.RegisterInput{
        Email:    "taken@example.com",
        Phone:    "+49555000222",
        Password: "securepassword",
        Role:     "user",
    })
    assert.ErrorIs(t, err, user.ErrEmailExists)
}
```

### Матчери gomock

```go
gomock.Any()                          // match any value
gomock.Eq("exact")                    // exact match
gomock.Not(gomock.Eq("excluded"))     // negation
gomock.Nil()                          // match nil
```

### Очікування порядку викликів

```go
first := mockRepo.EXPECT().GetByEmail(gomock.Any(), "a@b.com").Return(nil, errNotFound)
mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).After(first)
```

## 3. Інтеграційні тести з Testcontainers

Інтеграційні тести запускають реальні Docker-контейнери для Postgres та Redis.

### Структура тесту

```go
func TestUserRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Start real Postgres container
    dsn, cleanup := testutil.SetupPostgres(t)
    defer cleanup()

    // Connect to database
    db, err := database.New(dsn)
    require.NoError(t, err)
    defer db.Close()

    // Run migrations
    // ... apply schema ...

    // Test real queries
    repo := postgres.NewUserRepository(db.GetDB())
    ctx := context.Background()

    u := &user.User{
        ID:        ulid.New(),
        Email:     "test@example.com",
        Phone:     "+49555000111",
        Role:      user.RoleUser,
        IsActive:  true,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
    }

    err = repo.Create(ctx, u)
    require.NoError(t, err)

    got, err := repo.GetByEmail(ctx, "test@example.com")
    require.NoError(t, err)
    assert.Equal(t, u.ID, got.ID)
    assert.Equal(t, u.Email, got.Email)
}
```

### Допоміжна функція: testutil.SetupPostgres

Знаходиться у `internal/testutil/postgres.go`. Повертає DSN та функцію очищення:

```go
dsn, cleanup := testutil.SetupPostgres(t)
defer cleanup()
```

Автоматично пропускається, коли Docker недоступний.

### Допоміжна функція: testutil.SetupRedis

Знаходиться у `internal/testutil/redis.go`:

```go
addr, cleanup := testutil.SetupRedis(t)
defer cleanup()
```

## 4. Тести HTTP-обробників

Використовуйте `httptest` з маршрутизатором Gin для тестування HTTP-ендпоінтів без запущеного сервера.

```go
package server_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "go.uber.org/mock/gomock"

    mock_database "[package_name]/internal/database/mock"
)

func TestHealthEndpoint(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockDB := mock_database.NewMockDB(ctrl)
    s := server.New(cfg, mockDB, log, nil)

    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    s.Router().ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), `"status":"ok"`)
}
```

### Тестування POST-ендпоінтів

```go
func TestRegisterEndpoint(t *testing.T) {
    // ... setup server with mocks ...

    body := `{"email":"a@b.com","phone":"+49555111222","password":"securepass","role":"user"}`

    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    s.Router().ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)

    var resp auth.RegisterOutput
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    require.NoError(t, err)
    assert.NotEmpty(t, resp.UserID)
}
```

## 5. E2E: Реєстрація -> Вхід -> Оновлення токена -> Вихід

Повний потік автентифікації, протестований наскрізно через HTTP.

```go
func TestAuthFlow_E2E(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping e2e test")
    }

    dsn, cleanup := testutil.SetupPostgres(t)
    defer cleanup()

    // Setup real DB, run migrations, create server with real deps
    db, err := database.New(dsn)
    require.NoError(t, err)
    defer db.Close()

    // ... apply migrations, wire dependencies ...

    router := s.Router()

    // Step 1: Register
    regBody := `{"[package_name].kg","phone":"+49700111222","password":"MyStr0ngPass!","role":"user"}`
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(regBody))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusCreated, w.Code)

    // Step 2: Login
    loginBody := `{"[package_name].kg","password":"MyStr0ngPass!"}`
    w = httptest.NewRecorder()
    req = httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(loginBody))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)

    var loginResp auth.LoginOutput
    err = json.Unmarshal(w.Body.Bytes(), &loginResp)
    require.NoError(t, err)
    assert.NotEmpty(t, loginResp.Tokens.AccessToken)
    assert.NotEmpty(t, loginResp.Tokens.RefreshToken)

    // Step 3: Refresh token
    refreshBody := fmt.Sprintf(`{"refresh_token":"%s"}`, loginResp.Tokens.RefreshToken)
    w = httptest.NewRecorder()
    req = httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(refreshBody))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)

    var refreshResp auth.TokenPair
    err = json.Unmarshal(w.Body.Bytes(), &refreshResp)
    require.NoError(t, err)
    assert.NotEmpty(t, refreshResp.AccessToken)
    assert.NotEqual(t, loginResp.Tokens.RefreshToken, refreshResp.RefreshToken)

    // Step 4: Logout
    logoutBody := fmt.Sprintf(`{"refresh_token":"%s"}`, refreshResp.RefreshToken)
    w = httptest.NewRecorder()
    req = httptest.NewRequest(http.MethodPost, "/auth/logout", strings.NewReader(logoutBody))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)

    // Step 5: Refresh with old token should fail
    w = httptest.NewRecorder()
    req = httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(logoutBody))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}
```

## 6. Контрольний список критичних тестів

### Домен автентифікації

| Тест                            | Тип    | Що перевіряє                           |
| ------------------------------- | ------ | -------------------------------------- |
| Реєстрація з валідними даними   | Unit   | Користувач + облікові дані створені    |
| Реєстрація з дублікатом email   | Unit   | Повертає `ErrEmailExists`              |
| Реєстрація з дублікатом телефону| Unit   | Повертає `ErrPhoneExists`              |
| Реєстрація з невалідною роллю   | Unit   | Повертає `ErrInvalidRole`              |
| Вхід з валідними обліковими даними | Unit | Повертає пару токенів                  |
| Вхід з невірним паролем         | Unit   | Повертає `ErrInvalidCredentials`       |
| Вхід неактивного користувача    | Unit   | Повертає `ErrUserNotActive`            |
| Оновлення валідного токена      | Unit   | Стара сесія видалена, нова пара видана |
| Оновлення простроченої сесії    | Unit   | Повертає `ErrSessionExpired`           |
| Генерація + валідація токенів   | Unit   | Claims збігаються, термін дії працює   |
| Невідповідність типу токена     | Unit   | Access token відхилений як MFA        |
| Унікальність хешу пароля        | Unit   | Однаковий пароль -> різні хеші        |

### Інфраструктура

| Тест                    | Тип         | Що перевіряє                          |
| ----------------------- | ----------- | ------------------------------------- |
| Підключення до БД + ping| Integration | Testcontainer Postgres працює         |
| Health endpoint         | Unit        | Повертає 200 `{"status":"ok"}`        |
| Request ID middleware   | Unit        | Заголовок X-Request-ID — 26-символьний ULID |
| Graceful shutdown       | Unit        | Сервер завершується без помилки       |
| Значення конфігурації за замовчуванням | Unit | Адекватні значення завантажені |
| Перевизначення конфігурації через env | Unit | Змінні середовища перевизначають значення |
| Унікальність ULID       | Unit        | Немає колізій між горутинами          |

## 7. Організація тестів

```text
internal/
├── auth/
│   ├── jwt.go
│   ├── jwt_test.go          # unit: token gen/validation
│   ├── password.go
│   ├── password_test.go     # unit: hash/verify
│   ├── register.go
│   ├── register_test.go     # unit: use case with mocks
│   ├── login.go
│   └── login_test.go        # unit: use case with mocks
├── user/
│   ├── entity.go
│   ├── entity_test.go       # unit: domain logic (CanLogin, etc.)
│   ├── repository.go        # interfaces (mock source)
│   └── mock/
│       └── repository.go    # generated mocks
├── database/
│   ├── database.go
│   ├── database_test.go     # integration: testcontainers
│   └── mock/
│       └── database.go      # generated mock
└── testutil/
    ├── postgres.go           # testcontainer helper
    └── redis.go              # testcontainer helper
```

## 8. Швидка довідка

```mermaid
graph LR
    subgraph "Unit (fast, no Docker)"
        A[Domain Logic] --> B[testify assert/require]
        C[Use Cases] --> D[gomock mocks]
        E[HTTP Handlers] --> F[httptest.NewRecorder]
    end

    subgraph "Integration (requires Docker)"
        G[Repository] --> H[testutil.SetupPostgres]
        I[Cache] --> J[testutil.SetupRedis]
    end

    subgraph "Skip Control"
        K["testing.Short() → skip integration"]
        L["just test → -short flag"]
        M["just test-all → runs everything"]
    end
```
