import { useTranslation } from "react-i18next";

import { UISelect } from "@/components/ui-select";

const THEME_KEY = "observer-theme";
const LANG_KEY = "observer-lang";

const LANGUAGES = [
  { value: "ky", label: "Kyrgyzça" },
  { value: "en", label: "English" },
  { value: "ru", label: "Русский" },
  { value: "uk", label: "Українська" },
  { value: "de", label: "Deutsch" },
  { value: "tr", label: "Türkçe" },
];

function getThemeValue(): string {
  return localStorage.getItem(THEME_KEY) || "system";
}

function getLangValue(): string {
  return localStorage.getItem(LANG_KEY) || "ky";
}

export function AppFooter() {
  const { t, i18n } = useTranslation();

  const themeOptions = [
    { value: "system", label: t("common.themeSystem") },
    { value: "light", label: t("common.themeLight") },
    { value: "dark", label: t("common.themeDark") },
    { value: "light-hc", label: t("common.themeLightHc") },
    { value: "dark-hc", label: t("common.themeDarkHc") },
  ];

  function handleThemeChange(value: string) {
    if (value === "system") {
      delete document.documentElement.dataset.theme;
      localStorage.removeItem(THEME_KEY);
    } else {
      document.documentElement.dataset.theme = value;
      localStorage.setItem(THEME_KEY, value);
    }
  }

  function handleLangChange(value: string) {
    i18n.changeLanguage(value);
    document.documentElement.lang = value;
    localStorage.setItem(LANG_KEY, value);
  }

  return (
    <footer className="border-t border-border-secondary">
      <div className="mx-auto flex max-w-5xl items-center justify-end gap-4 px-4 py-3">
        <div className="flex items-center gap-1.5">
          <span className="text-xs text-fg-tertiary">{t("common.theme")}</span>
          <UISelect
            value={getThemeValue()}
            onValueChange={handleThemeChange}
            options={themeOptions}
          />
        </div>
        <div className="flex items-center gap-1.5">
          <span className="text-xs text-fg-tertiary">{t("common.language")}</span>
          <UISelect value={getLangValue()} onValueChange={handleLangChange} options={LANGUAGES} />
        </div>
      </div>
    </footer>
  );
}
