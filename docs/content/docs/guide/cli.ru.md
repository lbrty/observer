---
title: Справочник CLI
weight: 4
---

Observer предоставляет CLI для управления сервером, миграциями базы данных, генерацией ключей, администрированием пользователей и утилитами для разработки.

## Установка

```bash
# Build from source
just build

# Or install directly
go install github.com/lbrty/observer/cmd/observer@latest
```

## Команды

### serve

Запуск HTTP-сервера.

```bash
observer serve [flags]
```

| Flag     | Type   | Default     | Описание                                          |
| -------- | ------ | ----------- | ------------------------------------------------- |
| `--host` | string | `localhost` | Хост сервера (переопределяет переменную `SERVER_HOST`) |
| `--port` | int    | `9000`      | Порт сервера (переопределяет переменную `SERVER_PORT`) |

**Примеры:**

```bash
# Start with defaults (localhost:9000)
observer serve

# Custom host and port
observer serve --host 0.0.0.0 --port 8080

# With environment configuration
DATABASE_DSN="postgres://..." REDIS_URL="redis://..." observer serve
```

В production-сборках встроенные миграции применяются автоматически при запуске. Сервер корректно завершает работу по сигналу SIGINT/SIGTERM с таймаутом 30 секунд.

---

### migrate

Управление миграциями базы данных.

#### migrate up

Применить все ожидающие миграции.

```bash
observer migrate up [flags]
```

| Flag     | Type   | Default      | Описание                        |
| -------- | ------ | ------------ | ------------------------------- |
| `--path` | string | `migrations` | Путь к директории с миграциями  |

```bash
# Apply all pending migrations
observer migrate up

# Use a custom migrations directory
observer migrate up --path ./db/migrations
```

#### migrate create

Создать новый файл однонаправленной миграции.

```bash
observer migrate create [name] [flags]
```

| Flag     | Type   | Default      | Описание                        |
| -------- | ------ | ------------ | ------------------------------- |
| `--path` | string | `migrations` | Путь к директории с миграциями  |
| `--seq`  | uint   | auto         | Явный порядковый номер          |

```bash
# Create a migration (auto-numbered)
observer migrate create add_audit_log

# Create with explicit sequence number
observer migrate create add_audit_log --seq 25
```

#### migrate version

Показать текущую версию миграции.

```bash
observer migrate version [flags]
```

| Flag     | Type   | Default      | Описание                        |
| -------- | ------ | ------------ | ------------------------------- |
| `--path` | string | `migrations` | Путь к директории с миграциями  |

---

### keygen

Генерация пары RSA-ключей для подписи JWT.

```bash
observer keygen [flags]
```

| Flag       | Type   | Default | Описание                                |
| ---------- | ------ | ------- | --------------------------------------- |
| `--bits`   | int    | `4096`  | Размер RSA-ключа (минимум 4096)         |
| `--output` | string | `.`     | Директория для сохранения файлов ключей |

**Примеры:**

```bash
# Generate keys in the current directory
observer keygen

# Generate 8192-bit keys in the keys/ directory
observer keygen --bits 8192 --output keys
```

Выходные файлы: `private_key.pem` (0600) и `public_key.pem` (0644).

---

### create-admin

Создание учётной записи администратора платформы.

```bash
observer create-admin [flags]
```

| Flag           | Type   | Required | Описание                         |
| -------------- | ------ | -------- | -------------------------------- |
| `--email`      | string | yes      | Email администратора             |
| `--password`   | string | yes      | Пароль (минимум 8 символов)      |
| `--first-name` | string | no       | Имя                              |
| `--last-name`  | string | no       | Фамилия                          |
| `--phone`      | string | no       | Номер телефона                   |

**Примеры:**

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

Подключается к базе данных, хеширует пароль с помощью Argon2id и создаёт пользователя с ролью администратора (подтверждённый и активный). Отклоняет дублирующиеся email-адреса и номера телефонов.

---

### seed

Заполнение базы данных реалистичными тестовыми данными для разработки.

```bash
observer seed [flags]
```

| Flag         | Type  | Default | Описание                              |
| ------------ | ----- | ------- | ------------------------------------- |
| `--people`   | int   | `50`    | Количество людей на проект            |
| `--projects` | int   | `2`     | Количество проектов                   |
| `--seed`     | int64 | `0`     | Зерно генератора (0 = случайное)      |

**Примеры:**

```bash
# Seed with defaults (2 projects, 50 people each)
observer seed

# Custom counts
observer seed --projects 5 --people 200

# Reproducible seed
observer seed --seed 42
```

{{% callout type="warning" %}}
Эта команда очищает ВСЕ таблицы перед вставкой данных. Не запускайте её на production-базе данных.
{{% /callout %}}

Создаёт справочные данные (страны, области, населённые пункты, офисы, категории), пользователей с известными паролями (`password`), проекты с правами доступа, а также заполняет записи о людях с записями поддержки, миграционными записями, заметками, домашними животными и домохозяйствами.

---

### setup

Интерактивная первоначальная настройка проекта.

```bash
observer setup
```

Эта команда:

1. Создаёт файл `.env` с настройками по умолчанию (запрашивает подтверждение перед перезаписью)
2. Создаёт необходимые директории (`keys/`, `data/uploads/`)
3. Генерирует 4096-битную пару RSA-ключей для подписи JWT
4. Выводит инструкции по дальнейшим действиям

**Пример вывода:**

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

## Типичные сценарии использования

### Первоначальная настройка

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

Для полной демонстрационной настройки с примерами данных смотрите [руководство по демо-настройке](/docs/guide/demo-setup/).

### Добавление новой миграции

```bash
# Create the migration file
observer migrate create add_audit_log

# Edit the generated SQL file
$EDITOR migrations/000025_add_audit_log.up.sql

# Apply it
observer migrate up
```

### Заполнение данными для разработки

```bash
# Make sure migrations are applied first
observer migrate up

# Seed with defaults
observer seed

# Login with: admin@example.com / password
```

### Генерация новых JWT-ключей

```bash
# Generate new keys
observer keygen --output keys

# Update .env if paths differ from defaults
JWT_PRIVATE_KEY_PATH=keys/private_key.pem
JWT_PUBLIC_KEY_PATH=keys/public_key.pem

# Restart the server — existing tokens will be invalidated
observer serve
```
