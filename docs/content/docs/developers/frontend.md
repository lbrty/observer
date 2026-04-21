---
title: Frontend
weight: 3
---

## Stack

| Concern         | Choice                                   |
| --------------- | ---------------------------------------- |
| Framework       | React 19 + React Compiler                |
| Bundler         | Vite 6                                   |
| Package manager | Bun (workspace monorepo)                 |
| Routing         | TanStack Router (file-based)             |
| Data fetching   | TanStack Query v5                        |
| Styling         | Tailwind CSS v4                          |
| Headless UI     | Base UI (`@base-ui/react`)               |
| Icons           | Phosphor Icons (`@phosphor-icons/react`) |
| i18n            | i18next + react-i18next                  |
| Type checking   | TypeScript 5.7 (strict)                  |

## Project layout

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

## Running

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

## Import conventions

`@/` alias resolves to `src/`. Configured in both `tsconfig.json` and `vite.config.ts`.

Import order (blank line between groups):

1. `react`, `react-dom`
2. External libs (`@tanstack/*`, `@base-ui/*`, `@phosphor-icons/*`, `i18next`)
3. App aliases (`@/lib/*`, `@/stores/*`, `@/types/*`)
4. Colocated siblings (`./constants`, `./types`)
5. Styles (`.module.css`) — always last

## Adding a new route

1. Create a file under `src/routes/`. TanStack Router's Vite plugin auto-generates the route tree.
2. Protected routes go under `_app/` (requires authentication).
3. Public auth routes go under `_auth/` (redirects away if already authenticated).

Example:

```tsx
// src/routes/_app/settings.tsx → /settings (protected)
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_app/settings")({
  component: SettingsPage,
});

function SettingsPage() {
  return <div>Settings</div>;
}
```

## Adding translations

1. Add keys to both `src/locales/ky.json` and `src/locales/en.json`.
2. Use in components via `useTranslation()`:

```tsx
import { useTranslation } from "react-i18next";

function MyComponent() {
  const { t } = useTranslation();
  return <p>{t("namespace.key")}</p>;
}
```

Interpolation uses `{{variable}}` syntax in JSON:

```json
{ "greeting": "Salam, {{name}}" }
```

```tsx
t("greeting", { name: "Ali" }); // → "Salam, Ali"
```

## Environment variables

| Variable       | Default                 | Description          |
| -------------- | ----------------------- | -------------------- |
| `VITE_API_URL` | `http://localhost:9000` | Backend API base URL |

Vite only exposes variables prefixed with `VITE_`.
