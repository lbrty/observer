---
title: Frontend
weight: 3
---

## Stek

| Maanisi          | Tañdoo                                   |
| ---------------- | ---------------------------------------- |
| Framework        | React 19 + React Compiler                |
| Bundler          | Vite 6                                   |
| Paket başqaruuçu | Bun (workspace monorepo)                 |
| Routing          | TanStack Router (fajlğa tajanğan)        |
| Maalymat aluu    | TanStack Query v5                        |
| Stilder          | Tailwind CSS v4                          |
| Headless UI      | Base UI (`@base-ui/react`)               |
| Ikonalar         | Phosphor Icons (`@phosphor-icons/react`) |
| i18n             | i18next + react-i18next                  |
| Tip tekşerüü     | TypeScript 5.7 (strict)                  |

## Proekt cajğaşuusu

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

## Işletüü

```bash
just web-install    # install dependencies (bun)
just web-dev        # start dev server (http://localhost:5173)
just web-build      # production build
just web-preview    # preview production build
```

## Import konvensijalary

`@/` alias `src/` papkasuna cañyrat. `tsconfig.json` cana `vite.config.ts` eköösündö da tuuralanğan.

Import tartibi (toptor arasynda boş sap):

1. `react`, `react-dom`
2. Tyşqy librarijalar (`@tanstack/*`, `@base-ui/*`, `@phosphor-icons/*`, `i18next`)
3. Qoldonmo aliastary (`@/lib/*`, `@/stores/*`, `@/types/*`)
4. Cergeliktüü fajldar (`./constants`, `./types`)
5. Stilder (`.module.css`) — dajyma aqyrqy

## Cañy route qoşuu

1. `src/routes/` içine fajl tüzüñüz. TanStack Router'din Vite plugini route darağyn avtomattyq tüzöt.
2. Qorğolğon routelar `_app/` içine barat (autentifikasija talap qylat).
3. Açyq auth routelar `_auth/` içine barat (login bolğon bolso başqa cerge bağyttajt).

Misal:

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

## Qotormolor qoşuu

1. `src/locales/ky.json` cana `src/locales/en.json` eköösünö teñ açqyçtar qoşuñuz.
2. Komponentterde `useTranslation()` arqyluu qoldonuñuz:

```tsx
import { useTranslation } from "react-i18next";

function MyComponent() {
  const { t } = useTranslation();
  return <p>{t("namespace.key")}</p>;
}
```

Interpolasija JSON'do `{{variable}}` sintaksisin qoldonot:

```json
{ "greeting": "Salam, {{name}}" }
```

```tsx
t("greeting", { name: "Ali" }); // → "Salam, Ali"
```

## Çöjrö özgörmölörü

| Özgörmö        | Default                 | Taanıştıruu             |
| -------------- | ----------------------- | ----------------------- |
| `VITE_API_URL` | `http://localhost:9000` | Backend API bazalyq URL |

Vite `VITE_` prefiksi bar özgörmölördü ğana açat.
