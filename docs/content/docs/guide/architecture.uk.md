---
title: Архітектура
weight: 3
---

Ця сторінка призначена для розробників та технічного персоналу, які хочуть зрозуміти, як побудований Observer. Якщо ви адміністратор, що налаштовує Observer для вашої організації, можете пропустити це — перейдіть до [Розгортання](/docs/guide/deployment/).

## Загальний огляд

Кожний HTTP-запит проходить однаковий шлях: він надходить на сервер, проходить через middleware (автентифікація, логування), потрапляє до обробника, який делегує виконання use case, а use case взаємодіє з базою даних через репозиторій. Конфігурація та впровадження залежностей пов'язують все разом при запуску.

```mermaid
graph TD
    CLIENT[HTTP Client] --> SERVER[Server<br>internal/server]
    SERVER --> MW[Middleware<br>internal/middleware]
    MW --> HANDLER[Handlers<br>internal/handler]
    HANDLER --> USECASE[Use Cases<br>internal/usecase]
    USECASE --> IFACE[Repository Interfaces<br>internal/domain/*/repository.go]
    IFACE -.implements.-> IMPL[Repository Implementations<br>internal/postgres]
    IMPL --> DB[(PostgreSQL)]
    USECASE --> CRYPTO[Crypto<br>internal/crypto]
    MW --> CRYPTO
    APP[DI Container<br>internal/app] -.wires.-> SERVER
    APP -.wires.-> HANDLER
    APP -.wires.-> USECASE
    APP -.wires.-> IMPL
    CONFIG[Config<br>internal/config] --> APP
```

## Потік залежностей (Чиста архітектура)

Кодова база організована в шари. Внутрішні шари визначають правила, зовнішні надають інфраструктуру. Залежності завжди спрямовані всередину — бізнес-логіка ніколи не імпортує код бази даних або HTTP напряму. Це дозволяє тестувати use cases без працюючої бази даних.

```mermaid
graph LR
    subgraph Outer["Outer Layer (Infrastructure)"]
        PG[internal/postgres]
        SRV[internal/server]
        CFG[internal/config]
    end

    subgraph Middle["Middle Layer (Adapters)"]
        HDL[internal/handler]
        MDW[internal/middleware]
    end

    subgraph Inner["Inner Layer (Business)"]
        UC[internal/usecase]
    end

    subgraph Core["Core (Domain)"]
        ENT[Entities]
        IFACE[Repository Interfaces]
        ERR[Domain Errors]
    end

    PG -->|implements| IFACE
    HDL -->|calls| UC
    MDW -->|uses| IFACE
    UC -->|depends on| IFACE
    UC -->|uses| ENT
    UC -->|returns| ERR
    HDL -->|maps| ERR

    style Core fill:#e8f5e9
    style Inner fill:#fff3e0
    style Middle fill:#e3f2fd
    style Outer fill:#fce4ec
```

## Репозиторій: від інтерфейсу до реалізації

Доменний код визначає _які_ операції з даними потрібні (інтерфейси), тоді як шар PostgreSQL надає _як_ (реалізації). Це розділення означає, що ви можете замінити PostgreSQL на іншу базу даних, не торкаючись жодної бізнес-логіки. Кожна доменна область — користувачі, авторизація, проєкти, довідкові дані — має власний інтерфейс репозиторію.

