---
title: Ön Yüz
weight: 3
---

## Teknoloji Yığını

| Konu             | Tercih                                   |
| ---------------- | ---------------------------------------- |
| Framework        | React 19 + React Compiler                |
| Bundler          | Vite 6                                   |
| Paket yöneticisi | Bun (workspace monorepo)                 |
| Yönlendirme      | TanStack Router (dosya tabanlı)          |
| Veri çekme       | TanStack Query v5                        |
| Stillendirme     | Tailwind CSS v4                          |
| Headless UI      | Base UI (`@base-ui/react`)               |
| İkonlar          | Phosphor Icons (`@phosphor-icons/react`) |
| i18n             | i18next + react-i18next                  |
| Tür kontrolü     | TypeScript 5.7 (strict)                  |

## Proje yapısı

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

## Çalıştırma

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

## İçe aktarma kuralları

`@/` kısaltması `src/` dizinine çözümlenir. Hem `tsconfig.json` hem de `vite.config.ts` dosyasında yapılandırılmıştır.

İçe aktarma sırası (gruplar arası boş satır):

1. `react`, `react-dom`
2. Harici kütüphaneler (`@tanstack/*`, `@base-ui/*`, `@phosphor-icons/*`, `i18next`)
3. Uygulama kısaltmaları (`@/lib/*`, `@/stores/*`, `@/types/*`)
4. Aynı konumdaki dosyalar (`./constants`, `./types`)
5. Stiller (`.module.css`) — her zaman en sonda

## Yeni rota ekleme

1. `src/routes/` altında bir dosya oluşturun. TanStack Router'ın Vite eklentisi rota ağacını otomatik olarak oluşturur.
2. Korumalı rotalar `_app/` altına gider (kimlik doğrulama gerektirir).
3. Herkese açık kimlik doğrulama rotaları `_auth/` altına gider (zaten giriş yapılmışsa yönlendirir).

Örnek:

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

## Çeviri ekleme

1. Anahtarları hem `src/locales/ky.json` hem de `src/locales/en.json` dosyalarına ekleyin.
2. Bileşenlerde `useTranslation()` ile kullanın:

```tsx
import { useTranslation } from "react-i18next";

function MyComponent() {
  const { t } = useTranslation();
  return <p>{t("namespace.key")}</p>;
}
```

Enterpolasyon, JSON'da `{{variable}}` söz dizimini kullanır:

```json
{ "greeting": "Salam, {{name}}" }
```

```tsx
t("greeting", { name: "Ali" }); // → "Salam, Ali"
```

## Ortam değişkenleri

| Değişken       | Varsayılan              | Açıklama                |
| -------------- | ----------------------- | ----------------------- |
| `VITE_API_URL` | `http://localhost:9000` | Backend API temel URL'i |

Vite yalnızca `VITE_` önekli değişkenleri dışa açar.
