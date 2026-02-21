# ADR-002: Implementation Plan: User & JWT Authentication with RSA Signing

| Field      | Value      |
| ---------- | ---------- |
| Status     | Accepted   |
| Date       | 2026-02-21 |
| Supersedes | —          |
| Components | auth, mfa  |

---

Important: use variables defined in: `../variables.md`

## Phase 1: Database Schema & Migrations

### 1.1 Create Migration Files Structure

```text
migrations/
├── 000002_create_users_table.up.sql
├── 000002_create_users_table.down.sql
├── 000003_create_credentials_table.up.sql
├── 000003_create_credentials_table.down.sql
├── 000004_create_mfa_configs_table.up.sql
├── 000004_create_mfa_configs_table.down.sql
├── 000005_create_sessions_table.up.sql
├── 000005_create_sessions_table.down.sql
├── 000006_create_verification_tokens_table.up.sql
└── 000006_create_verification_tokens_table.down.sql
```

### 1.2 Migration Content (PostgreSQL)

**000002_create_users_table.up.sql:**

```sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,  -- ULID as text
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('user', 'seller', 'admin')),
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_users_email UNIQUE (email),
    CONSTRAINT uq_users_phone UNIQUE (phone)
);

CREATE INDEX ix_users_email ON users(email);
CREATE INDEX ix_users_phone ON users(phone);
CREATE INDEX ix_users_role ON users(role);
```

**000002_create_users_table.down.sql:**

```sql
DROP TABLE IF EXISTS users;
```

**000003_create_credentials_table.up.sql:**

```sql
CREATE TABLE credentials (
    user_id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL,
    salt TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
);
```

**000003_create_credentials_table.down.sql:**

```sql
DROP TABLE IF EXISTS credentials;
```

**000004_create_mfa_configs_table.up.sql:**

```sql
CREATE TABLE mfa_configs (
    user_id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    method VARCHAR(10) NOT NULL CHECK (method IN ('totp', 'sms')),
    secret TEXT,  -- TOTP secret (store encrypted in Phase 2)
    phone VARCHAR(20),
    is_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
);
```

**000004_create_mfa_configs_table.down.sql:**

```sql
DROP TABLE IF EXISTS mfa_configs;
```

**000005_create_sessions_table.up.sql:**

```sql
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,  -- ULID
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL,
    user_agent TEXT,
    ip VARCHAR(45),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_sessions_refresh_token UNIQUE (refresh_token)
);

CREATE INDEX ix_sessions_user_id ON sessions(user_id);
CREATE INDEX ix_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX ix_sessions_expires_at ON sessions(expires_at);
```

**000005_create_sessions_table.down.sql:**

```sql
DROP TABLE IF EXISTS sessions;
```

**000006_create_verification_tokens_table.up.sql:**

```sql
CREATE TABLE verification_tokens (
    id TEXT PRIMARY KEY,  -- ULID
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('verification', 'password_reset')),
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_verification_tokens_token UNIQUE (token)
);

CREATE INDEX ix_verification_tokens_token ON verification_tokens(token);
CREATE INDEX ix_verification_tokens_user_id ON verification_tokens(user_id);
```

**000006_create_verification_tokens_table.down.sql:**

```sql
DROP TABLE IF EXISTS verification_tokens;
```

## Phase 2: RSA Key Setup & Configuration

### 2.1 Generate RSA Keys

```bash
# Create keys directory
mkdir -p keys

# Generate private key (4096 bits)
openssl genrsa -out keys/jwt_rsa 4096

# Generate public key
openssl rsa -in keys/jwt_rsa -pubout -out keys/jwt_rsa.pub

# Set proper permissions
chmod 600 keys/jwt_rsa
chmod 644 keys/jwt_rsa.pub
```

### 2.2 Configuration Structure

```go
// internal/config/config.go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
    Redis    RedisConfig
}

type JWTConfig struct {
    PrivateKeyPath string        `env:"JWT_PRIVATE_KEY_PATH" envDefault:"keys/jwt_rsa"`
    PublicKeyPath  string        `env:"JWT_PUBLIC_KEY_PATH" envDefault:"keys/jwt_rsa.pub"`
    AccessTTL      time.Duration `env:"JWT_ACCESS_TTL" envDefault:"15m"`
    RefreshTTL     time.Duration `env:"JWT_REFRESH_TTL" envDefault:"168h"` // 7 days
    MFATempTTL     time.Duration `env:"JWT_MFA_TEMP_TTL" envDefault:"5m"`
    Issuer         string        `env:"JWT_ISSUER" envDefault:"[project_name]"`
}
```

