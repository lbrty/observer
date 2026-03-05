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
  index.html
  vite.config.ts
  tsconfig.json
  vite-env.d.ts
  src/
    main.tsx                  # app bootstrap (Router + Query + i18n)
    main.css                  # Tailwind entry
    lib/
      api.ts                  # fetch wrapper (credentials: include, 401 auto-refresh)
      i18n.ts                 # i18next setup
    types/
      auth.ts                 # auth DTOs matching backend
    stores/
      auth.tsx                # AuthProvider context + useAuth hook
    locales/
      ky.json                 # Kyrgyz Latin (default)
      en.json                 # English
    routes/
      __root.tsx              # root layout (AuthProvider wraps Outlet)
      _auth.tsx               # public layout — redirects to / if authenticated
      _auth/
        login.tsx             # /login
        register.tsx          # /register
      _app.tsx                # protected layout — redirects to /login if not
      _app/
        index.tsx             # / (dashboard stub)
```

## Running

```bash
just web-install    # install dependencies (bun)
just web-dev        # start dev server (http://localhost:5173)
just web-build      # production build
just web-preview    # preview production build
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
