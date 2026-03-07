---
title: Ön Yüz
weight: 3
---

## Teknoloji Yığını

| Konu            | Tercih                                   |
| --------------- | ---------------------------------------- |
| Framework       | React 19 + React Compiler                |
| Bundler         | Vite 6                                   |
| Paket yöneticisi | Bun (workspace monorepo)                 |
| Yönlendirme     | TanStack Router (dosya tabanlı)          |
| Veri çekme      | TanStack Query v5                        |
| Stillendirme    | Tailwind CSS v4                          |
| Headless UI     | Base UI (`@base-ui/react`)               |
| İkonlar         | Phosphor Icons (`@phosphor-icons/react`) |
| i18n            | i18next + react-i18next                  |
| Tür kontrolü    | TypeScript 5.7 (strict)                  |

## Proje yapısı

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

## Çalıştırma

```bash
just web-install    # bağımlılıkları yükle (bun)
just web-dev        # geliştirme sunucusunu başlat (http://localhost:5173)
just web-build      # üretim derlemesi
just web-preview    # üretim derlemesini önizle
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

| Değişken       | Varsayılan              | Açıklama             |
| -------------- | ----------------------- | -------------------- |
| `VITE_API_URL` | `http://localhost:9000` | Backend API temel URL'i |

Vite yalnızca `VITE_` önekli değişkenleri dışa açar.
