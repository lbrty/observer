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

### Saqtoo

| Özgörmö           | Default        | Taanyştyruu                                                     |
| ----------------- | -------------- | --------------------------------------------------------------- |
| `STORAGE_PATH`    | `data/uploads` | Lokal fajl sistemesinin tübü (`STORAGE_BACKEND=local` bolğondo) |
| `STORAGE_BACKEND` | `local`        | Saqtoo backend: `local` ce `s3`                                 |
| `S3_ENDPOINT`     | `""`           | S3 endpoint URL (boş = AWS defolt)                              |
| `S3_BUCKET`       | `""`           | S3 bucket aty (backend `s3` bolğondo mildet)                    |
| `S3_REGION`       | `us-east-1`    | S3 region                                                       |
| `S3_ACCESS_KEY`   | `""`           | AWS access key (mildet emes — SDK chain qoldoloт)               |
| `S3_SECRET_KEY`   | `""`           | AWS secret key (mildet emes — SDK chain qoldoloт)               |

### Başqalar

| Özgörmö                     | Default                    | Taanyştyruu                              |
| --------------------------- | -------------------------- | ---------------------------------------- |
| `DEV_MODE`                  | `false`                    | Öndürüştü rejimdi qoşuu                  |
| `LOG_LEVEL`                 | `info`                     | Log deñgeeli                             |
| `REDIS_URL`                 | `redis://localhost:6379/0` | Redis bajlanyş URL                       |
| `SWAGGER_ENABLED`           | `false`                    | `/swagger/` da Swagger UI açuu           |
| `RATE_LIMIT_LOGIN`          | `10`                       | Minutasyna iriştüünün maks sanasy        |
| `RATE_LIMIT_REGISTER`       | `5`                        | Minutasyna kattaluunun maks sanasy       |
| `SENTRY_DSN`                | `""`                       | Sentry DSN (boş bolso Sentry öçürülğön)  |
| `SENTRY_TRACES_SAMPLE_RATE` | `0.1`                      | Sentry ishlöö körünüştörünün ülgü dereji |

## Frontend çöjrö özgörmölörü

| Özgörmö        | Default                 | Taanyştyruu             |
| -------------- | ----------------------- | ----------------------- |
| `VITE_API_URL` | `http://localhost:9000` | Backend API bazalyq URL |
