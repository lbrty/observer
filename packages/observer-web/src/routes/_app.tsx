import { CheckIcon, SignOutIcon, UserCircleIcon } from "@/components/icons";
import { Menu } from "@base-ui/react/menu";
import { createFileRoute, Link, Navigate, Outlet } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { LANG_KEY, LANGUAGES, THEME_KEY } from "@/lib/constants";
import { useAuth } from "@/stores/auth";

export const Route = createFileRoute("/_app")({
  component: AppLayout,
});

function getStoredTheme(): string {
  return localStorage.getItem(THEME_KEY) || "system";
}

function getStoredLang(): string {
  return localStorage.getItem(LANG_KEY) || "ky";
}

function AppLayout() {
  const { t } = useTranslation();
  const { isAuthenticated, isLoading, user, logout } = useAuth();

  if (isLoading) return null;

  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  return (
    <div className="flex min-h-screen flex-col bg-bg">
      <header className="glass sticky top-0 z-50 border-b border-border-secondary">
        <div className="flex h-13 items-center justify-between px-5">
          <Link
            to="/"
            className="flex items-center gap-2.5 text-sm font-semibold text-fg hover:text-fg"
          >
            <span className="brand-icon inline-flex size-7 items-center justify-center rounded-lg text-xs font-bold text-white">
              O
            </span>
            {t("common.appName")}
          </Link>
          <AvatarMenu email={user?.email ?? ""} onLogout={logout} />
        </div>
      </header>
      <div className="flex flex-1">
        <Outlet />
      </div>
    </div>
  );
}

function AvatarMenu({ email, onLogout }: { email: string; onLogout: () => void }) {
  const { t, i18n } = useTranslation();
  const [theme, setTheme] = useState(getStoredTheme);
  const [lang, setLang] = useState(getStoredLang);

  const themeOptions = [
    { value: "system", label: t("common.themeSystem") },
    { value: "light", label: t("common.themeLight") },
    { value: "dark", label: t("common.themeDark") },
    { value: "light-hc", label: t("common.themeLightHc") },
    { value: "dark-hc", label: t("common.themeDarkHc") },
  ];

  function handleThemeChange(value: unknown) {
    const v = value as string;
    setTheme(v);
    if (v === "system") {
      delete document.documentElement.dataset.theme;
      localStorage.removeItem(THEME_KEY);
    } else {
      document.documentElement.dataset.theme = v;
      localStorage.setItem(THEME_KEY, v);
    }
  }

  function handleLangChange(value: unknown) {
    const v = value as string;
    setLang(v);
    i18n.changeLanguage(v);
    document.documentElement.lang = v;
    localStorage.setItem(LANG_KEY, v);
  }

  return (
    <Menu.Root>
      <Menu.Trigger className="inline-flex size-7 cursor-pointer items-center justify-center rounded-full bg-bg-tertiary text-[11px] font-semibold text-fg-secondary transition-shadow hover:ring-2 hover:ring-accent/30">
        {email.charAt(0).toUpperCase()}
      </Menu.Trigger>
      <Menu.Portal>
        <Menu.Positioner sideOffset={6} align="end" className="z-[100]">
          <Menu.Popup className="w-52 origin-(--transform-origin) rounded-xl border border-border-secondary bg-bg-secondary py-1 shadow-elevated transition-[transform,scale,opacity] data-ending-style:scale-95 data-ending-style:opacity-0 data-starting-style:scale-95 data-starting-style:opacity-0">
            <Menu.Group>
              <Menu.GroupLabel className="px-3 pt-2 pb-1 text-[11px] font-semibold uppercase tracking-wide text-fg-tertiary">
                {t("common.theme")}
              </Menu.GroupLabel>
              <Menu.RadioGroup value={theme} onValueChange={handleThemeChange}>
                {themeOptions.map((opt) => (
                  <Menu.RadioItem
                    key={opt.value}
                    value={opt.value}
                    closeOnClick={false}
                    className="flex cursor-pointer items-center gap-2 px-3 py-1.5 text-sm text-fg outline-none select-none data-highlighted:bg-bg-tertiary"
                  >
                    <span className="inline-flex w-4 items-center justify-center text-accent">
                      <Menu.RadioItemIndicator>
                        <CheckIcon size={14} weight="bold" />
                      </Menu.RadioItemIndicator>
                    </span>
                    {opt.label}
                  </Menu.RadioItem>
                ))}
              </Menu.RadioGroup>
            </Menu.Group>

            <Menu.Separator className="my-1 h-px bg-border-secondary" />

            <Menu.Group>
              <Menu.GroupLabel className="px-3 pt-2 pb-1 text-[11px] font-semibold uppercase tracking-wide text-fg-tertiary">
                {t("common.language")}
              </Menu.GroupLabel>
              <Menu.RadioGroup value={lang} onValueChange={handleLangChange}>
                {LANGUAGES.map((opt) => (
                  <Menu.RadioItem
                    key={opt.value}
                    value={opt.value}
                    closeOnClick={false}
                    className="flex cursor-pointer items-center gap-2 px-3 py-1.5 text-sm text-fg outline-none select-none data-highlighted:bg-bg-tertiary"
                  >
                    <span className="inline-flex w-4 items-center justify-center text-accent">
                      <Menu.RadioItemIndicator>
                        <CheckIcon size={14} weight="bold" />
                      </Menu.RadioItemIndicator>
                    </span>
                    {opt.label}
                  </Menu.RadioItem>
                ))}
              </Menu.RadioGroup>
            </Menu.Group>

            <Menu.Separator className="my-1 h-px bg-border-secondary" />

            <Menu.Item
              render={<Link to="/profile" />}
              className="flex cursor-pointer items-center gap-2 px-3 py-1.5 text-sm text-fg outline-none select-none data-highlighted:bg-bg-tertiary"
            >
              <span className="inline-flex w-4 items-center justify-center text-fg-tertiary">
                <UserCircleIcon size={14} />
              </span>
              {t("profile.title")}
            </Menu.Item>

            <Menu.Item
              onClick={onLogout}
              className="flex cursor-pointer items-center gap-2 px-3 py-1.5 text-sm text-fg outline-none select-none data-highlighted:bg-bg-tertiary"
            >
              <span className="inline-flex w-4 items-center justify-center text-fg-tertiary">
                <SignOutIcon size={14} />
              </span>
              {t("common.logout")}
            </Menu.Item>
          </Menu.Popup>
        </Menu.Positioner>
      </Menu.Portal>
    </Menu.Root>
  );
}
