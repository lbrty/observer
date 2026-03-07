---
title: Erste Schritte
weight: 2
---

## In 5 Minuten lauffähig

Sie brauchen keinen Server, keinen Hosting-Anbieter und keine IT-Abteilung. Wenn Sie einen Laptop mit installiertem Docker haben, können Sie Observer sofort ausprobieren.

```bash
git clone https://github.com/lbrty/observer.git
cd observer
cp .env.example .env
just generate-keys
just docker-up
just run
```

Öffnen Sie `http://localhost:9000/health` in Ihrem Browser. Wenn Sie `"status":"healthy"` sehen, läuft das Backend.

Starten Sie dann die Weboberfläche:

```bash
just web-dev
```

Öffnen Sie `http://localhost:5173` — Sie sehen Observer.

## Was Sie gerade gestartet haben

- Ein **Backend**, das die API bereitstellt — übernimmt Authentifizierung, Datenspeicherung und Berichte
- Eine **Datenbank** mit Tabellen für Personen, Haushalte, Unterstützungseinträge, Migrationshistorie, Dokumente und Haustiere
- Eine **Weboberfläche** mit Projektverwaltung, rollenbasierter Zugriffskontrolle und integrierter Berichterstattung
- **Automatische Anmeldesicherheit** — Tokens werden bei jeder Sitzungsaktualisierung rotiert

All das läuft auf einem einzelnen Rechner. In der Produktion wird es zu einer einzelnen Datei kompiliert, die Sie auf jeden Server kopieren können.

## Bereit für den Produktiveinsatz?

Um von „Ausprobieren" zu „Mein Team nutzt das täglich" zu kommen, brauchen Sie:

| Was | Warum |
| --- | --- |
| Einen Server (VPS oder lokal) | Observer ist selbst gehostet — Ihre Daten verlassen nie Ihre Infrastruktur |
| PostgreSQL | Der einzige externe Dienst, den Observer benötigt |
| Etwa 30 Minuten | Führen Sie `docker compose up` auf einem Server aus, auf den Ihre Domain verweist |

Kein Abonnement. Keine nutzerbezogenen Gebühren. Kein Vendor Lock-in. Die Daten und die Bereitstellung gehören Ihnen.

Siehe [Bereitstellung](/docs/guide/deployment/) für die schrittweise Produktionseinrichtung.

## Für Entwickler: Lokale Einrichtung

Wenn Sie an Observer selbst arbeiten möchten, benötigen Sie folgende Tools:

| Tool | Version | Installation |
| --- | --- | --- |
| Go | 1.25.* | https://go.dev/dl/ |
| Bun | latest | https://bun.sh/ |
| Docker + Compose | latest | https://docs.docker.com/get-docker/ |
| Just | latest | https://github.com/casey/just#installation |

### 1. Klonen und Abhängigkeiten installieren

```bash
git clone https://github.com/lbrty/observer.git
cd observer
go mod download
bun install
```

### 2. Umgebung konfigurieren

```bash
cp .env.example .env
```

Die Standardwerte funktionieren sofort mit der mitgelieferten `docker-compose.yml`. Keine Bearbeitung nötig.

### 3. Signaturschlüssel generieren

```bash
just generate-keys
```

Dies erstellt ein Schlüsselpaar, das Observer zum Signieren von Login-Tokens verwendet. Das `keys/`-Verzeichnis ist in `.gitignore` — jeder Entwickler generiert sein eigenes.

### 4. Alles starten

```bash
just docker-up    # startet PostgreSQL und Redis
just run          # startet das Backend auf :9000 (führt Migrationen automatisch aus)
just web-dev      # startet das Frontend auf :5173
```

## Etwas funktioniert nicht?

**Port 5432 bereits belegt** — Wahrscheinlich läuft ein lokales PostgreSQL. Stoppen Sie es oder ändern Sie den Port in `docker-compose.yml`.

**"no such file or directory" für Schlüsselpfade** — Sie müssen zuerst `just generate-keys` ausführen.

**Migration schlägt mit "connection refused" fehl** — Der Datenbank-Container ist möglicherweise noch nicht bereit. Warten Sie einige Sekunden nach `just docker-up` und versuchen Sie es erneut.
