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

## Işletüü

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
