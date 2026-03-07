---
title: Bereitstellung
weight: 4
---

Diese Anleitung führt Sie durch die Installation von Observer auf einem Server, den Ihr Team nutzen kann. Sie brauchen kein tiefes technisches Wissen — wenn Sie sich per SSH auf einem Server einloggen und ein paar Befehle ausführen können, schaffen Sie das.

## Docker (empfohlen)

Dies ist der einfachste Weg. Observer wird als einzelnes Docker-Image mit bereits integrierter Weboberfläche ausgeliefert — es muss nichts zusätzlich installiert oder auf der Frontend-Seite konfiguriert werden.

### Was Sie brauchen

- Einen Server mit installiertem Docker und Docker Compose
- Einen Domainnamen, der auf Ihren Server verweist (für HTTPS)

### Schritt 1: Signaturschlüssel generieren

Observer verwendet RSA-Schlüssel zum Signieren von Login-Tokens. Führen Sie diese Befehle auf Ihrem Server aus, um sie zu erstellen:

```bash
mkdir -p keys
openssl genrsa -out keys/jwt_rsa 4096
openssl rsa -in keys/jwt_rsa -pubout -out keys/jwt_rsa.pub
```

Bewahren Sie diese Schlüssel sicher auf. Wenn Sie sie verlieren, müssen sich alle erneut anmelden.

### Schritt 2: Umgebung konfigurieren

Kopieren Sie die Beispiel-Umgebungsdatei und passen Sie sie an Ihre Einrichtung an:

```bash
cp .env.example .env
```

Die wichtigsten Variablen:

| Variable               | Wofür sie da ist                                  | Standard                   |
| ---------------------- | ------------------------------------------------- | -------------------------- |
| `DATABASE_DSN`         | Wie Observer sich mit PostgreSQL verbindet         | _(muss gesetzt werden)_    |
| `REDIS_URL`            | Wie Observer sich mit Redis verbindet              | `redis://localhost:6379/0` |
| `JWT_PRIVATE_KEY_PATH` | Wo Sie den privaten Schlüssel aus Schritt 1 abgelegt haben | `keys/jwt_rsa`     |
| `JWT_PUBLIC_KEY_PATH`  | Wo Sie den öffentlichen Schlüssel aus Schritt 1 abgelegt haben | `keys/jwt_rsa.pub` |
| `CORS_ORIGINS`         | Ihre Domain (z.B. `https://observer.yourorg.org`) | `http://localhost:5173`    |
| `COOKIE_SECURE`        | Auf `true` setzen bei HTTPS (sollten Sie)         | `true`                     |
| `SERVER_HOST`          | Auf welcher Adresse gelauscht wird                | `localhost`                |
| `SERVER_PORT`          | Auf welchem Port gelauscht wird                   | `9000`                     |

Siehe [Umgebungsvariablen](/docs/developers/reference/variables/) für die vollständige Liste.

### Schritt 3: Observer starten

```bash
docker compose up -d
```

Dies startet PostgreSQL, Redis und Observer. Das Datenbankschema wird beim ersten Start automatisch erstellt — kein manueller Migrationsschritt nötig.

### Schritt 4: Prüfen ob es läuft

```bash
curl http://localhost:9000/health
```

Sie sollten sehen:

```json
{ "status": "healthy", "database": "connected", "timestamp": "..." }
```

Wenn Sie das sehen, ist Observer bereit. Öffnen Sie Ihre Domain im Browser, um auf die Weboberfläche zuzugreifen.

## Ohne Docker (VPS / Bare Metal)

Wenn Sie Observer lieber direkt ausführen möchten, bauen Sie die Binärdatei:

```bash
CGO_ENABLED=0 go build -tags production -ldflags="-s -w" -o observer ./cmd/observer
```

Das `-tags production`-Flag bettet die Weboberfläche in die Binärdatei ein. Sie erhalten eine einzelne Datei, die Sie überallhin kopieren können.

Starten Sie sie:

```bash
./observer serve --host 0.0.0.0
```

Sie benötigen PostgreSQL und Redis als separate Dienste. Verweisen Sie `DATABASE_DSN` und `REDIS_URL` darauf.

## HTTPS einrichten

Sie sollten Observer immer hinter einem Reverse Proxy betreiben, der HTTPS übernimmt. So bleiben Anmeldedaten und persönliche Daten während der Übertragung verschlüsselt.

[Caddy](https://caddyserver.com/) ist die einfachste Option — Zertifikate werden automatisch verwaltet:

```
observer.yourorg.org {
    reverse_proxy localhost:9000
}
```

Wenn Sie Nginx oder einen anderen Proxy verwenden, stellen Sie sicher, dass Sie Folgendes setzen:

- `COOKIE_SECURE=true` in Ihrer Umgebung
- `CORS_ORIGINS` auf Ihre tatsächliche Domain (z.B. `https://observer.yourorg.org`)