## Phase 3: Domain Layer Implementation

### 3.1 Update Domain Entities with ULID

```go
// internal/domain/user/entity.go
import "github.com/oklog/ulid/v2"

type Role string

const (
    RoleUser 				Role = "user"
    RoleSeller      Role = "seller"
    RoleAdmin       Role = "admin"
)

type User struct {
    ID         ulid.ULID
    Email      string
    Phone      string
    Role       Role
    IsVerified bool
    IsActive   bool
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

### 3.2 Implement ULID Generator

```go
// pkg/ulid/generator.go
package ulid

import (
    "crypto/rand"
    "github.com/oklog/ulid/v2"
    "sync"
)

var (
    entropy = ulid.Monotonic(rand.Reader, 0)
    mu      sync.Mutex
)

func New() ulid.ULID {
    mu.Lock()
    defer mu.Unlock()
    return ulid.MustNew(ulid.Now(), entropy)
}

func Parse(s string) (ulid.ULID, error) {
    return ulid.Parse(s)
}
```

### 3.3 Repository Interfaces (updated signatures)

```go
// internal/domain/user/repository.go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id ulid.ULID) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    GetByPhone(ctx context.Context, phone string) (*User, error)
    Update(ctx context.Context, user *User) error
    UpdateVerified(ctx context.Context, id ulid.ULID, verified bool) error
}
```

## Phase 4: Infrastructure Layer - RSA JWT Implementation

### 4.1 RSA Key Loader

```go
// internal/infrastructure/auth/rsa_keys.go
package auth

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "os"
)

type RSAKeys struct {
    PrivateKey *rsa.PrivateKey
    PublicKey  *rsa.PublicKey
}

func LoadRSAKeys(privateKeyPath, publicKeyPath string) (*RSAKeys, error) {
    // Load private key
    privateKeyData, err := os.ReadFile(privateKeyPath)
    if err != nil {
        return nil, fmt.Errorf("read private key: %w", err)
    }

    privateBlock, _ := pem.Decode(privateKeyData)
    if privateBlock == nil {
        return nil, fmt.Errorf("failed to decode private key PEM")
    }

    privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
    if err != nil {
        return nil, fmt.Errorf("parse private key: %w", err)
    }

    // Load public key
    publicKeyData, err := os.ReadFile(publicKeyPath)
    if err != nil {
        return nil, fmt.Errorf("read public key: %w", err)
    }

    publicBlock, _ := pem.Decode(publicKeyData)
    if publicBlock == nil {
        return nil, fmt.Errorf("failed to decode public key PEM")
    }

    publicKeyInterface, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
    if err != nil {
        return nil, fmt.Errorf("parse public key: %w", err)
    }

    publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
    if !ok {
        return nil, fmt.Errorf("not an RSA public key")
    }

    return &RSAKeys{
        PrivateKey: privateKey,
        PublicKey:  publicKey,
    }, nil
}
```

### 4.2 RSA JWT Token Generator

```go
// internal/infrastructure/auth/jwt_rsa.go
package auth

import (
    "crypto/rsa"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/oklog/ulid/v2"
)

type Claims struct {
    jwt.RegisteredClaims
    UserID string `json:"uid"`
    Role   string `json:"role"`
    Type   string `json:"type"` // access, refresh, mfa_pending
}

type RSATokenGenerator struct {
    privateKey *rsa.PrivateKey
    publicKey  *rsa.PublicKey
    accessTTL  time.Duration
    refreshTTL time.Duration
    mfaTTL     time.Duration
    issuer     string
}

func NewRSATokenGenerator(keys *RSAKeys, accessTTL, refreshTTL, mfaTTL time.Duration, issuer string) *RSATokenGenerator {
    return &RSATokenGenerator{
        privateKey: keys.PrivateKey,
        publicKey:  keys.PublicKey,
        accessTTL:  accessTTL,
        refreshTTL: refreshTTL,
        mfaTTL:     mfaTTL,
        issuer:     issuer,
    }
}