```mermaid
classDiagram
    direction LR

    namespace domain_user {
        class UserRepository {
            <<interface>>
            +Create(ctx, *User) error
            +GetByID(ctx, ulid.ULID) (*User, error)
            +GetByEmail(ctx, string) (*User, error)
            +GetByPhone(ctx, string) (*User, error)
            +Update(ctx, *User) error
            +UpdateVerified(ctx, ulid.ULID, bool) error
            +List(ctx, UserListFilter) ([]*User, int, error)
        }
        class CredentialsRepository {
            <<interface>>
            +Create(ctx, *Credentials) error
            +GetByUserID(ctx, ulid.ULID) (*Credentials, error)
        }
        class MFARepository {
            <<interface>>
            +Create(ctx, *MFAConfig) error
            +GetByUserID(ctx, ulid.ULID) (*MFAConfig, error)
        }
    }

    namespace domain_auth {
        class SessionRepository {
            <<interface>>
            +Create(ctx, *Session) error
            +GetByRefreshToken(ctx, string) (*Session, error)
            +Delete(ctx, ulid.ULID) error
            +DeleteByRefreshToken(ctx, string) error
        }
    }

    namespace domain_project {
        class PermissionLoader {
            <<interface>>
            +GetPermission(ctx, ulid.ULID, string) (*Permission, error)
            +IsProjectOwner(ctx, ulid.ULID, string) (bool, error)
        }
        class PermissionRepository {
            <<interface>>
            +List(ctx, string) ([]*ProjectPermission, error)
            +GetByID(ctx, string) (*ProjectPermission, error)
            +Create(ctx, *ProjectPermission) error
            +Update(ctx, *ProjectPermission) error
            +Delete(ctx, string) error
        }
    }

    namespace domain_reference {
        class CountryRepository {
            <<interface>>
            +List(ctx) ([]*Country, error)
            +GetByID(ctx, string) (*Country, error)
            +Create(ctx, *Country) error
            +Update(ctx, *Country) error
            +Delete(ctx, string) error
        }
        class StateRepository {
            <<interface>>
        }
        class PlaceRepository {
            <<interface>>
        }
        class OfficeRepository {
            <<interface>>
        }
        class CategoryRepository {
            <<interface>>
        }
    }

    namespace postgres {
        class pg_UserRepository {
            -db *sqlx.DB
        }
        class pg_CredentialsRepository {
            -db *sqlx.DB
        }
        class pg_SessionRepository {
            -db *sqlx.DB
        }
        class pg_MFARepository {
            -db *sqlx.DB
        }
        class pg_PermissionRepository {
            -db *sqlx.DB
        }
        class pg_ProjectPermissionRepository {
            -db *sqlx.DB
        }
        class pg_CountryRepository {
            -db *sqlx.DB
        }
    }

    pg_UserRepository ..|> UserRepository
    pg_CredentialsRepository ..|> CredentialsRepository
    pg_SessionRepository ..|> SessionRepository
    pg_MFARepository ..|> MFARepository
    pg_PermissionRepository ..|> PermissionLoader
    pg_ProjectPermissionRepository ..|> PermissionRepository
    pg_CountryRepository ..|> CountryRepository
```

## Use Cases: хто від чого залежить

Кожна дія користувача — вхід у систему, перегляд списку людей, призначення дозволів — обробляється окремим use case. Use cases координують роботу між репозиторіями та криптографічними сервісами, але самі не містять HTTP або код бази даних. Діаграма нижче показує, від яких репозиторіїв залежить кожний use case.

```mermaid
graph TD
    subgraph auth_usecases["Auth Use Cases"]
        REG[RegisterUseCase]
        LOG[LoginUseCase]
        REF[RefreshTokenUseCase]
        OUT[LogoutUseCase]
    end

    subgraph admin_usecases["Admin Use Cases"]
        LU[ListUsersUseCase]
        GU[GetUserUseCase]
        UU[UpdateUserUseCase]
    end

    subgraph perm_usecases["Permission Use Cases"]
        LP[ListPermissionsUseCase]
        AP[AssignPermissionUseCase]
        UP[UpdatePermissionUseCase]
        RP[RevokePermissionUseCase]
    end

    subgraph ref_usecases["Reference Use Cases"]
        COU[CountryUseCase]
        STA[StateUseCase]
        PLA[PlaceUseCase]
        OFF[OfficeUseCase]
        CAT[CategoryUseCase]
    end

    subgraph interfaces["Repository Interfaces"]
        UR[UserRepository]
        CR[CredentialsRepository]
        SR[SessionRepository]
        MR[MFARepository]
        PR[PermissionRepository]
        COR[CountryRepository]
        STR[StateRepository]
        PLR[PlaceRepository]
        OFR[OfficeRepository]
        CAR[CategoryRepository]
    end

    subgraph crypto["Crypto Services"]
        PH[PasswordHasher]
        TG[TokenGenerator]
    end

    REG --> UR
    REG --> CR
    REG --> PH

    LOG --> UR
    LOG --> CR
    LOG --> SR
    LOG --> MR
    LOG --> PH
    LOG --> TG

    REF --> SR
    REF --> TG

    OUT --> SR

    LU --> UR
    GU --> UR
    UU --> UR

    LP --> PR
    AP --> PR
    UP --> PR
    RP --> PR

    COU --> COR
    STA --> STR
    PLA --> PLR
    OFF --> OFR
    CAT --> CAR
```

## Потік HTTP-запиту

