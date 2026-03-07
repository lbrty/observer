---
title: Фронтенд
weight: 3
---

## Стек

| Область         | Выбор                                    |
| --------------- | ---------------------------------------- |
| Фреймворк       | React 19 + React Compiler                |
| Сборщик         | Vite 6                                   |
| Пакетный менеджер | Bun (workspace monorepo)               |
| Маршрутизация   | TanStack Router (file-based)             |
| Загрузка данных | TanStack Query v5                        |
| Стилизация      | Tailwind CSS v4                          |
| Headless UI     | Base UI (`@base-ui/react`)               |
| Иконки          | Phosphor Icons (`@phosphor-icons/react`) |
| i18n            | i18next + react-i18next                  |
| Проверка типов  | TypeScript 5.7 (strict)                  |

## Структура проекта

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

## Запуск

```bash
just web-install    # install dependencies (bun)
just web-dev        # start dev server (http://localhost:5173)
just web-build      # production build
just web-preview    # preview production build
```

## Соглашения по импортам

Алиас `@/` указывает на `src/`. Настроен в `tsconfig.json` и `vite.config.ts`.

Порядок импортов (пустая строка между группами):

1. `react`, `react-dom`
2. Внешние библиотеки (`@tanstack/*`, `@base-ui/*`, `@phosphor-icons/*`, `i18next`)
3. Алиасы приложения (`@/lib/*`, `@/stores/*`, `@/types/*`)
4. Соседние модули (`./constants`, `./types`)
5. Стили (`.module.css`) — всегда последними

## Добавление нового маршрута

1. Создайте файл в `src/routes/`. Vite-плагин TanStack Router автоматически генерирует дерево маршрутов.
2. Защищённые маршруты размещаются в `_app/` (требуют аутентификации).
3. Публичные маршруты аутентификации размещаются в `_auth/` (перенаправляют, если уже аутентифицирован).

Пример:

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

## Добавление переводов

1. Добавьте ключи в `src/locales/ky.json` и `src/locales/en.json`.
2. Используйте в компонентах через `useTranslation()`:

```tsx
import { useTranslation } from "react-i18next";

function MyComponent() {
  const { t } = useTranslation();
  return <p>{t("namespace.key")}</p>;
}
```

Интерполяция использует синтаксис `{{variable}}` в JSON:

```json
{ "greeting": "Salam, {{name}}" }
```

```tsx
t("greeting", { name: "Ali" }); // → "Salam, Ali"
```

## Переменные окружения

| Переменная     | По умолчанию            | Описание                |
| -------------- | ----------------------- | ----------------------- |
| `VITE_API_URL` | `http://localhost:9000` | Базовый URL бэкенд API  |

Vite предоставляет доступ только к переменным с префиксом `VITE_`.