func (g *RSATokenGenerator) GenerateAccessToken(userID ulid.ULID, role string) (string, time.Time, error) {
    now := time.Now().UTC()
    expiresAt := now.Add(g.accessTTL)

    claims := Claims{
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    g.issuer,
            Subject:   userID.String(),
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
        },
        UserID: userID.String(),
        Role:   role,
        Type:   "access",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    tokenString, err := token.SignedString(g.privateKey)
    if err != nil {
        return "", time.Time{}, fmt.Errorf("sign token: %w", err)
    }

    return tokenString, expiresAt, nil
}

func (g *RSATokenGenerator) GenerateRefreshToken() (string, error) {
    // Generate a secure random ULID for refresh token
    return ulid.New().String(), nil
}

func (g *RSATokenGenerator) GenerateMFAToken(userID ulid.ULID) (string, error) {
    now := time.Now().UTC()
    expiresAt := now.Add(g.mfaTTL)

    claims := Claims{
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    g.issuer,
            Subject:   userID.String(),
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(now),
        },
        UserID: userID.String(),
        Type:   "mfa_pending",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    return token.SignedString(g.privateKey)
}

func (g *RSATokenGenerator) ValidateAccessToken(tokenString string) (*Claims, error) {
    return g.validateToken(tokenString, "access")
}

func (g *RSATokenGenerator) ValidateMFAToken(tokenString string) (*Claims, error) {
    return g.validateToken(tokenString, "mfa_pending")
}

func (g *RSATokenGenerator) validateToken(tokenString, expectedType string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return g.publicKey, nil
    })

    if err != nil {
        return nil, fmt.Errorf("parse token: %w", err)
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, fmt.Errorf("invalid token claims")
    }

    if claims.Type != expectedType {
        return nil, fmt.Errorf("invalid token type: expected %s, got %s", expectedType, claims.Type)
    }

    return claims, nil
}
```

### 4.3 Password Hasher (Argon2id)

```go
// internal/infrastructure/auth/password.go
package auth

import (
    "crypto/rand"
    "crypto/subtle"
    "encoding/base64"
    "fmt"

    "golang.org/x/crypto/argon2"
)

type PasswordHasher interface {
    Hash(password string) (hash, salt string, err error)
    Verify(password, hash, salt string) error
}

type ArgonHasher struct {
    time    uint32
    memory  uint32
    threads uint8
    keyLen  uint32
}

func NewArgonHasher() *ArgonHasher {
    return &ArgonHasher{
        time:    1,
        memory:  64 * 1024, // 64 MB
        threads: 4,
        keyLen:  32,
    }
}

func (h *ArgonHasher) Hash(password string) (string, string, error) {
    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return "", "", fmt.Errorf("generate salt: %w", err)
    }

    hash := argon2.IDKey([]byte(password), salt, h.time, h.memory, h.threads, h.keyLen)

    return base64.RawStdEncoding.EncodeToString(hash),
           base64.RawStdEncoding.EncodeToString(salt),
           nil
}

