---
title: Фронтенд
weight: 3
---

## Стек

| Напрямок        | Вибір                                    |
| --------------- | ---------------------------------------- |
| Фреймворк       | React 19 + React Compiler                |
| Збірник         | Vite 6                                   |
| Пакетний менеджер | Bun (workspace monorepo)               |
| Маршрутизація   | TanStack Router (file-based)             |
| Отримання даних | TanStack Query v5                        |
| Стилізація      | Tailwind CSS v4                          |
| Headless UI     | Base UI (`@base-ui/react`)               |
| Іконки          | Phosphor Icons (`@phosphor-icons/react`) |
| i18n            | i18next + react-i18next                  |
| Перевірка типів | TypeScript 5.7 (strict)                  |

## Структура проєкту

```
packages/observer-web/
  index.html
  vite.config.ts
  tsconfig.json
  vite-env.d.ts
  src/
    main.tsx                  # ініціалізація застосунку (Router + Query + i18n)
    main.css                  # точка входу Tailwind
    lib/
      api.ts                  # обгортка fetch (credentials: include, 401 auto-refresh)
      i18n.ts                 # налаштування i18next
    types/
      auth.ts                 # auth DTO, що відповідають бекенду
    stores/
      auth.tsx                # AuthProvider контекст + useAuth хук
    locales/
      ky.json                 # киргизька латиниця (за замовчуванням)
      en.json                 # англійська
    routes/
      __root.tsx              # кореневий макет (AuthProvider обгортає Outlet)
      _auth.tsx               # публічний макет — перенаправляє на / якщо автентифікований
      _auth/
        login.tsx             # /login
        register.tsx          # /register
      _app.tsx                # захищений макет — перенаправляє на /login якщо ні
      _app/
        index.tsx             # / (заглушка панелі керування)
```

## Запуск

```bash
just web-install    # install dependencies (bun)
just web-dev        # start dev server (http://localhost:5173)
just web-build      # production build
just web-preview    # preview production build
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

| Змінна         | За замовчуванням        | Опис                   |
| -------------- | ----------------------- | ---------------------- |
| `VITE_API_URL` | `http://localhost:9000` | Базова URL-адреса API бекенду |

Vite експонує лише змінні з префіксом `VITE_`.