Ось що відбувається, коли користувач входить у систему. Запит надходить через маршрутизатор, проходить через middleware, який призначає ідентифікатор запиту та логер, потім потрапляє до обробника авторизації. Обробник розбирає JSON-тіло та викликає use case входу, який шукає користувача, перевіряє пароль за допомогою Argon2, генерує JWT-токени та створює сесію.

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router
    participant MW as Middleware
    participant H as Handler
    participant UC as UseCase
    participant Repo as Repository
    participant DB as PostgreSQL

    C->>R: POST /auth/login
    R->>MW: requestID + logger + recovery
    MW->>H: AuthHandler.Login
    H->>H: Bind JSON request
    H->>UC: LoginUseCase.Execute(input)
    UC->>Repo: UserRepo.GetByEmail(email)
    Repo->>DB: SELECT ... FROM users
    DB-->>Repo: row
    Repo-->>UC: *User
    UC->>Repo: CredRepo.GetByUserID(id)
    Repo->>DB: SELECT ... FROM credentials
    DB-->>Repo: row
    Repo-->>UC: *Credentials
    UC->>UC: Verify password (Argon2)
    UC->>UC: Generate tokens (RSA)
    UC->>Repo: SessionRepo.Create(session)
    Repo->>DB: INSERT INTO sessions
    UC-->>H: LoginOutput
    H-->>C: 200 JSON response
```

## Потік захищених маршрутів (Admin + Project RBAC)

Захищені маршрути проходять додаткові перевірки. Адміністративні маршрути перевіряють платформну роль користувача (admin, staff тощо). Маршрути на рівні проєкту завантажують дозволи користувача на рівні проєкту та перевіряють, чи достатня його проєктна роль для запитуваної дії. Middleware також встановлює прапорці конфіденційності, що контролюють, чи включає відповідь контактну інформацію, персональні дані або дані документів.

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router
    participant Auth as AuthMiddleware
    participant Role as RequireRole
    participant ProjAuth as ProjectAuthMiddleware
    participant H as Handler
    participant UC as UseCase

    C->>R: GET /admin/projects/:id/permissions
    R->>Auth: Authenticate()
    Auth->>Auth: Parse Bearer JWT
    Auth->>Auth: Set CtxUserID, CtxUserRole
    Auth-->>Role: next
    Role->>Role: Check user.Role in [admin]
    Role-->>H: next (or 403)
    H->>UC: ListPermissionsUseCase.Execute(projectID)
    UC-->>H: []PermissionDTO
    H-->>C: 200 JSON

    Note over C,R: Project-scoped route (future)
    C->>R: GET /projects/:id/people
    R->>Auth: Authenticate()
    Auth-->>ProjAuth: next
    ProjAuth->>ProjAuth: Load PermissionLoader.GetPermission()
    ProjAuth->>ProjAuth: Check role rank >= MinRoleForAction
    ProjAuth->>ProjAuth: Set CtxProjectRole + sensitivity flags
    ProjAuth-->>H: next (or 403)
```

## Зв'язування DI-контейнера

При запуску додаток зчитує конфігурацію та підключається до бази даних, потім зв'язує все разом у контейнері впровадження залежностей. Контейнер створює репозиторії, криптографічні сервіси та use cases, передаючи кожному компоненту його залежності. Повністю зібраний контейнер передається серверу, який впроваджує обробники та middleware у маршрутизатор.

```mermaid
graph TD
    subgraph inputs["Inputs"]
        CFG[Config]
        DB[Database]
    end

    subgraph container["Container wires everything"]
        direction TB
        KEYS[RSA Keys] --> TG[TokenGenerator]
        HASH[ArgonHasher]

        DB --> |sqlxDB| REPOS[All Repositories]

        REPOS --> UC_AUTH[Auth Use Cases]
        HASH --> UC_AUTH
        TG --> UC_AUTH

        REPOS --> UC_ADMIN[Admin Use Cases]
        REPOS --> UC_REF[Reference Use Cases]
        REPOS --> UC_PERM[Permission Use Cases]
    end

    subgraph output["Output"]
        CONT[Container struct]
    end

    CFG --> KEYS
    CFG --> container
    DB --> container
    container --> CONT

    CONT --> SERVER[Server.setupRoutes]
    SERVER --> |injects into| HANDLERS[Handlers]
    SERVER --> |injects into| MIDDLEWARE[Middleware]
```