func (h *ArgonHasher) Verify(password, hashStr, saltStr string) error {
    salt, err := base64.RawStdEncoding.DecodeString(saltStr)
    if err != nil {
        return fmt.Errorf("decode salt: %w", err)
    }

    expectedHash, err := base64.RawStdEncoding.DecodeString(hashStr)
    if err != nil {
        return fmt.Errorf("decode hash: %w", err)
    }

    computedHash := argon2.IDKey([]byte(password), salt, h.time, h.memory, h.threads, h.keyLen)

    if subtle.ConstantTimeCompare(computedHash, expectedHash) == 1 {
        return nil
    }

    return fmt.Errorf("invalid password")
}
```

## Phase 5: Repository Implementation

### 5.1 PostgreSQL User Repository

```go
// internal/infrastructure/persistence/postgres/user_repository.go
package postgres

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

    "github.com/oklog/ulid/v2"
    "[package_name]/internal/domain/user"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
    query := `
        INSERT INTO users (id, email, phone, role, is_verified, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

    _, err := r.db.ExecContext(ctx, query,
        u.ID.String(),
        u.Email,
        u.Phone,
        u.Role,
        u.IsVerified,
        u.IsActive,
        u.CreatedAt.UTC(),
        u.UpdatedAt.UTC(),
    )

    if err != nil {
        // Check for unique constraint violations
        // Return appropriate domain errors
        return fmt.Errorf("create user: %w", err)
    }

    return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
    query := `
        SELECT id, email, phone, role, is_verified, is_active, created_at, updated_at
        FROM users
        WHERE email = $1
    `

    var u user.User
    var idStr string
    var createdAt, updatedAt time.Time

    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &idStr,
        &u.Email,
        &u.Phone,
        &u.Role,
        &u.IsVerified,
        &u.IsActive,
        &createdAt,
        &updatedAt,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, user.ErrUserNotFound
        }
        return nil, fmt.Errorf("get user by email: %w", err)
    }

    u.ID, err = ulid.Parse(idStr)
    if err != nil {
        return nil, fmt.Errorf("parse user ID: %w", err)
    }

    u.CreatedAt = createdAt.UTC()
    u.UpdatedAt = updatedAt.UTC()

    return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id ulid.ULID) (*user.User, error) {
    query := `
        SELECT id, email, phone, role, is_verified, is_active, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    var u user.User
    var idStr string
    var createdAt, updatedAt time.Time

    err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
        &idStr,
        &u.Email,
        &u.Phone,
        &u.Role,
        &u.IsVerified,
        &u.IsActive,
        &createdAt,
        &updatedAt,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, user.ErrUserNotFound
        }
        return nil, fmt.Errorf("get user by id: %w", err)
    }

    u.ID, err = ulid.Parse(idStr)
    if err != nil {
        return nil, fmt.Errorf("parse user ID: %w", err)
    }

    u.CreatedAt = createdAt.UTC()
    u.UpdatedAt = updatedAt.UTC()

    return &u, nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
    query := `
        SELECT id, email, phone, role, is_verified, is_active, created_at, updated_at
        FROM users
        WHERE phone = $1
    `

    var u user.User
    var idStr string
    var createdAt, updatedAt time.Time

    err := r.db.QueryRowContext(ctx, query, phone).Scan(
        &idStr,
        &u.Email,
        &u.Phone,
        &u.Role,
        &u.IsVerified,
        &u.IsActive,
        &createdAt,
        &updatedAt,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, user.ErrUserNotFound
        }
        return nil, fmt.Errorf("get user by phone: %w", err)
    }

    u.ID, err = ulid.Parse(idStr)
    if err != nil {
        return nil, fmt.Errorf("parse user ID: %w", err)
    }

    u.CreatedAt = createdAt.UTC()
    u.UpdatedAt = updatedAt.UTC()

    return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
    query := `
        UPDATE users
        SET email = $2, phone = $3, role = $4, is_verified = $5, is_active = $6, updated_at = $7
        WHERE id = $1
    `

    result, err := r.db.ExecContext(ctx, query,
        u.ID.String(),
        u.Email,
        u.Phone,
        u.Role,
        u.IsVerified,
        u.IsActive,
        time.Now().UTC(),
    )

    if err != nil {
        return fmt.Errorf("update user: %w", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("get rows affected: %w", err)
    }

    if rows == 0 {
        return user.ErrUserNotFound
    }

    return nil
}

