---
title: CLI Şilteme
weight: 4
---

Observer serverdi, maалymatter bazasyn, açqyç generatsijasyn, qoldonuuçu başqaruusun cana iştep çyğaruu quraldaryn başqaruu üçün CLI beret.

## Ornotuu

```bash
# Baştapqy koddor curuu
just build

# Ce tüzdön-tüz ornotuu
go install github.com/lbrty/observer/cmd/observer@latest
```

## Bujruqtar

### serve

HTTP serverdi iştetüü.

```bash
observer serve [flags]
```

| Flag     | Type   | Default     | Süröttömö                                          |
| -------- | ------ | ----------- | -------------------------------------------------- |
| `--host` | string | `localhost` | Server hostu (`SERVER_HOST` env özgörtöt)          |
| `--port` | int    | `9000`      | Server portu (`SERVER_PORT` env özgörtöt)          |

**Misaldar:**

```bash
# Demejki cöndöölör menen iştetüü (localhost:9000)
observer serve

# Başqa host cana port menen
observer serve --host 0.0.0.0 --port 8080

# Çöjrö özgörmölörü menen
DATABASE_DSN="postgres://..." REDIS_URL="redis://..." observer serve
```

Prodakşn qurulmalarda qamtylğan migratsijalar iştetilgende avtomattyq türdö qoldondulat. Server SIGINT/SIGTERM signaldarynda 30 sekunddyq tajmaut menen toqtojt.

---

### migrate

Maalymatter bazasynyñ migratsijasyn başqaruu.

#### migrate up

Bardyq kütülgön migratsijalardy qoldonuu.

```bash
observer migrate up [flags]
```

| Flag     | Type   | Default      | Süröttömö                       |
| -------- | ------ | ------------ | ------------------------------- |
| `--path` | string | `migrations` | Migratsijalar qataloqunun colu  |

```bash
# Bardyq kütülgön migratsijalardy qoldonuu
observer migrate up

# Başqa migratsijalar qataloqun qoldonuu
observer migrate up --path ./db/migrations
```

#### migrate create

Cañy aldyğa ğana migratsija fajlyn tüzüü.

```bash
observer migrate create [name] [flags]
```

| Flag     | Type   | Default      | Süröttömö                       |
| -------- | ------ | ------------ | ------------------------------- |
| `--path` | string | `migrations` | Migratsijalar qataloqunun colu  |
| `--seq`  | uint   | auto         | Açyq yraattululuq nomeri        |

```bash
# Migratsija tüzüü (avtomattyq nomerlengen)
observer migrate create add_audit_log

# Açyq yraattululuq nomeri menen tüzüü
observer migrate create add_audit_log --seq 25
```

#### migrate version

Uçurdağy migratsija versijasyn körsötüü.

```bash
observer migrate version [flags]
```

| Flag     | Type   | Default      | Süröttömö                       |
| -------- | ------ | ------------ | ------------------------------- |
| `--path` | string | `migrations` | Migratsijalar qataloqunun colu  |

---

### keygen

JWT qol qojuu üçün RSA açqyç cupttaryn generatsijaaloo.

```bash
observer keygen [flags]
```

| Flag       | Type   | Default | Süröttömö                              |
| ---------- | ------ | ------- | -------------------------------------- |
| `--bits`   | int    | `4096`  | RSA açqyç ölçömü (minimum 4096)       |
| `--output` | string | `.`     | Açqyç fajldary üçün çyğaruu qataloqu  |

**Misaldar:**

```bash
# Uçurdağy qatalogdo açqyçtardy generatsijaaloo
observer keygen

# keys/ qatalogunda 8192-bit açqyçtardy generatsijaaloo
observer keygen --bits 8192 --output keys
```

Çyğaruu fajldary: `private_key.pem` (0600) cana `public_key.pem` (0644).

---

### create-admin

Platforma administrator akkauntun tüzüü.

```bash
observer create-admin [flags]
```

| Flag           | Type   | Required | Süröttömö                          |
| -------------- | ------ | -------- | ---------------------------------- |
| `--email`      | string | ooba     | Administrator email                |
| `--password`   | string | ooba     | Administrator syrsözü (min 8 belgi)|
| `--first-name` | string | coq      | Aty                                |
| `--last-name`  | string | coq      | Familijasy                         |
| `--phone`      | string | coq      | Telefon nomeri                     |

**Misaldar:**

