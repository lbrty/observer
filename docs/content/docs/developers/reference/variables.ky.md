---
title: Çöjrö Özgörmölörü
weight: 1
---

## Proekt

- proekt aty: `observer`
- paket aty: `github.com/lbrty/observer`
- go versijasy: 1.25.\*
- defolt UI tili: Qyrğyz Latin (`ky`)

## Backend çöjrö özgörmölörü

### Server

| Özgörmö                | Default     | Taanyştyruu         |
| ---------------------- | ----------- | ------------------- |
| `SERVER_HOST`          | `localhost` | Bind darek          |
| `SERVER_PORT`          | `9000`      | Tuñdoo portu        |
| `SERVER_READ_TIMEOUT`  | `30s`       | HTTP oquu tajmautu  |
| `SERVER_WRITE_TIMEOUT` | `30s`       | HTTP cazuu tajmautu |

### Maalymat bazasy

| Özgörmö        | Default | Taanyştyruu                  |
| -------------- | ------- | ---------------------------- |
| `DATABASE_DSN` | `""`    | PostgreSQL bajlanyş strokasy |

### JWT

| Özgörmö                | Default            | Taanyştyruu                   |
| ---------------------- | ------------------ | ----------------------------- |
| `JWT_PRIVATE_KEY_PATH` | `keys/jwt_rsa`     | RSA ceke açqyç colu           |
| `JWT_PUBLIC_KEY_PATH`  | `keys/jwt_rsa.pub` | RSA açyq açqyç colu           |
| `JWT_ACCESS_TTL`       | `15m`              | Access token möönötü          |
| `JWT_REFRESH_TTL`      | `168h`             | Refresh token möönötü (7 kün) |
| `JWT_MFA_TEMP_TTL`     | `5m`               | MFA kütüü token möönötü       |
| `JWT_ISSUER`           | `observer`         | Token issuer claim            |

### Cookie

| Özgörmö            | Default              | Taanyştyruu                         |
| ------------------ | -------------------- | ----------------------------------- |
| `COOKIE_DOMAIN`    | `""` (ağymdağy host) | Cookie domeni                       |
| `COOKIE_SECURE`    | `false`              | Produksionda `true` qojuñuz (HTTPS) |
| `COOKIE_SAME_SITE` | `lax`                | `lax`, `strict` ce `none`           |
| `COOKIE_MAX_AGE`   | `2h`                 | Cookie möönötü                      |

### CORS

| Özgörmö        | Default                 | Taanyştyruu                                    |
| -------------- | ----------------------- | ---------------------------------------------- |
| `CORS_ORIGINS` | `http://localhost:5173` | Ütür menen bölünğön uruqsat berilgen originter |

### Başqalar

| Özgörmö           | Default          | Taanyştyruu                    |
| ----------------- | ---------------- | ------------------------------ |
| `LOG_LEVEL`       | `info`           | Log deñgeeli                   |
| `REDIS_URI`       | `localhost:6379` | Redis bajlanyş URI             |
| `SWAGGER_ENABLED` | `false`          | `/swagger/` da Swagger UI açuu |

## Frontend çöjrö özgörmölörü

| Özgörmö        | Default                 | Taanyştyruu             |
| -------------- | ----------------------- | ----------------------- |
| `VITE_API_URL` | `http://localhost:9000` | Backend API bazalyq URL |