func (r *UserRepository) UpdateVerified(ctx context.Context, id ulid.ULID, verified bool) error {
    query := `
        UPDATE users
        SET is_verified = $2, updated_at = $3
        WHERE id = $1
    `

    result, err := r.db.ExecContext(ctx, query, id.String(), verified, time.Now().UTC())
    if err != nil {
        return fmt.Errorf("update verified status: %w", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("get rows affected: %w", err)
    }

    if rows == 0 {
        return user.ErrUserNotFound
    }

    return nil
}
```

## Phase 6: Application Layer - Use Cases

### 6.1 Register Use Case

```go
// internal/application/auth/register_usecase.go
package auth

import (
    "context"
    "fmt"
    "time"

    "github.com/oklog/ulid/v2"
    "[package_name]/internal/domain/user"
    pkgulid "[package_name]/pkg/ulid"
)

type RegisterUseCase struct {
    userRepo  user.UserRepository
    credRepo  user.CredentialsRepository
    tokenRepo user.VerificationTokenRepository
    hasher    PasswordHasher
    // emailSvc  notification.EmailService (Phase 2)
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
    // Validate input
    if err := uc.validateInput(input); err != nil {
        return nil, fmt.Errorf("invalid input: %w", err)
    }

    // Check if email exists
    if _, err := uc.userRepo.GetByEmail(ctx, input.Email); err == nil {
        return nil, user.ErrEmailExists
    }

    // Hash password
    hash, salt, err := uc.hasher.Hash(input.Password)
    if err != nil {
        return nil, fmt.Errorf("hash password: %w", err)
    }

    // Create user
    userID := pkgulid.New()
    now := time.Now().UTC()

    newUser := &user.User{
        ID:         userID,
        Email:      input.Email,
        Phone:      input.Phone,
        Role:       user.Role(input.Role),
        IsVerified: false,
        IsActive:   true,
        CreatedAt:  now,
        UpdatedAt:  now,
    }

    if err := uc.userRepo.Create(ctx, newUser); err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }

    // Create credentials
    cred := &user.Credentials{
        UserID:       userID,
        PasswordHash: hash,
        Salt:         salt,
        UpdatedAt:    now,
    }

    if err := uc.credRepo.Create(ctx, cred); err != nil {
        return nil, fmt.Errorf("create credentials: %w", err)
    }

    // Generate verification token (implement later)
    // Send verification email (implement later)

    return &RegisterOutput{
        UserID:  userID,
        Message: "Registration successful. Please check your email to verify your account.",
    }, nil
}
```

### 6.2 Login Use Case

```go
// internal/application/auth/login_usecase.go
package auth

import (
    "context"
    "fmt"
    "time"

    "github.com/oklog/ulid/v2"
    "[package_name]/internal/domain/auth"
    "[package_name]/internal/domain/user"
    pkgulid "[package_name]/pkg/ulid"
)

type LoginUseCase struct {
    userRepo    user.UserRepository
    credRepo    user.CredentialsRepository
    sessionRepo auth.SessionRepository
    mfaRepo     user.MFARepository
    hasher      PasswordHasher
    tokenGen    TokenGenerator
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput, userAgent, ip string) (*LoginOutput, error) {
    // Get user by email
    u, err := uc.userRepo.GetByEmail(ctx, input.Email)
    if err != nil {
        return nil, user.ErrInvalidCredentials
    }

    // Check if user can login
    if err := u.CanLogin(); err != nil {
        return nil, err
    }

    // Get credentials
    cred, err := uc.credRepo.GetByUserID(ctx, u.ID)
    if err != nil {
        return nil, user.ErrInvalidCredentials
    }

    // Verify password
    if err := uc.hasher.Verify(input.Password, cred.PasswordHash, cred.Salt); err != nil {
        return nil, user.ErrInvalidCredentials
    }

    // Check MFA status
    mfaConfig, err := uc.mfaRepo.GetByUserID(ctx, u.ID)
    if err == nil && mfaConfig.IsEnabled {
        // MFA required - generate temporary token
        mfaToken, err := uc.tokenGen.GenerateMFAToken(u.ID)
        if err != nil {
            return nil, fmt.Errorf("generate mfa token: %w", err)
        }

        return &LoginOutput{
            RequiresMFA: true,
            MFAToken:    mfaToken,
        }, nil
    }

    // No MFA - create session and return tokens
    tokens, err := uc.createSession(ctx, u, userAgent, ip)
    if err != nil {
        return nil, fmt.Errorf("create session: %w", err)
    }

    return &LoginOutput{
        RequiresMFA: false,
        Tokens:      tokens,
        User:        uc.toUserDTO(u),
    }, nil
}

func (uc *LoginUseCase) createSession(ctx context.Context, u *user.User, userAgent, ip string) (*TokenPair, error) {
    accessToken, expiresAt, err := uc.tokenGen.GenerateAccessToken(u.ID, string(u.Role))
    if err != nil {
        return nil, fmt.Errorf("generate access token: %w", err)
    }

    refreshToken, err := uc.tokenGen.GenerateRefreshToken()
    if err != nil {
        return nil, fmt.Errorf("generate refresh token: %w", err)
    }

    session := &auth.Session{
        ID:           pkgulid.New(),
        UserID:       u.ID,
        RefreshToken: refreshToken,
        UserAgent:    userAgent,
        IP:           ip,
        ExpiresAt:    expiresAt.Add(7 * 24 * time.Hour), // 7 days
        CreatedAt:    time.Now().UTC(),
    }

    if err := uc.sessionRepo.Create(ctx, session); err != nil {
        return nil, fmt.Errorf("create session: %w", err)
    }

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresAt:    expiresAt,
    }, nil
}

