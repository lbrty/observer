# Observer

Case management platform for NGOs working with forcibly displaced people and animals. Replaces spreadsheets with a structured, auditable system for registering individuals, tracking support history, managing documents, and reporting to donors.

**Stack:** Go · PostgreSQL · Redis · React 19

---

## Quick start

### 1. Generate RSA keys

```sh
mkdir keys
docker run --rm -v "$PWD/keys:/keys" sultaniman/observer keygen \
  --private-key /keys/jwt_rsa \
  --public-key  /keys/jwt_rsa.pub
```

### 2. Run with Docker Compose

Save the following as `docker-compose.yml`:

```yaml
services:
  postgres:
    image: postgres:18-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: observer
    volumes:
      - postgres_data:/var/lib/postgresql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:8-alpine
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  observer:
    image: sultaniman/observer:latest
    ports:
      - "9000:9000"
    environment:
      DATABASE_DSN: "postgres://postgres:postgres@postgres:5432/observer?sslmode=disable"
      REDIS_URL: "redis://redis:6379/0"
      JWT_PRIVATE_KEY_PATH: "/keys/jwt_rsa"
      JWT_PUBLIC_KEY_PATH:  "/keys/jwt_rsa.pub"
      CORS_ORIGINS: "https://your-domain.example"
      COOKIE_SECURE: "true"
    volumes:
      - ./keys:/keys:ro
      - uploads:/data/uploads
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy

volumes:
  postgres_data:
  uploads:
```

```sh
docker compose up -d
```

### 3. Apply migrations

```sh
docker compose exec observer migrate up
```

### 4. Create the first admin account

```sh
docker compose exec observer create-admin admin@example.com yourpassword
```

The web UI is now available at `http://localhost:9000`.

---

## Environment variables

| Variable               | Default                 | Description                              |
| ---------------------- | ----------------------- | ---------------------------------------- |
| `DATABASE_DSN`         | —                       | Postgres connection string (required)    |
| `REDIS_URL`            | `redis://localhost:6379/0` | Redis connection URL                  |
| `JWT_PRIVATE_KEY_PATH` | `keys/jwt_rsa`          | Path to RSA private key (min 4096 bits)  |
| `JWT_PUBLIC_KEY_PATH`  | `keys/jwt_rsa.pub`      | Path to RSA public key                   |
| `SERVER_HOST`          | `0.0.0.0`               | Bind address                             |
| `SERVER_PORT`          | `9000`                  | HTTP port                                |
| `STORAGE_BACKEND`      | `local`                 | `local` or `s3`                          |
| `STORAGE_PATH`         | `data/uploads`          | Upload directory (local backend)         |
| `S3_BUCKET`            | —                       | S3 bucket name (s3 backend)              |
| `S3_REGION`            | `us-east-1`             | S3 region                                |
| `S3_ENDPOINT`          | —                       | Custom S3 endpoint (MinIO, etc.)         |
| `CORS_ORIGINS`         | `http://localhost:5173` | Comma-separated allowed origins          |
| `COOKIE_SECURE`        | `true`                  | Set `false` only for local HTTP dev      |
| `SWAGGER_ENABLED`      | `false`                 | Expose Swagger UI at `/swagger/`         |
| `SENTRY_DSN`           | —                       | Sentry error tracking DSN                |

---

## CLI commands

The container entrypoint is the `observer` binary. Available sub-commands:

| Command        | Description                                  |
| -------------- | -------------------------------------------- |
| `serve`        | Start the HTTP server (default)              |
| `migrate up`   | Apply all pending migrations                 |
| `keygen`       | Generate an RSA key pair for JWT signing     |
| `create-admin` | Create the first platform admin account      |
| `seed`         | Load demo data (development/staging only)    |

---

## File storage

By default uploads are written to `STORAGE_PATH` inside the container. Mount a volume there to persist files across restarts.

For production, set `STORAGE_BACKEND=s3` and supply `S3_BUCKET`, `S3_REGION`, and AWS credentials via the standard `AWS_*` environment variables. Any S3-compatible store (MinIO, Tigris, Cloudflare R2) works by setting `S3_ENDPOINT`.

---

## Upgrading

1. Pull the new image.
2. Run `migrate up` — migrations are forward-only and safe to run on each deploy.
3. Restart the service.

---

## Source & support

- GitHub: [github.com/lbrty/observer](https://github.com/lbrty/observer)
- Contributions and [donations](https://paypal.me/SultanIman) are welcome.
