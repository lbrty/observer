---
title: Ortam Değişkenleri
weight: 1
---

## Proje

- proje adı: `observer`
- paket adı: `github.com/lbrty/observer`
- go sürümü: 1.25.*
- varsayılan arayüz dili: Kırgızca Latin (`ky`)

## Backend ortam değişkenleri

### Server

| Değişken               | Varsayılan  | Açıklama           |
| ---------------------- | ----------- | ------------------ |
| `SERVER_HOST`          | `localhost` | Bağlanma adresi    |
| `SERVER_PORT`          | `9000`      | Dinleme portu      |
| `SERVER_READ_TIMEOUT`  | `30s`       | HTTP okuma zaman aşımı  |
| `SERVER_WRITE_TIMEOUT` | `30s`       | HTTP yazma zaman aşımı  |

### Veritabanı

| Değişken       | Varsayılan | Açıklama                     |
| -------------- | ---------- | ---------------------------- |
| `DATABASE_DSN` | `""`       | PostgreSQL bağlantı dizesi   |

### JWT

| Değişken               | Varsayılan         | Açıklama                        |
| ---------------------- | ------------------ | ------------------------------- |
| `JWT_PRIVATE_KEY_PATH` | `keys/jwt_rsa`     | RSA özel anahtar yolu           |
| `JWT_PUBLIC_KEY_PATH`  | `keys/jwt_rsa.pub` | RSA genel anahtar yolu          |
| `JWT_ACCESS_TTL`       | `15m`              | Access token ömrü               |
| `JWT_REFRESH_TTL`      | `168h`             | Refresh token ömrü (7 gün)     |
| `JWT_MFA_TEMP_TTL`     | `5m`               | MFA bekleyen token ömrü         |
| `JWT_ISSUER`           | `observer`         | Token issuer claim              |

### Cookie

| Değişken           | Varsayılan          | Açıklama                         |
| ------------------ | ------------------- | -------------------------------- |
| `COOKIE_DOMAIN`    | `""` (mevcut host)  | Cookie alan adı                  |
| `COOKIE_SECURE`    | `false`             | Üretimde (HTTPS) `true` yapın   |
| `COOKIE_SAME_SITE` | `lax`               | `lax`, `strict` veya `none`     |
| `COOKIE_MAX_AGE`   | `2h`                | Cookie ömrü                      |

### CORS

| Değişken       | Varsayılan              | Açıklama                        |
| -------------- | ----------------------- | ------------------------------- |
| `CORS_ORIGINS` | `http://localhost:5173` | Virgülle ayrılmış izinli kaynaklar |

### Diğer

| Değişken          | Varsayılan       | Açıklama                         |
| ----------------- | ---------------- | -------------------------------- |
| `LOG_LEVEL`       | `info`           | Log seviyesi                     |
| `REDIS_URI`       | `localhost:6379` | Redis bağlantı URI'si            |
| `SWAGGER_ENABLED` | `false`          | `/swagger/` adresinde Swagger UI etkinleştir |

## Ön yüz ortam değişkenleri

| Değişken       | Varsayılan              | Açıklama             |
| -------------- | ----------------------- | -------------------- |
| `VITE_API_URL` | `http://localhost:9000` | Backend API temel URL'i |