func (uc *LoginUseCase) toUserDTO(u *user.User) *UserDTO {
    return &UserDTO{
        ID:         u.ID,
        Email:      u.Email,
        Phone:      u.Phone,
        Role:       string(u.Role),
        IsVerified: u.IsVerified,
        CreatedAt:  u.CreatedAt,
    }
}
```

## Phase 7: HTTP Handler Implementation

### 7.1 Auth Middleware

```go
// internal/interfaces/http/middleware/auth.go
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/oklog/ulid/v2"
)

type AuthMiddleware struct {
    tokenGen TokenGenerator
}

func NewAuthMiddleware(tokenGen TokenGenerator) *AuthMiddleware {
    return &AuthMiddleware{tokenGen: tokenGen}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
            c.Abort()
            return
        }

        claims, err := m.tokenGen.ValidateAccessToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        userID, err := ulid.Parse(claims.UserID)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
            c.Abort()
            return
        }

        c.Set("user_id", userID)
        c.Set("user_role", claims.Role)
        c.Next()
    }
}

func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("user_role")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"error": "role not found"})
            c.Abort()
            return
        }

        roleStr := userRole.(string)
        for _, role := range roles {
            if roleStr == role {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
        c.Abort()
    }
}
```

### 7.2 Auth Handler

```go
// internal/interfaces/http/handlers/auth_handler.go
package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/oklog/ulid/v2"
    "[package_name]/internal/application/auth"
)

type AuthHandler struct {
    registerUC  *auth.RegisterUseCase
    loginUC     *auth.LoginUseCase
    refreshUC   *auth.RefreshTokenUseCase
    logoutUC    *auth.LogoutUseCase
    // ... other use cases
}

func (h *AuthHandler) Register(c *gin.Context) {
    var input auth.RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    output, err := h.registerUC.Execute(c.Request.Context(), input)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusCreated, output)
}

