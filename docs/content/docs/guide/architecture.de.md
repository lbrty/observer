---
title: Architektur
weight: 3
---

Diese Seite richtet sich an Entwickler und technische Mitarbeiter, die verstehen möchten, wie Observer aufgebaut ist. Wenn Sie als Administrator Observer für Ihre Organisation einrichten, können Sie dies überspringen — gehen Sie stattdessen zur [Bereitstellung](/docs/guide/deployment/).

## Überblick

Jede HTTP-Anfrage durchläuft denselben Pfad: Sie erreicht den Server, passiert Middleware (Authentifizierung, Logging), gelangt zu einem Handler, der an einen Use Case delegiert, und der Use Case kommuniziert über ein Repository mit der Datenbank. Konfiguration und Dependency Injection verbinden beim Start alles miteinander.

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

## Abhängigkeitsfluss (Clean Architecture)

Die Codebasis ist in Schichten organisiert. Innere Schichten definieren die Regeln, äußere Schichten stellen die Infrastruktur bereit. Abhängigkeiten zeigen immer nach innen — Geschäftslogik importiert nie direkt Datenbank- oder HTTP-Code. Dadurch können Use Cases ohne laufende Datenbank getestet werden.

```mermaid
graph LR
    subgraph Outer["Äußere Schicht (Infrastruktur)"]
        PG[internal/postgres]
        SRV[internal/server]
        CFG[internal/config]
    end

    subgraph Middle["Mittlere Schicht (Adapter)"]
        HDL[internal/handler]
        MDW[internal/middleware]
    end

    subgraph Inner["Innere Schicht (Geschäftslogik)"]
        UC[internal/usecase]
    end

    subgraph Core["Kern (Domain)"]
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

## Repository: Interface zu Implementierung

Domain-Code definiert, _welche_ Datenoperationen benötigt werden (Interfaces), während die PostgreSQL-Schicht das _Wie_ bereitstellt (Implementierungen). Diese Trennung bedeutet, dass Sie PostgreSQL gegen eine andere Datenbank austauschen könnten, ohne Geschäftslogik zu berühren. Jeder Domain-Bereich — Benutzer, Auth, Projekte, Referenzdaten — hat sein eigenes Repository-Interface.

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

## Use Cases: Wer hängt wovon ab

Jede Benutzeraktion — Anmelden, Personen auflisten, Berechtigungen zuweisen — wird von einem dedizierten Use Case behandelt. Use Cases koordinieren zwischen Repositories und Crypto-Services, enthalten aber selbst keinen HTTP- oder Datenbankcode. Das folgende Diagramm zeigt, von welchen Repositories jeder Use Case abhängt.

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

## HTTP-Anfragefluss

Hier sehen Sie, was passiert, wenn sich ein Benutzer anmeldet. Die Anfrage kommt über den Router, passiert Middleware, die eine Request-ID und einen Logger zuweist, und erreicht dann den Auth-Handler. Der Handler parst den JSON-Body und ruft den Login-Use-Case auf, der den Benutzer nachschlägt, das Passwort mit Argon2 verifiziert, JWT-Tokens generiert und eine Session erstellt.

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

## Geschützte Routen (Admin + Project RBAC)

Geschützte Routen durchlaufen zusätzliche Prüfungen. Admin-Routen verifizieren die Plattformrolle des Benutzers (Admin, Mitarbeiter usw.). Projektbezogene Routen laden die Projektberechtigung des Benutzers und prüfen, ob seine Projektrolle für die angeforderte Aktion ausreicht. Die Middleware setzt außerdem Sensitivitätsstufen, die steuern, ob die Antwort Kontaktdaten, persönliche Details oder Dokumentdaten enthält.

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

Beim Start liest die Anwendung die Konfiguration und verbindet sich mit der Datenbank, dann wird alles in einem Dependency-Injection-Container zusammengebaut. Der Container erstellt Repositories, Crypto-Services und Use Cases und übergibt jeder Komponente ihre Abhängigkeiten. Der vollständig zusammengebaute Container wird an den Server übergeben, der Handler und Middleware in den Router einbindet.

```mermaid
graph TD
    subgraph inputs["Eingaben"]
        CFG[Config]
        DB[Database]
    end

    subgraph container["Container verbindet alles"]
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

    subgraph output["Ausgabe"]
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
