---
title: Фронтенд
weight: 3
---

## Стек

| Напрямок          | Вибір                                    |
| ----------------- | ---------------------------------------- |
| Фреймворк         | React 19 + React Compiler                |
| Збірник           | Vite 6                                   |
| Пакетний менеджер | Bun (workspace monorepo)                 |
| Маршрутизація     | TanStack Router (file-based)             |
| Отримання даних   | TanStack Query v5                        |
| Стилізація        | Tailwind CSS v4                          |
| Headless UI       | Base UI (`@base-ui/react`)               |
| Іконки            | Phosphor Icons (`@phosphor-icons/react`) |
| i18n              | i18next + react-i18next                  |
| Перевірка типів   | TypeScript 5.7 (strict)                  |

## Структура проєкту

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

## Запуск

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

## Конвенції імпортів

Аліас `@/` вказує на `src/`. Налаштовано як у `tsconfig.json`, так і у `vite.config.ts`.

Порядок імпортів (порожній рядок між групами):

1. `react`, `react-dom`
2. Зовнішні бібліотеки (`@tanstack/*`, `@base-ui/*`, `@phosphor-icons/*`, `i18next`)
3. Аліаси застосунку (`@/lib/*`, `@/stores/*`, `@/types/*`)
4. Сусідні файли (`./constants`, `./types`)
5. Стилі (`.module.css`) — завжди останніми

## Додавання нового маршруту

1. Створіть файл у `src/routes/`. Vite-плагін TanStack Router автоматично генерує дерево маршрутів.
2. Захищені маршрути розміщуються у `_app/` (потребують автентифікації).
3. Публічні маршрути автентифікації розміщуються у `_auth/` (перенаправляють, якщо вже автентифікований).

Приклад:

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

## Додавання перекладів

1. Додайте ключі до обох файлів `src/locales/ky.json` та `src/locales/en.json`.
2. Використовуйте в компонентах через `useTranslation()`:

```tsx
import { useTranslation } from "react-i18next";

function MyComponent() {
  const { t } = useTranslation();
  return <p>{t("namespace.key")}</p>;
}
```

Інтерполяція використовує синтаксис `{{variable}}` у JSON:

```json
{ "greeting": "Salam, {{name}}" }
```

```tsx
t("greeting", { name: "Ali" }); // → "Salam, Ali"
```

## Змінні середовища

| Змінна         | За замовчуванням        | Опис                          |
| -------------- | ----------------------- | ----------------------------- |
| `VITE_API_URL` | `http://localhost:9000` | Базова URL-адреса API бекенду |

Vite експонує лише змінні з префіксом `VITE_`.