```bash
# Mildetüü talaalar menen administrator tüzüü
observer create-admin --email admin@example.com --password "s3cure-p4ss"

# Qoşumça profil taalary menen
observer create-admin \
  --email admin@example.com \
  --password "s3cure-p4ss" \
  --first-name Admin \
  --last-name User \
  --phone "+1234567890"
```

Maalymatter bazasyna tutaşat, syrsözdü Argon2id menen heştejt cana qoldonuuçunu administrator rolu menen qoşot (tekşerilgen + aktivdüü). Qajtalanğan emailderdi cana telefon nomerlerin çetke qağat.

---

### seed

İştep çyğaruu üçün maalymatter bazasyn realduu test maалymattary menen toltuuruu.

```bash
observer seed [flags]
```

| Flag         | Type  | Default | Süröttömö                               |
| ------------ | ----- | ------- | --------------------------------------- |
| `--people`   | int   | `50`    | Ar bir proekt üçün adamdardyn sany      |
| `--projects` | int   | `2`     | Proektterdiin sany                      |
| `--seed`     | int64 | `0`     | Tuş keldi seed (0 = tuş keldi)          |

**Misaldar:**

```bash
# Demejki cöndöölör menen toltuuruu (2 proekt, ar birine 50 adam)
observer seed

# Başqa sandar menen
observer seed --projects 5 --people 200

# Qajtalanma seed
observer seed --seed 42
```

{{% callout type="warning" %}}
Bul bujruq maalymattardy qoşuudan murda BARDYQ tablitsalardy tazalajt. Prodakşn maalymatter bazasynda iştetpeñiz.
{{% /callout %}}

Şilteme maalymattaryn (ölkölör, ştattar, cerler, keñseler, kategorijalar), belgilüü syrsözdörü bar qoldonuuçulardy (`password`), uruqsattary bar proektterdi tüzöt cana adamdardy qoldoo cazuulary, migratsija cazuulary, eskertüülör, üj canybarlar cana üj çarbalary menen tolturat.

---

### setup

Birinçi colqu proekt ornotuusun interaktivdüü iştetüü.

```bash
observer setup
```

Bul bujruq:

1. Demejki cöndöölör menen `.env` fajlyn tüzöt (üstünön cazuudan murda surajt)
2. Kerektüü qatalogdordu tüzöt (`keys/`, `data/uploads/`)
3. JWT qol qojuu üçün 4096-bit RSA açqyç cupttaryn generatsijaalajt
4. Kijinki qadamdardy nusqamalaryn körsötöt

**Çyğaruu misaly:**

```
Created .env with default configuration.
Created directory: keys
Created directory: data/uploads
Generating 4096-bit RSA key pair...
Private key written to: keys/private_key.pem
Public key written to: keys/public_key.pem

Setup complete!

Next steps:
  1. Start Postgres and Redis:
     docker compose up -d

  2. Run database migrations:
     observer migrate up

  3. Create an admin user:
     observer create-admin --email admin@example.com --password "your-password"

  4. Start the server:
     observer serve
```

---

## Calpy iş protsesster

### Birinçi colqu ornotuu

```bash
# 1. İnteraktivdüü ornotuunu iştetüü (.env, açqyçtar, qatalogdor tüzöt)
observer setup

# 2. Postgres cana Redis iştetüü
docker compose up -d

# 3. Migratsijalardy qoldonuu
observer migrate up

# 4. Administrator qoldonuuçu tüzüü
observer create-admin --email admin@example.com --password "your-password"

# 5. Serverdi iştetüü
observer serve
```

Toluq demo ornotuu üçün [demo ornotuu qoldonmosun](/docs/guide/demo-setup/) qarañyz.

### Cañy migratsija qoşuu

```bash
# Migratsija fajlyn tüzüü
observer migrate create add_audit_log

# Generatsijaalanğan SQL fajldy oñdoo
$EDITOR migrations/000025_add_audit_log.up.sql

# Qoldonuu
observer migrate up
```

### İştep çyğaruu maalymattaryn toltuuruu

```bash
# Adegene migratsijalar qoldondulğanyn tekşeriñiz
observer migrate up

# Demejki cöndöölör menen toltuuruu
observer seed

# Kirüü: admin@example.com / password
```

### Cañy JWT açqyçtaryn generatsijaaloo

```bash
# Cañy açqyçtardy generatsijaaloo
observer keygen --output keys

# Coldor demejkiden ajyrmadansa .env cañyrtuu
JWT_PRIVATE_KEY_PATH=keys/private_key.pem
JWT_PUBLIC_KEY_PATH=keys/public_key.pem

# Serverdi qajra iştetüü — uçurdağy tokender caraqsyz bolot
observer serve
```
