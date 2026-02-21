# Observer

This project aims to support forcibly displaced individuals and animals in various ways,
including legal consulting and humanitarian support. It advocates for compassion and
understanding towards those affected by displacement and works towards creating a more
just and equitable world for everyone.

## Example case

To get more context, here is a real-world example (people and places are made up).

Let's say Lee had to flee his home country because of war. Now he is in Ukraine. He found
an NGO and activist group helping with accommodation and paperwork. All displaced people
and their private information is stored in spreadsheets without any protection — and it
grew large. Activities became hard to perform, so more spreadsheets were created. Tracking
provided support is stored in the same spreadsheets, with no automatic system to aggregate
data. Staff has to manually gather information about any person case by case.

Observer replaces that chaos with a secure, structured platform.

## Features

- **Encryption of sensitive information** — Personal information and documents are encrypted and protected.
- **Human & pet friendly** — People and pets can be registered and receive all sorts of support.
- **Two-factor authentication** — TOTP-based MFA with backup codes.
- **Audit logs** — Tracks actions users perform to support auditing and capture important changes.
- **Invite-only mode** — Restrict registration to invited or known users only.
- **Role-based permissions** — Fine-grained access control per project.

## Tech stack

Observer is built using the following major frameworks and libraries:

- [Poetry](https://python-poetry.org/) — dependency and environment management
- [FastAPI](https://fastapi.tiangolo.com/) — API server
- [SQLAlchemy](https://www.sqlalchemy.org/) — interact with Postgres
- [Alembic](https://alembic.sqlalchemy.org/en/latest/) — database schema migrations
- [Typer](https://typer.tiangolo.com/) — CLI utilities
- [Pydantic](https://docs.pydantic.dev/) — entity validation and settings
- [Python Cryptography](https://github.com/pyca/cryptography) — encryption/decryption
- [Pytest](https://docs.pytest.org/en/latest/) — testing
- [Ruff](https://github.com/charliermarsh/ruff) — linting
- [Black](https://black.readthedocs.io/en/stable/) — code formatting
- [Postgres](https://www.postgresql.org/) — main database

## Development setup

Requires Python 3.10+.

1. Install [Poetry](https://python-poetry.org/)
2. Install dependencies: `poetry install`
3. Generate encryption keys (see [Encryption keys](#encryption-keys))
4. Install and start Postgres

Minimal `.env` for local development:

```sh
DB_URI=postgresql+asyncpg://postgres:postgres@localhost:5432/observer
STORAGE_KIND=fs
KEYSTORE_PATH=/ABS/PATH/TO/keys
STORAGE_ROOT=/ABS/PATH/TO/uploads
```

### Postgres with Docker

```yml
version: "3.1"

services:
  db:
    image: postgres:14.5
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "127.0.0.1:5432:5432"

volumes:
  postgres_data: {}
```

### Running

Apply migrations then start the server:

```sh
python -m observer db upgrade
python -m observer server start --port 3000
# OR
make serve
```

### Running tests

```sh
make test
```

### Generating OpenAPI schema

```sh
make swagger
```

## CLI

Observer ships a built-in CLI powered by [Typer](https://github.com/tiangolo/typer):

```sh
python -m observer --help

╭─ Commands ───────────────╮
│ db                       │
│ keys                     │
│ server                   │
│ swagger                  │
╰──────────────────────────╯
```

### Migrations

```sh
# Generate a new migration
python -m observer db revision -m "migration name"

# Apply migrations
python -m observer db upgrade
```

## Deployment

Before deploying, decide and configure the following:

1. **Domain** — `APP_DOMAIN` used to build links in emails
2. **Mailer** — `MAILER_TYPE`: `gmail`, `sendgrid`, or `dummy`
3. **Storage** — `STORAGE_KIND`: `fs` or `s3`
4. **Database** — `DB_URI` pointing to Postgres 12+
5. **Encryption keys** — generate and place keys before starting

Run from source or use the Docker image [`sultaniman/observer`](https://hub.docker.com/r/sultaniman/observer).

### Storage configuration

| Variable          | Description                              |
| ----------------- | ---------------------------------------- |
| `STORAGE_KIND`    | `fs` or `s3`                             |
| `STORAGE_ROOT`    | Absolute path (fs) or S3 bucket key (s3) |
| `DOCUMENTS_PATH`  | Relative to `STORAGE_ROOT`               |
| `KEYSTORE_PATH`   | Relative to `STORAGE_ROOT`               |
| `MAX_UPLOAD_SIZE` | Defaults to 5 MB                         |

For S3, keys and documents must reside in the same bucket under different paths:

```text
/storage/
    - keys/
    - documents/
```

### Example `.env`

```sh
DB_URI=postgresql+asyncpg://postgres:postgres@localhost:5432/observer
STORAGE_KIND=fs
STORAGE_ROOT=/uploads/
DOCUMENTS_PATH=documents
KEYSTORE_PATH=keys
```

## Configuration reference

| Variable                          | Default              | Description                          |
| --------------------------------- | -------------------- | ------------------------------------ |
| `DEBUG`                           | `false`              | Enable FastAPI debug mode            |
| `PORT`                            | `3000`               | Server port                          |
| `DB_URI`                          | —                    | Postgres DSN                         |
| `POOL_SIZE`                       | `5`                  | DB connection pool size              |
| `MAX_OVERFLOW`                    | `10`                 | Max overflow connections             |
| `POOL_TIMEOUT`                    | `30`                 | Pool timeout in seconds              |
| `ECHO`                            | `false`              | Echo SQL queries                     |
| `APP_DOMAIN`                      | `observer.app`       | Frontend domain for email links      |
| `INVITE_ONLY`                     | `false`              | Enable invite-only mode              |
| `ADMIN_EMAILS`                    | `admin@examples.com` | Comma-separated admin emails         |
| `KEYSTORE_PATH`                   | `keys`               | Path to keystore folder              |
| `KEY_SIZE`                        | `2048`               | RSA key size                         |
| `AES_KEY_BITS`                    | `32`                 | AES key bits                         |
| `ACCESS_TOKEN_EXPIRATION_MINUTES` | `15`                 | Access token TTL                     |
| `REFRESH_TOKEN_EXPIRATION_DAYS`   | `180`                | Refresh token TTL                    |
| `TOTP_LEEWAY`                     | `10`                 | OTP validation leeway in seconds     |
| `NUM_BACKUP_CODES`                | `6`                  | Number of MFA backup codes           |
| `CORS_ORIGINS`                    | `["*"]`              | Allowed CORS origins                 |
| `GZIP_LEVEL`                      | `8`                  | Gzip compression level               |
| `MAILER_TYPE`                     | `dummy`              | Mailer: `gmail`, `sendgrid`, `dummy` |
| `FROM_EMAIL`                      | `no-reply@email.com` | Sender email address                 |
| `STORAGE_KIND`                    | `fs`                 | Storage backend: `fs` or `s3`        |
| `MAX_UPLOAD_SIZE`                 | `5242880`            | Max upload size in bytes             |
| `AUDIT_EVENT_EXPIRATION_DAYS`     | `365`                | Audit log retention in days          |
| `LOGIN_EVENT_EXPIRATION_DAYS`     | `7`                  | Login event retention in days        |

### Mailer backends

**Gmail:**

- `GMAIL_USERNAME`
- `GMAIL_PASSWORD`
- `GMAIL_PORT` (default: `465`)
- `GMAIL_HOSTNAME` (default: `smtp.gmail.com`)

**Sendgrid:**

- `SENDGRID_API_KEY`

### S3 storage

| Variable      | Example                              |
| ------------- | ------------------------------------ |
| `S3_ENDPOINT` | `https://s3.aws.amazon.com/observer` |
| `S3_REGION`   | `eu-central-1`                       |
| `S3_BUCKET`   | `observer-keys`                      |

## Encryption

Observer uses a hybrid encryption approach:

- **Personal information** — encrypted with RSA private keys
- **Documents** — AES symmetric encryption; the AES secret is itself encrypted with an RSA key

### Encryption keys

RSA keys are used to:

1. Encrypt personal information
2. Issue and validate JWT access/refresh tokens
3. Encrypt AES secrets used for document encryption

Generate a key with `openssl`:

```sh
openssl genrsa -out key.pem 2048
```

Or use the built-in CLI:

```sh
python -m observer keys generate --size 2048 -o key.pem
```

## Sessions

Observer uses stateless JWT session management with two cookies:

| Cookie          | Expiration | Notes                                      |
| --------------- | ---------- | ------------------------------------------ |
| `access_token`  | 15 minutes | Short-lived                                |
| `refresh_token` | 180 days   | HTTP-only; used to issue new access tokens |

Passwords are hashed with `passlib` + `bcrypt`. MFA uses `totp`.

## Two-factor authentication (TOTP)

Login flow with MFA:

1. Submit email and password
2. If credentials are valid and MFA is enabled, server returns `HTTP 417`
3. Client presents TOTP code alongside credentials
4. If TOTP is valid, auth tokens are issued

To enable MFA:

1. `POST /mfa/enable` — generates TOTP secret and QR code
2. User scans QR code and enters the TOTP code
3. `POST /mfa/setup` — confirms setup, returns encrypted backup codes

## Invite-only mode

Restrict the system to invited users only:

```sh
INVITE_ONLY=true
ADMIN_EMAILS=admin@examples.com,admin-staff@examples.com
INVITE_EXPIRATION_MINUTES=15
INVITE_URL=/account/invites/{code}
```

## Roles and permissions

Permissions are bound to projects and assigned per user.

### Roles

| Role         | Description                                      |
| ------------ | ------------------------------------------------ |
| `admin`      | Full access                                      |
| `consultant` | Full access                                      |
| `staff`      | Create, read, update; no delete or personal info |
| `guest`      | Read only                                        |

### Granting permissions

**Invite a new user with permissions:**

```json
POST /admin/invites

{
  "email": "user@example.com",
  "role": "staff",
  "permissions": [
    {
      "can_create": false,
      "can_read": true,
      "can_update": false,
      "can_delete": false,
      "can_read_documents": false,
      "can_read_personal_info": false,
      "can_invite_members": false,
      "project_id": "405d8375-3514-403b-8c43-83ae74cfe0e9"
    }
  ]
}
```

**Add a project member:**

```json
POST /projects/{project_id}/members

{
  "can_create": false,
  "can_read": true,
  "can_update": false,
  "can_delete": false,
  "can_read_documents": false,
  "can_read_personal_info": false,
  "can_invite_members": false,
  "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5",
  "project_id": "405d8375-3514-403b-8c43-83ae74cfe0e9"
}
```

## Audit logs

Audit logs track:

- Important events and database changes
- History of actions (create, update, delete, login, password reset)

Each log entry includes a `ref` string and optional `data` payload:

```ini
endpoint=create_place,action=create:place,place_id=11111111-1111-1111-1111-111111111111,ref_id=523608c6-23ff-421b-a0ed-a1aec17dadc6
```

Logs can be configured to expire after a set number of days via `AUDIT_EVENT_EXPIRATION_DAYS`.

## Allowed document formats

Uploads default to a 5 MB limit (`MAX_UPLOAD_SIZE`). Supported formats:

`jpg`, `png`, `csv`, `txt`, `md`, `doc`, `docx`, `xls`, `xlsx`, `rtf`, `pdf`, `mp4`, `mpeg`, `mp3`

## Supporting the project

If you believe in this project, you can support it by contributing, donating, or spreading the word.

- ADA: `addr1qx9p6vjp88fnmyxz7krcs860ufh6scyflclfs9e9m5gh7fcpype89hv85xg228tv8pgndzy4wjawxu723ttn9kfycwsswp5u2u`
- ETH: `0x75F774c9583820bC72aFC53844B02656089fd17f`
- [PayPal](https://paypal.me/SultanIman)
- [GitHub Sponsors](https://github.com/lbrty/observer)
