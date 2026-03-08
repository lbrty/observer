---
title: Довідник CLI
weight: 4
---

Observer надає CLI для керування сервером, міграціями бази даних, генерацією ключів, адмініструванням користувачів та утилітами для розробки.

## Встановлення

```bash
# Build from source
just build

# Or install directly
go install github.com/lbrty/observer/cmd/observer@latest
```

## Команди

### serve

Запустити HTTP-сервер.

```bash
observer serve [flags]
```

| Flag     | Type   | Default     | Опис                                              |
| -------- | ------ | ----------- | ------------------------------------------------- |
| `--host` | string | `localhost` | Хост сервера (перевизначає змінну `SERVER_HOST`)  |
| `--port` | int    | `9000`      | Порт сервера (перевизначає змінну `SERVER_PORT`)  |

**Приклади:**

```bash
# Start with defaults (localhost:9000)
observer serve

# Custom host and port
observer serve --host 0.0.0.0 --port 8080

# With environment configuration
DATABASE_DSN="postgres://..." REDIS_URL="redis://..." observer serve
```

У production-збірках вбудовані міграції застосовуються автоматично при запуску. Сервер завершує роботу плавно при отриманні SIGINT/SIGTERM з таймаутом 30 секунд.

---

### migrate

Керування міграціями бази даних.

#### migrate up

Застосувати всі незастосовані міграції.

```bash
observer migrate up [flags]
```

| Flag     | Type   | Default      | Опис                          |
| -------- | ------ | ------------ | ----------------------------- |
| `--path` | string | `migrations` | Шлях до директорії міграцій   |

```bash
# Apply all pending migrations
observer migrate up

# Use a custom migrations directory
observer migrate up --path ./db/migrations
```

#### migrate create

Створити новий файл міграції (лише вперед).

```bash
observer migrate create [name] [flags]
```

| Flag     | Type   | Default      | Опис                          |
| -------- | ------ | ------------ | ----------------------------- |
| `--path` | string | `migrations` | Шлях до директорії міграцій   |
| `--seq`  | uint   | auto         | Явний порядковий номер        |

```bash
# Create a migration (auto-numbered)
observer migrate create add_audit_log

# Create with explicit sequence number
observer migrate create add_audit_log --seq 25
```

#### migrate version

Показати поточну версію міграції.

```bash
observer migrate version [flags]
```

| Flag     | Type   | Default      | Опис                          |
| -------- | ------ | ------------ | ----------------------------- |
| `--path` | string | `migrations` | Шлях до директорії міграцій   |

---

### keygen

Згенерувати пару RSA-ключів для підпису JWT.

```bash
observer keygen [flags]
```

| Flag       | Type   | Default | Опис                                  |
| ---------- | ------ | ------- | ------------------------------------- |
| `--bits`   | int    | `4096`  | Розмір RSA-ключа (мінімум 4096)       |
| `--output` | string | `.`     | Директорія для збереження файлів      |

**Приклади:**

```bash
# Generate keys in the current directory
observer keygen

# Generate 8192-bit keys in the keys/ directory
observer keygen --bits 8192 --output keys
```

Вихідні файли: `private_key.pem` (0600) та `public_key.pem` (0644).

---

### create-admin

Створити обліковий запис адміністратора платформи.

```bash
observer create-admin [flags]
```

| Flag           | Type   | Required | Опис                                 |
| -------------- | ------ | -------- | ------------------------------------ |
| `--email`      | string | yes      | Email адміністратора                 |
| `--password`   | string | yes      | Пароль (мінімум 8 символів)          |
| `--first-name` | string | no       | Ім'я                                 |
| `--last-name`  | string | no       | Прізвище                             |
| `--phone`      | string | no       | Номер телефону                       |

**Приклади:**

```bash
# Create an admin with required fields
observer create-admin --email admin@example.com --password "s3cure-p4ss"

# With optional profile fields
observer create-admin \
  --email admin@example.com \
  --password "s3cure-p4ss" \
  --first-name Admin \
  --last-name User \
  --phone "+1234567890"
```

Підключається до бази даних, хешує пароль за допомогою Argon2id та додає користувача з роллю адміністратора (верифікований + активний). Відхиляє дублікати email та номерів телефонів.

---

### seed

Заповнити базу даних реалістичними тестовими даними для розробки.

```bash
observer seed [flags]
```

| Flag         | Type  | Default | Опис                                |
| ------------ | ----- | ------- | ----------------------------------- |
| `--people`   | int   | `50`    | Кількість людей на проєкт          |
| `--projects` | int   | `2`     | Кількість проєктів                 |
| `--seed`     | int64 | `0`     | Початкове значення (0 = випадкове)  |

**Приклади:**

```bash
# Seed with defaults (2 projects, 50 people each)
observer seed

# Custom counts
observer seed --projects 5 --people 200

# Reproducible seed
observer seed --seed 42
```

{{% callout type="warning" %}}
Ця команда очищає ВСІ таблиці перед додаванням даних. Не запускайте її на production базі даних.
{{% /callout %}}

Створює довідникові дані (країни, області, населені пункти, офіси, категорії), користувачів з відомими паролями (`password`), проєкти з дозволами, а також заповнює дані про людей із записами підтримки, міграційними записами, нотатками, тваринами та домогосподарствами.

---

### setup

Запустити інтерактивне початкове налаштування проєкту.

```bash
observer setup
```

Ця команда:

1. Створює файл `.env` зі стандартними налаштуваннями (запитує підтвердження перед перезаписом)
2. Створює необхідні директорії (`keys/`, `data/uploads/`)
3. Генерує 4096-бітну пару RSA-ключів для підпису JWT
4. Виводить інструкції щодо наступних кроків

**Приклад виводу:**

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

## Типові робочі процеси

### Перше налаштування

```bash
# 1. Run interactive setup (creates .env, keys, directories)
observer setup

# 2. Start Postgres and Redis
docker compose up -d

# 3. Run migrations
observer migrate up

# 4. Create an admin user
observer create-admin --email admin@example.com --password "your-password"

# 5. Start the server
observer serve
```

Для повного демонстраційного налаштування з прикладами даних дивіться [посібник з демо-налаштування](/docs/guide/demo-setup/).

### Додавання нової міграції

```bash
# Create the migration file
observer migrate create add_audit_log

# Edit the generated SQL file
$EDITOR migrations/000025_add_audit_log.up.sql

# Apply it
observer migrate up
```

### Заповнення даними для розробки

```bash
# Make sure migrations are applied first
observer migrate up

# Seed with defaults
observer seed

# Login with: admin@example.com / password
```

### Генерація нових JWT-ключів

```bash
# Generate new keys
observer keygen --output keys

# Update .env if paths differ from defaults
JWT_PRIVATE_KEY_PATH=keys/private_key.pem
JWT_PUBLIC_KEY_PATH=keys/public_key.pem

# Restart the server — existing tokens will be invalidated
observer serve
```
