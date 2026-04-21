---
title: Frontend
weight: 3
---

## Stack

| Bereich      | Auswahl                                  |
| ------------ | ---------------------------------------- |
| Framework    | React 19 + React Compiler                |
| Bundler      | Vite 6                                   |
| Paketmanager | Bun (Workspace-Monorepo)                 |
| Routing      | TanStack Router (dateibasiert)           |
| Datenabruf   | TanStack Query v5                        |
| Styling      | Tailwind CSS v4                          |
| Headless UI  | Base UI (`@base-ui/react`)               |
| Icons        | Phosphor Icons (`@phosphor-icons/react`) |
| i18n         | i18next + react-i18next                  |
| Typprüfung   | TypeScript 5.7 (strict)                  |

## Projektstruktur

```
packages/observer-web/
  src/
    main.tsx          # app bootstrap (Router + Query + i18n)
    main.css          # Tailwind entry
    routes/           # TanStack Router file-based routes (_app/, _auth/)
    components/       # UI components, grouped by domain
      ui/             # atoms: button, badge, icons, toast…
      layout/         # app shell: page-header, sidebar-link…
      forms/          # form-field, filter-bar, comboboxes
      table/          # data-table, pagination, row-actions
      dialogs/        # confirm-dialog, form-dialog
      drawer/         # drawer-shell
      people/ pets/ support/ households/ migration/
      documents/ tags/ users/ permissions/ profile/
      reports/ charts/ auth/ date-picker/ search-palette/
    hooks/            # React Query hooks, grouped by domain
      reference/      # countries, states, places, offices, categories
      people/ pets/ support/ households/ migration/
      documents/ tags/ notes/ users/ projects/ reports/
    stores/           # Zustand: auth, toast
    types/            # TypeScript types matching API responses
    lib/              # api client, i18n, export, params helpers
    constants/        # enum → i18n key maps
    locales/          # ky.json (default), en.json, uk.json, ru.json…
```

## Ausführen

```bash
# Install dependencies
cd packages/observer-web && bun install

# Backend + frontend dev servers concurrently
just dev

# Frontend only (http://localhost:5173)
cd packages/observer-web && bun run dev

# Production build (frontend + embedded Go binary)
just build-prod

# Format frontend code
cd packages/observer-web && bun run fmt
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

| Variable       | Standard                | Beschreibung          |
| -------------- | ----------------------- | --------------------- |
| `VITE_API_URL` | `http://localhost:9000` | Backend-API-Basis-URL |

Vite stellt nur Variablen mit dem Präfix `VITE_` bereit.
