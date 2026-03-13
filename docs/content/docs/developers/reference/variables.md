---
title: Environment Variables
weight: 1
---

## Project

- project name: `observer`
- package name: `github.com/lbrty/observer`
- go version: 1.25.\*
- default UI language: Kyrgyz Latin (`ky`)

## Backend environment variables

### Server

| Variable               | Default     | Description        |
| ---------------------- | ----------- | ------------------ |
| `SERVER_HOST`          | `localhost` | Bind address       |
| `SERVER_PORT`          | `9000`      | Listen port        |
| `SERVER_READ_TIMEOUT`  | `30s`       | HTTP read timeout  |
| `SERVER_WRITE_TIMEOUT` | `30s`       | HTTP write timeout |

### Database

| Variable       | Default | Description                  |
| -------------- | ------- | ---------------------------- |
| `DATABASE_DSN` | `""`    | PostgreSQL connection string |

### JWT

| Variable               | Default            | Description                     |
| ---------------------- | ------------------ | ------------------------------- |
| `JWT_PRIVATE_KEY_PATH` | `keys/jwt_rsa`     | RSA private key path            |
| `JWT_PUBLIC_KEY_PATH`  | `keys/jwt_rsa.pub` | RSA public key path             |
| `JWT_ACCESS_TTL`       | `15m`              | Access token lifetime           |
| `JWT_REFRESH_TTL`      | `168h`             | Refresh token lifetime (7 days) |
| `JWT_MFA_TEMP_TTL`     | `5m`               | MFA pending token lifetime      |
| `JWT_ISSUER`           | `observer`         | Token issuer claim              |

### Cookie

| Variable           | Default             | Description                      |
| ------------------ | ------------------- | -------------------------------- |
| `COOKIE_DOMAIN`    | `""` (current host) | Cookie domain                    |
| `COOKIE_SECURE`    | `false`             | Set `true` in production (HTTPS) |
| `COOKIE_SAME_SITE` | `lax`               | `lax`, `strict`, or `none`       |
| `COOKIE_MAX_AGE`   | `2h`                | Cookie lifetime                  |

### CORS

| Variable       | Default                 | Description                     |
| -------------- | ----------------------- | ------------------------------- |
| `CORS_ORIGINS` | `http://localhost:5173` | Comma-separated allowed origins |

### Storage

| Variable          | Default        | Description                                         |
| ----------------- | -------------- | --------------------------------------------------- |
| `STORAGE_PATH`    | `data/uploads` | Local filesystem root (used when `STORAGE_BACKEND=local`) |
| `STORAGE_BACKEND` | `local`        | Storage backend: `local` or `s3`                    |
| `S3_ENDPOINT`     | `""`           | S3 endpoint URL (empty = AWS default)               |
| `S3_BUCKET`       | `""`           | S3 bucket name (required when backend is `s3`)      |
| `S3_REGION`       | `us-east-1`    | S3 region                                           |
| `S3_ACCESS_KEY`   | `""`           | AWS access key (optional — falls back to SDK chain) |
| `S3_SECRET_KEY`   | `""`           | AWS secret key (optional — falls back to SDK chain) |

### Other

| Variable                    | Default          | Description                              |
| --------------------------- | ---------------- | ---------------------------------------- |
| `DEV_MODE`                  | `false`          | Enable development mode                  |
| `LOG_LEVEL`                 | `info`           | Log level                                |
| `REDIS_URL`                 | `redis://localhost:6379/0` | Redis connection URL          |
| `SWAGGER_ENABLED`           | `false`          | Enable Swagger UI at `/swagger/`         |
| `RATE_LIMIT_LOGIN`          | `10`             | Max login attempts per minute            |
| `RATE_LIMIT_REGISTER`       | `5`              | Max registration attempts per minute     |
| `SENTRY_DSN`                | `""`             | Sentry DSN (empty disables Sentry)       |
| `SENTRY_TRACES_SAMPLE_RATE` | `0.1`            | Sentry performance traces sample rate    |

## Frontend environment variables

| Variable       | Default                 | Description          |
| -------------- | ----------------------- | -------------------- |
| `VITE_API_URL` | `http://localhost:9000` | Backend API base URL |
