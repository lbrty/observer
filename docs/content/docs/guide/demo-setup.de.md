---
title: Demo-Einrichtung
weight: 3
---

Probieren Sie Observer mit einer vollständig befüllten Datenbank in unter 5 Minuten aus. Der `seed`-Befehl füllt die Datenbank mit realistischen Testdaten — Personen, Projekte, Unterstützungseinträge, Migrationsverläufe, Haushalte und mehr — so können Sie alle Funktionen erkunden, ohne Daten manuell eingeben zu müssen.

## Schnellstart

```bash
git clone https://github.com/lbrty/observer.git
cd observer
just build
./observer setup
docker compose up -d
./observer migrate up
./observer seed
./observer serve
```

Öffnen Sie `http://localhost:9000/health` — Sie sollten `{"status":"healthy"}` sehen.

Starten Sie die Weboberfläche:

```bash
cd packages/observer-web
bun install
bun run dev
```

Öffnen Sie `http://localhost:5173` und melden Sie sich mit `admin@example.com` / `password` an.

## Was der seed-Befehl erstellt

| Was                        | Details                                                |
| -------------------------- | ------------------------------------------------------ |
| **Referenzdaten**          | Länder, Bundesländer, Orte, Büros, Kategorien          |
| **Benutzer**               | Admin- + Mitarbeiterkonten, alle mit Passwort `password`|
| **Projekte**               | 2 Projekte (konfigurierbar mit `--projects`)           |
| **Personen**               | 50 pro Projekt (konfigurierbar mit `--people`)         |
| **Unterstützungseinträge** | Beratungsdatensätze, verknüpft mit Personen            |
| **Migrationseinträge**     | Bewegungsverlauf mit Herkunfts-/Zielorten              |
| **Haushalte**              | Familiengruppen mit Mitgliedern                        |
| **Notizen**                | Fallnotizen zu Personen                                |
| **Haustiere**              | Haustiereinträge mit Tags                              |
| **Tags**                   | Projektbezogene Labels zur Kategorisierung             |

### Standard-Anmeldedaten

| Email               | Password   | Rolle |
| ------------------- | ---------- | ----- |
| `admin@example.com` | `password` | Admin |

### Benutzerdefinierte Seed-Optionen

```bash
# Mehr Projekte und Personen
./observer seed --projects 5 --people 200

# Reproduzierbare Daten (gleicher Seed = gleiche Ausgabe)
./observer seed --seed 42
```

{{% callout type="warning" %}}
Der seed-Befehl **löscht ALLE Tabellen**, bevor Daten eingefügt werden. Führen Sie ihn niemals gegen eine Produktionsdatenbank aus.
{{% /callout %}}

## Schritt-für-Schritt-Anleitung

### 1. Build

```bash
just build
```

Kompiliert die `observer`-Binary.

### 2. Setup

```bash
./observer setup
```

Erstellt `.env` mit sinnvollen Standardwerten, die Verzeichnisse `keys/` und `data/uploads/` sowie ein 4096-Bit RSA-Schlüsselpaar für die JWT-Signierung. Falls `.env` bereits existiert, wird vor dem Überschreiben nachgefragt.

### 3. Dienste starten

```bash
docker compose up -d
```

Startet PostgreSQL und Redis im Hintergrund.

### 4. Migrieren

```bash
./observer migrate up
```

Wendet alle Datenbankmigrationen an.

### 5. Seed

```bash
./observer seed
```

Füllt die Datenbank mit Demodaten. Dieser Schritt macht den Unterschied zwischen einer leeren Instanz und einer funktionierenden Demo, die Sie Stakeholdern zeigen können.

### 6. Starten

```bash
./observer serve
```

Startet den API-Server unter `http://localhost:9000`.

### 7. Frontend (optional)

```bash
cd packages/observer-web
bun install
bun run dev
```

Startet die Weboberfläche unter `http://localhost:5173`.

## Zurücksetzen

Um alles zu löschen und neu zu beginnen:

```bash
docker compose down -v
```

Dann ab Schritt 3 erneut ausführen.
