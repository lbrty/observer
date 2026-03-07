---
title: Test
weight: 4
---

[package_name] üç test katmanı kullanır: birim testleri, testcontainers ile entegrasyon testleri ve HTTP üzerinden uçtan uca testler.

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

## Testleri Çalıştırma

```bash
just test              # yalnızca birim testleri (hızlı, Docker gerektirmez)
just test-all          # entegrasyon dahil tüm testler
just test-coverage     # HTML kapsam raporu oluştur
just test-race         # yarış durumlarını tespit et
just bench             # karşılaştırma testleri
just generate-mocks    # mock dosyalarını yeniden oluştur
```

## Test Dosyası Kuralları

| Desen             | Amaç                                 |
| ----------------- | ------------------------------------- |
| `*_test.go`       | Aynı paketteki test dosyası          |
| `testing.Short()` | `just test`'te entegrasyon testlerini atla |
| `t.Helper()`      | Fonksiyonları test yardımcısı olarak işaretle |
| `testutil.Setup*` | Entegrasyon testleri için konteyner kurulumu |

## 1. Testify ile Birim Testleri

Ölümcül olmayan kontroller için `assert` (test devam eder), ölümcül kontroller için `require` (test durur) kullanın.

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

### Yaygın Testify Doğrulamaları

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

### Tablo Güdümlü Testler

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

## 2. gomock ile Mock Kullanımı

Mock'lar, `go:generate` direktifleri kullanılarak arayüzlerden oluşturulur.

### Mock Oluşturma

`internal/user/repository.go` dosyasındaki direktif:

```go
//go:generate mockgen -destination=mock/repository.go -package=mock [package_name]/internal/user UserRepository,CredentialsRepository,SessionRepository,MFARepository,VerificationTokenRepository
```

Arayüz değişikliklerinden sonra tüm mock'ları yeniden oluşturmak için `just generate-mocks` komutunu çalıştırın.

### Testlerde Mock Kullanımı

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

### gomock Eşleştiricileri

```go
gomock.Any()                          // match any value
gomock.Eq("exact")                    // exact match
gomock.Not(gomock.Eq("excluded"))     // negation
gomock.Nil()                          // match nil
```

### Çağrı Sırası Beklentisi

```go
first := mockRepo.EXPECT().GetByEmail(gomock.Any(), "a@b.com").Return(nil, errNotFound)
mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).After(first)
```

## 3. Testcontainers ile Entegrasyon Testleri

Entegrasyon testleri, Postgres ve Redis için gerçek Docker konteynerleri başlatır.

### Test Yapısı

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

### Yardımcı: testutil.SetupPostgres

`internal/testutil/postgres.go` dosyasında bulunur. Bir DSN ve temizleme fonksiyonu döndürür:

```go
dsn, cleanup := testutil.SetupPostgres(t)
defer cleanup()
```

Docker kullanılamadığında otomatik olarak atlanır.

### Yardımcı: testutil.SetupRedis

`internal/testutil/redis.go` dosyasında bulunur:

```go
addr, cleanup := testutil.SetupRedis(t)
defer cleanup()
```

## 4. HTTP Handler Testleri

Çalışan bir sunucu olmadan HTTP uç noktalarını test etmek için Gin router ile `httptest` kullanın.

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

### POST Uç Noktalarını Test Etme

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

## 5. Uçtan Uca Başarılı Akış: Kayıt -> Giriş -> Yenileme -> Çıkış

Bu, tam kimlik doğrulama akışının HTTP üzerinden uçtan uca test edilmesini gösterir.

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

## 6. Kritik Test Kontrol Listesi

### Auth Domain

| Test                          | Tür  | Neyi doğrular                        |
| ----------------------------- | ---- | ------------------------------------ |
| Geçerli verilerle kayıt       | Birim | Kullanıcı + kimlik bilgileri oluşturuldu |
| Tekrar eden e-posta ile kayıt | Birim | `ErrEmailExists` döner               |
| Tekrar eden telefon ile kayıt | Birim | `ErrPhoneExists` döner               |
| Geçersiz rol ile kayıt        | Birim | `ErrInvalidRole` döner               |
| Geçerli kimlik bilgileriyle giriş | Birim | Token çifti döner                |
| Yanlış parola ile giriş       | Birim | `ErrInvalidCredentials` döner        |
| Aktif olmayan kullanıcı girişi | Birim | `ErrUserNotActive` döner            |
| Geçerli token yenileme        | Birim | Eski oturum silinir, yeni çift verilir |
| Süresi dolmuş oturum yenileme | Birim | `ErrSessionExpired` döner            |
| Token oluşturma + doğrulama   | Birim | Claim'ler eşleşir, süre dolumu çalışır |
| Token türü uyumsuzluğu        | Birim | Access token MFA olarak reddedilir   |
| Parola hash benzersizliği     | Birim | Aynı parola -> farklı hash'ler       |

### Altyapı

| Test                  | Tür         | Neyi doğrular                       |
| --------------------- | ----------- | ----------------------------------- |
| DB bağlantı + ping    | Entegrasyon | Testcontainer Postgres çalışır      |
| Health uç noktası     | Birim       | 200 `{"status":"ok"}` döner         |
| Request ID middleware  | Birim       | X-Request-ID başlığı 26 karakterlik ULID |
| Düzgün kapanma        | Birim       | Sunucu hatasız kapanır              |
| Config varsayılanları  | Birim       | Anlamlı varsayılanlar yüklenir      |
| Config env geçersiz kılma | Birim   | Ortam değişkenleri varsayılanları geçersiz kılar |
| ULID benzersizliği     | Birim       | Goroutine'ler arasında çakışma yok  |

## 7. Test Organizasyonu

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

## 8. Hızlı Referans

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
