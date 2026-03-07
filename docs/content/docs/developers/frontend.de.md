---
title: Frontend
weight: 3
---

## Stack

| Bereich         | Auswahl                                  |
| --------------- | ---------------------------------------- |
| Framework       | React 19 + React Compiler                |
| Bundler         | Vite 6                                   |
| Paketmanager    | Bun (Workspace-Monorepo)                 |
| Routing         | TanStack Router (dateibasiert)           |
| Datenabruf      | TanStack Query v5                        |
| Styling         | Tailwind CSS v4                          |
| Headless UI     | Base UI (`@base-ui/react`)               |
| Icons           | Phosphor Icons (`@phosphor-icons/react`) |
| i18n            | i18next + react-i18next                  |
| Typprüfung      | TypeScript 5.7 (strict)                  |

## Projektstruktur

```
packages/observer-web/
  index.html
  vite.config.ts
  tsconfig.json
  vite-env.d.ts
  src/
    main.tsx                  # App-Bootstrap (Router + Query + i18n)
    main.css                  # Tailwind-Einstiegspunkt
    lib/
      api.ts                  # Fetch-Wrapper (credentials: include, 401 auto-refresh)
      i18n.ts                 # i18next-Konfiguration
    types/
      auth.ts                 # Auth-DTOs passend zum Backend
    stores/
      auth.tsx                # AuthProvider-Context + useAuth-Hook
    locales/
      ky.json                 # Kirgisisch Latein (Standard)
      en.json                 # Englisch
    routes/
      __root.tsx              # Root-Layout (AuthProvider umschließt Outlet)
      _auth.tsx               # Öffentliches Layout — leitet zu / um wenn authentifiziert
      _auth/
        login.tsx             # /login
        register.tsx          # /register
      _app.tsx                # Geschütztes Layout — leitet zu /login um wenn nicht authentifiziert
      _app/
        index.tsx             # / (Dashboard-Platzhalter)
```

## Ausführen

```bash
just web-install    # Abhängigkeiten installieren (bun)
just web-dev        # Entwicklungsserver starten (http://localhost:5173)
just web-build      # Produktions-Build
just web-preview    # Produktions-Build-Vorschau
```

## Import-Konventionen

Der `@/`-Alias verweist auf `src/`. Konfiguriert in `tsconfig.json` und `vite.config.ts`.

Import-Reihenfolge (Leerzeile zwischen Gruppen):

1. `react`, `react-dom`
2. Externe Bibliotheken (`@tanstack/*`, `@base-ui/*`, `@phosphor-icons/*`, `i18next`)
3. App-Aliase (`@/lib/*`, `@/stores/*`, `@/types/*`)
4. Gleichgeordnete Module (`./constants`, `./types`)
5. Styles (`.module.css`) — immer zuletzt

## Neue Route hinzufügen

1. Erstellen Sie eine Datei unter `src/routes/`. Das Vite-Plugin von TanStack Router generiert den Route-Baum automatisch.
2. Geschützte Routen gehören unter `_app/` (erfordert Authentifizierung).
3. Öffentliche Auth-Routen gehören unter `_auth/` (leitet um, wenn bereits authentifiziert).

Beispiel:

```tsx
// src/routes/_app/settings.tsx → /settings (geschützt)
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_app/settings")({
  component: SettingsPage,
});

function SettingsPage() {
  return <div>Settings</div>;
}
```

## Übersetzungen hinzufügen

1. Fügen Sie Schlüssel sowohl in `src/locales/ky.json` als auch in `src/locales/en.json` hinzu.
2. Verwenden Sie sie in Komponenten über `useTranslation()`:

```tsx
import { useTranslation } from "react-i18next";

function MyComponent() {
  const { t } = useTranslation();
  return <p>{t("namespace.key")}</p>;
}
```

Interpolation verwendet `{{variable}}`-Syntax in JSON:

```json
{ "greeting": "Salam, {{name}}" }
```

```tsx
t("greeting", { name: "Ali" }); // → "Salam, Ali"
```

## Umgebungsvariablen

| Variable       | Standard                | Beschreibung             |
| -------------- | ----------------------- | ------------------------ |
| `VITE_API_URL` | `http://localhost:9000` | Backend-API-Basis-URL    |

Vite stellt nur Variablen mit dem Präfix `VITE_` bereit.
