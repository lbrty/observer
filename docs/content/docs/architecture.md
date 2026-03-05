---
title: Architecture
weight: 2
---

## High-Level Overview

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

## Dependency Flow (Clean Architecture)

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

## Repository: Interface to Implementation

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

## Use Cases: Who Depends on What

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

## HTTP Request Flow

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

## Protected Route Flow (Admin + Project RBAC)

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

## DI Container Wiring

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

## Package Import Graph

```
cmd/observer/
  --> internal/config
  --> internal/database
  --> internal/logger
  --> internal/app         (DI container)
  --> internal/server

internal/app/
  --> internal/config
  --> internal/crypto
  --> internal/database
  --> internal/domain/user       (interfaces + entities)
  --> internal/domain/auth       (interfaces + entities)
  --> internal/domain/project    (interfaces + entities)
  --> internal/domain/reference  (interfaces + entities)
  --> internal/postgres          (implementations)
  --> internal/usecase/auth
  --> internal/usecase/admin

internal/server/
  --> internal/app
  --> internal/handler
  --> internal/middleware
  --> internal/domain/user       (for user.RoleAdmin constant)

internal/handler/
  --> internal/usecase/auth      (use case types)
  --> internal/usecase/admin     (use case types)
  --> internal/domain/user       (error sentinels)
  --> internal/domain/project    (error sentinels)
  --> internal/domain/reference  (error sentinels)

internal/middleware/
  --> internal/crypto            (TokenGenerator for JWT)
  --> internal/domain/user       (Role type + constants)
  --> internal/domain/project    (PermissionLoader, Action, entities)

internal/usecase/auth/
  --> internal/domain/user       (UserRepository, CredentialsRepository, entities)
  --> internal/domain/auth       (SessionRepository, entities)
  --> internal/crypto            (PasswordHasher, TokenGenerator)

internal/usecase/admin/
  --> internal/domain/user       (UserRepository, entities)
  --> internal/domain/project    (PermissionRepository, entities)
  --> internal/domain/reference  (all reference repos, entities)

internal/postgres/
  --> internal/domain/user       (entities + errors)
  --> internal/domain/auth       (entities + errors)
  --> internal/domain/project    (entities + errors)
  --> internal/domain/reference  (entities + errors)
```
