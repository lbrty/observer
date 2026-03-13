---
title: Ortam Değişkenleri
weight: 1
---

## Proje

- proje adı: `observer`
- paket adı: `github.com/lbrty/observer`
- go sürümü: 1.25.\*
- varsayılan arayüz dili: Kırgızca Latin (`ky`)

## Backend ortam değişkenleri

### Server

| Değişken               | Varsayılan  | Açıklama               |
| ---------------------- | ----------- | ---------------------- |
| `SERVER_HOST`          | `localhost` | Bağlanma adresi        |
| `SERVER_PORT`          | `9000`      | Dinleme portu          |
| `SERVER_READ_TIMEOUT`  | `30s`       | HTTP okuma zaman aşımı |
| `SERVER_WRITE_TIMEOUT` | `30s`       | HTTP yazma zaman aşımı |

### Veritabanı

| Değişken       | Varsayılan | Açıklama                   |
| -------------- | ---------- | -------------------------- |
| `DATABASE_DSN` | `""`       | PostgreSQL bağlantı dizesi |

### JWT

| Değişken               | Varsayılan         | Açıklama                   |
| ---------------------- | ------------------ | -------------------------- |
| `JWT_PRIVATE_KEY_PATH` | `keys/jwt_rsa`     | RSA özel anahtar yolu      |
| `JWT_PUBLIC_KEY_PATH`  | `keys/jwt_rsa.pub` | RSA genel anahtar yolu     |
| `JWT_ACCESS_TTL`       | `15m`              | Access token ömrü          |
| `JWT_REFRESH_TTL`      | `168h`             | Refresh token ömrü (7 gün) |
| `JWT_MFA_TEMP_TTL`     | `5m`               | MFA bekleyen token ömrü    |
| `JWT_ISSUER`           | `observer`         | Token issuer claim         |

### Cookie

| Değişken           | Varsayılan         | Açıklama                      |
| ------------------ | ------------------ | ----------------------------- |
| `COOKIE_DOMAIN`    | `""` (mevcut host) | Cookie alan adı               |
| `COOKIE_SECURE`    | `false`            | Üretimde (HTTPS) `true` yapın |
| `COOKIE_SAME_SITE` | `lax`              | `lax`, `strict` veya `none`   |
| `COOKIE_MAX_AGE`   | `2h`               | Cookie ömrü                   |

### CORS

| Değişken       | Varsayılan              | Açıklama                           |
| -------------- | ----------------------- | ---------------------------------- |
| `CORS_ORIGINS` | `http://localhost:5173` | Virgülle ayrılmış izinli kaynaklar |

### Depolama

| Değişken          | Varsayılan     | Açıklama                                                              |
| ----------------- | -------------- | --------------------------------------------------------------------- |
| `STORAGE_PATH`    | `data/uploads` | Yerel dosya sistemi kökü (`STORAGE_BACKEND=local` olduğunda kullanılır) |
| `STORAGE_BACKEND` | `local`        | Depolama arka ucu: `local` veya `s3`                                  |
| `S3_ENDPOINT`     | `""`           | S3 endpoint URL'i (boş = AWS varsayılanı)                             |
| `S3_BUCKET`       | `""`           | S3 bucket adı (backend `s3` olduğunda zorunlu)                        |
| `S3_REGION`       | `us-east-1`    | S3 bölgesi                                                            |
| `S3_ACCESS_KEY`   | `""`           | AWS access key (isteğe bağlı — SDK zincirine geri döner)              |
| `S3_SECRET_KEY`   | `""`           | AWS secret key (isteğe bağlı — SDK zincirine geri döner)              |

### Diğer

| Değişken                    | Varsayılan                 | Açıklama                                              |
| --------------------------- | -------------------------- | ----------------------------------------------------- |
| `DEV_MODE`                  | `false`                    | Geliştirme modunu etkinleştir                         |
| `LOG_LEVEL`                 | `info`                     | Log seviyesi                                          |
| `REDIS_URL`                 | `redis://localhost:6379/0` | Redis bağlantı URL'i                                  |
| `SWAGGER_ENABLED`           | `false`                    | `/swagger/` adresinde Swagger UI etkinleştir          |
| `RATE_LIMIT_LOGIN`          | `10`                       | Dakika başına maksimum giriş denemesi                 |
| `RATE_LIMIT_REGISTER`       | `5`                        | Dakika başına maksimum kayıt denemesi                 |
| `SENTRY_DSN`                | `""`                       | Sentry DSN (boş — Sentry devre dışı)                  |
| `SENTRY_TRACES_SAMPLE_RATE` | `0.1`                      | Sentry performans izleme örnekleme oranı              |

## Ön yüz ortam değişkenleri

| Değişken       | Varsayılan              | Açıklama                |
| -------------- | ----------------------- | ----------------------- |
| `VITE_API_URL` | `http://localhost:9000` | Backend API temel URL'i |