func (h *AuthHandler) Login(c *gin.Context) {
    var input auth.LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userAgent := c.GetHeader("User-Agent")
    ip := c.ClientIP()

    output, err := h.loginUC.Execute(c.Request.Context(), input, userAgent, ip)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, output)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
    var input auth.RefreshTokenInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    tokens, err := h.refreshUC.Execute(c.Request.Context(), input)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Logout(c *gin.Context) {
    userID := c.MustGet("user_id").(ulid.ULID)

    var input auth.LogoutInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.logoutUC.Execute(c.Request.Context(), userID, input.RefreshToken); err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
    // Map domain errors to HTTP status codes
    // Implement error mapping logic
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
```

## Phase 8: Bootstrap & Main

### 8.1 Dependency Injection Container

```go
// internal/app/container.go
package app

import (
    "database/sql"
    "fmt"

    "[package_name]/internal/config"
    "[package_name]/internal/infrastructure/auth"
    "[package_name]/internal/infrastructure/persistence/postgres"
)

type Container struct {
    // Repositories
    UserRepo        user.UserRepository
    CredRepo        user.CredentialsRepository
    SessionRepo     auth.SessionRepository

    // Services
    PasswordHasher  auth.PasswordHasher
    TokenGenerator  auth.TokenGenerator

    // Use Cases
    RegisterUC      *auth.RegisterUseCase
    LoginUC         *auth.LoginUseCase
    RefreshTokenUC  *auth.RefreshTokenUseCase
    LogoutUC        *auth.LogoutUseCase
}

func NewContainer(cfg *config.Config, db *sql.DB) (*Container, error) {
    // Load RSA keys
    rsaKeys, err := auth.LoadRSAKeys(cfg.JWT.PrivateKeyPath, cfg.JWT.PublicKeyPath)
    if err != nil {
        return nil, fmt.Errorf("load RSA keys: %w", err)
    }

    // Initialize repositories
    userRepo := postgres.NewUserRepository(db)
    credRepo := postgres.NewCredentialsRepository(db)
    sessionRepo := postgres.NewSessionRepository(db)

    // Initialize services
    passwordHasher := auth.NewArgonHasher()
    tokenGenerator := auth.NewRSATokenGenerator(
        rsaKeys,
        cfg.JWT.AccessTTL,
        cfg.JWT.RefreshTTL,
        cfg.JWT.MFATempTTL,
        cfg.JWT.Issuer,
    )

    // Initialize use cases
    registerUC := auth.NewRegisterUseCase(userRepo, credRepo, passwordHasher)
    loginUC := auth.NewLoginUseCase(userRepo, credRepo, sessionRepo, passwordHasher, tokenGenerator)
    refreshTokenUC := auth.NewRefreshTokenUseCase(sessionRepo, tokenGenerator)
    logoutUC := auth.NewLogoutUseCase(sessionRepo)

    return &Container{
        UserRepo:       userRepo,
        CredRepo:       credRepo,
        SessionRepo:    sessionRepo,
        PasswordHasher: passwordHasher,
        TokenGenerator: tokenGenerator,
        RegisterUC:     registerUC,
        LoginUC:        loginUC,
        RefreshTokenUC: refreshTokenUC,
        LogoutUC:       logoutUC,
    }, nil
}
```

### 8.2 Main Entry Point

```go
// cmd/api/main.go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "[package_name]/internal/app"
    "[package_name]/internal/config"
    "[package_name]/internal/interfaces/http"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize database
    db, err := app.InitDatabase(cfg.Database)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Initialize container
    container, err := app.NewContainer(cfg, db)
    if err != nil {
        log.Fatalf("Failed to initialize container: %v", err)
    }

    // Initialize HTTP server
    server := http.NewServer(cfg.Server, container)

    // Start server
    go func() {
        if err := server.Start(); err != nil {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server exited")
}
```

## Phase 9: Testing

### 9.1 Unit Tests Structure

```text
internal/
├── application/auth/
│   ├── register_usecase_test.go
│   ├── login_usecase_test.go
│   └── refresh_token_usecase_test.go
├── infrastructure/auth/
│   ├── jwt_rsa_test.go
│   └── password_test.go
└── interfaces/http/handlers/
    └── auth_handler_test.go
```

### 9.2 Integration Tests

```go
// tests/integration/auth_test.go
package integration

import (
    "testing"
    "context"

    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestAuthFlow(t *testing.T) {
    ctx := context.Background()

    // Start PostgreSQL container
    pgContainer, err := postgres.Run(ctx,
        "postgres:18-alpine",
        postgres.WithDatabase("observer_test"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    if err != nil {
        t.Fatal(err)
    }
    defer pgContainer.Terminate(ctx)

    // Run migrations
    // Initialize test server
    // Test registration
    // Test login
    // Test token refresh
    // Test logout
}
```

## Phase 10: Justfile Commands Addition

Add these commands to your existing Justfile:

```justfile
# Generate RSA keys for JWT signing
generate-keys:
    #!/usr/bin/env bash
    mkdir -p keys
    echo "Generating RSA private key (4096 bits)..."
    openssl genrsa -out keys/jwt_rsa 4096
    echo "Generating RSA public key..."
    openssl rsa -in keys/jwt_rsa -pubout -out keys/jwt_rsa.pub
    echo "Setting permissions..."
    chmod 600 keys/jwt_rsa
    chmod 644 keys/jwt_rsa.pub
    echo "✓ RSA keys generated successfully in keys/ directory"
```

## Implementation Order Summary

1. **Setup** (Day 1)
   - Generate RSA keys using `just generate-keys`
   - Create migration files
   - Set up project structure

2. **Core Infrastructure** (Days 2-3)
   - Implement ULID generator
   - Implement password hasher
   - Implement RSA JWT generator
   - Create PostgreSQL repositories

3. **Domain & Application Logic** (Days 4-5)
   - Implement domain entities
   - Implement use cases
   - Create DTOs

4. **HTTP Layer** (Day 6)
   - Implement handlers
   - Create middleware
   - Set up routes

5. **Integration** (Day 7)
   - Wire up dependency injection
   - Create bootstrap logic
   - Test end-to-end flow

6. **Testing** (Days 8-9)
   - Write unit tests
   - Write integration tests
   - Manual API testing

7. **Documentation** (Day 10)
   - API documentation
   - Setup instructions
   - Deployment guide
