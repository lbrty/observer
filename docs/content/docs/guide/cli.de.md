---
title: CLI-Referenz
weight: 4
---

Observer stellt eine CLI zur Verwaltung des Servers, der Datenbankmigrationen, der Schlüsselgenerierung, der Benutzerverwaltung und der Entwicklungswerkzeuge bereit.

## Installation

```bash
# Build from source
just build

# Or install directly
go install github.com/lbrty/observer/cmd/observer@latest
```

## Befehle

### serve

HTTP-Server starten.

```bash
observer serve [flags]
```

| Flag     | Type   | Default     | Beschreibung                                       |
| -------- | ------ | ----------- | -------------------------------------------------- |
| `--host` | string | `localhost` | Server-Host (überschreibt die Umgebungsvariable `SERVER_HOST`) |
| `--port` | int    | `9000`      | Server-Port (überschreibt die Umgebungsvariable `SERVER_PORT`) |

**Beispiele:**

```bash
# Start with defaults (localhost:9000)
observer serve

# Custom host and port
observer serve --host 0.0.0.0 --port 8080

# With environment configuration
DATABASE_DSN="postgres://..." REDIS_URL="redis://..." observer serve
```

Bei Produktions-Builds werden eingebettete Migrationen beim Start automatisch angewendet. Der Server fährt bei SIGINT/SIGTERM mit einem 30-Sekunden-Timeout ordnungsgemäß herunter.

---

### migrate

Verwaltung von Datenbankmigrationen.

#### migrate up

Alle ausstehenden Migrationen anwenden.

```bash
observer migrate up [flags]
```

| Flag     | Type   | Default      | Beschreibung                  |
| -------- | ------ | ------------ | ----------------------------- |
| `--path` | string | `migrations` | Pfad zum Migrationsverzeichnis |

```bash
# Apply all pending migrations
observer migrate up

# Use a custom migrations directory
observer migrate up --path ./db/migrations
```

#### migrate create

Eine neue vorwärtsgerichtete Migrationsdatei erstellen.

```bash
observer migrate create [name] [flags]
```

| Flag     | Type   | Default      | Beschreibung                     |
| -------- | ------ | ------------ | -------------------------------- |
| `--path` | string | `migrations` | Pfad zum Migrationsverzeichnis   |
| `--seq`  | uint   | auto         | Explizite Sequenznummer          |

```bash
# Create a migration (auto-numbered)
observer migrate create add_audit_log

# Create with explicit sequence number
observer migrate create add_audit_log --seq 25
```

#### migrate version

Aktuelle Migrationsversion anzeigen.

```bash
observer migrate version [flags]
```

| Flag     | Type   | Default      | Beschreibung                  |
| -------- | ------ | ------------ | ----------------------------- |
| `--path` | string | `migrations` | Pfad zum Migrationsverzeichnis |

---

### keygen

RSA-Schlüsselpaar für JWT-Signierung generieren.

```bash
observer keygen [flags]
```

| Flag       | Type   | Default | Beschreibung                          |
| ---------- | ------ | ------- | ------------------------------------- |
| `--bits`   | int    | `4096`  | RSA-Schlüsselgröße (mindestens 4096) |
| `--output` | string | `.`     | Ausgabeverzeichnis für Schlüsseldateien |

**Beispiele:**

```bash
# Generate keys in the current directory
observer keygen

# Generate 8192-bit keys in the keys/ directory
observer keygen --bits 8192 --output keys
```

Ausgabedateien: `private_key.pem` (0600) und `public_key.pem` (0644).

---

### create-admin

Ein Plattform-Administratorkonto erstellen.

```bash
observer create-admin [flags]
```

| Flag           | Type   | Erforderlich | Beschreibung                         |
| -------------- | ------ | ------------ | ------------------------------------ |
| `--email`      | string | ja           | Admin-E-Mail                         |
| `--password`   | string | ja           | Admin-Passwort (min. 8 Zeichen)      |
| `--first-name` | string | nein         | Vorname                              |
| `--last-name`  | string | nein         | Nachname                             |
| `--phone`      | string | nein         | Telefonnummer                        |

**Beispiele:**

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

Stellt eine Verbindung zur Datenbank her, hasht das Passwort mit Argon2id und fügt den Benutzer mit der Admin-Rolle ein (verifiziert + aktiv). Doppelte E-Mail-Adressen und Telefonnummern werden abgelehnt.

---

### seed

Datenbank mit realistischen Testdaten für die Entwicklung befüllen.

```bash
observer seed [flags]
```

| Flag         | Type  | Default | Beschreibung                       |
| ------------ | ----- | ------- | ---------------------------------- |
| `--people`   | int   | `50`    | Anzahl der Personen pro Projekt    |
| `--projects` | int   | `2`     | Anzahl der Projekte                |
| `--seed`     | int64 | `0`     | Zufallsseed (0 = zufällig)        |

**Beispiele:**

```bash
# Seed with defaults (2 projects, 50 people each)
observer seed

# Custom counts
observer seed --projects 5 --people 200

# Reproducible seed
observer seed --seed 42
```

{{% callout type="warning" %}}
Dieser Befehl leert ALLE Tabellen, bevor Daten eingefügt werden. Nicht gegen eine Produktionsdatenbank ausführen.
{{% /callout %}}

Erstellt Referenzdaten (Länder, Regionen, Orte, Büros, Kategorien), Benutzer mit bekannten Passwörtern (`password`), Projekte mit Berechtigungen und befüllt Personen mit Unterstützungsdatensätzen, Migrationsdatensätzen, Notizen, Haustieren und Haushalten.

---

### setup

Erstmalige Projekteinrichtung interaktiv durchführen.

```bash
observer setup
```

Dieser Befehl:

1. Erstellt eine `.env`-Datei mit sinnvollen Standardwerten (fragt vor dem Überschreiben nach)
2. Erstellt erforderliche Verzeichnisse (`keys/`, `data/uploads/`)
3. Generiert ein 4096-Bit-RSA-Schlüsselpaar für die JWT-Signierung
4. Gibt Anweisungen für die nächsten Schritte aus

**Beispielausgabe:**

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

## Häufige Arbeitsabläufe

### Ersteinrichtung

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

Eine vollständige Demo-Einrichtung mit Beispieldaten finden Sie in der [Demo-Einrichtungsanleitung](/docs/guide/demo-setup/).

### Neue Migration hinzufügen

```bash
# Create the migration file
observer migrate create add_audit_log

# Edit the generated SQL file
$EDITOR migrations/000025_add_audit_log.up.sql

# Apply it
observer migrate up
```

### Entwicklungsdaten befüllen

```bash
# Make sure migrations are applied first
observer migrate up

# Seed with defaults
observer seed

# Login with: admin@example.com / password
```

### Neue JWT-Schlüssel generieren

```bash
# Generate new keys
observer keygen --output keys

# Update .env if paths differ from defaults
JWT_PRIVATE_KEY_PATH=keys/private_key.pem
JWT_PUBLIC_KEY_PATH=keys/public_key.pem

# Restart the server — existing tokens will be invalidated
observer serve
```
