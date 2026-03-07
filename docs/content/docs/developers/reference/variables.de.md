---
title: Umgebungsvariablen
weight: 1
---

## Projekt

- Projektname: `observer`
- Paketname: `github.com/lbrty/observer`
- Go-Version: 1.25.*
- Standard-UI-Sprache: Kirgisisch Latein (`ky`)

## Backend-Umgebungsvariablen

### Server

| Variable               | Standard    | Beschreibung       |
| ---------------------- | ----------- | ------------------ |
| `SERVER_HOST`          | `localhost` | Bind-Adresse       |
| `SERVER_PORT`          | `9000`      | Lausch-Port        |
| `SERVER_READ_TIMEOUT`  | `30s`       | HTTP-Lese-Timeout  |
| `SERVER_WRITE_TIMEOUT` | `30s`       | HTTP-Schreib-Timeout |

### Datenbank

| Variable       | Standard | Beschreibung                     |
| -------------- | -------- | -------------------------------- |
| `DATABASE_DSN` | `""`     | PostgreSQL-Verbindungszeichenkette |

### JWT

| Variable               | Standard           | Beschreibung                        |
| ---------------------- | ------------------ | ----------------------------------- |
| `JWT_PRIVATE_KEY_PATH` | `keys/jwt_rsa`     | RSA-Private-Key-Pfad                |
| `JWT_PUBLIC_KEY_PATH`  | `keys/jwt_rsa.pub` | RSA-Public-Key-Pfad                 |
| `JWT_ACCESS_TTL`       | `15m`              | Lebensdauer des Access-Tokens       |
| `JWT_REFRESH_TTL`      | `168h`             | Lebensdauer des Refresh-Tokens (7 Tage) |
| `JWT_MFA_TEMP_TTL`     | `5m`               | Lebensdauer des MFA-Pending-Tokens  |
| `JWT_ISSUER`           | `observer`         | Token-Issuer-Claim                  |

### Cookie

| Variable           | Standard            | Beschreibung                              |
| ------------------ | ------------------- | ----------------------------------------- |
| `COOKIE_DOMAIN`    | `""` (aktueller Host) | Cookie-Domain                           |
| `COOKIE_SECURE`    | `false`             | Auf `true` setzen in Produktion (HTTPS)   |
| `COOKIE_SAME_SITE` | `lax`               | `lax`, `strict` oder `none`               |
| `COOKIE_MAX_AGE`   | `2h`                | Cookie-Lebensdauer                        |

### CORS

| Variable       | Standard                | Beschreibung                          |
| -------------- | ----------------------- | ------------------------------------- |
| `CORS_ORIGINS` | `http://localhost:5173` | Kommagetrennte erlaubte Origins       |

### Sonstiges

| Variable          | Standard         | Beschreibung                              |
| ----------------- | ---------------- | ----------------------------------------- |
| `LOG_LEVEL`       | `info`           | Log-Level                                 |
| `REDIS_URI`       | `localhost:6379` | Redis-Verbindungs-URI                     |
| `SWAGGER_ENABLED` | `false`          | Swagger UI unter `/swagger/` aktivieren   |

## Frontend-Umgebungsvariablen

| Variable       | Standard                | Beschreibung             |
| -------------- | ----------------------- | ------------------------ |
| `VITE_API_URL` | `http://localhost:9000` | Backend-API-Basis-URL    |
