# Variables

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

### Other

| Variable          | Default          | Description                      |
| ----------------- | ---------------- | -------------------------------- |
| `LOG_LEVEL`       | `info`           | Log level                        |
| `REDIS_URI`       | `localhost:6379` | Redis connection URI             |
| `SWAGGER_ENABLED` | `false`          | Enable Swagger UI at `/swagger/` |

## Frontend environment variables

| Variable       | Default                 | Description          |
| -------------- | ----------------------- | -------------------- |
| `VITE_API_URL` | `http://localhost:9000` | Backend API base URL |
