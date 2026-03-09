import { useState } from "react";

import { useTranslation } from "react-i18next";

import { CheckIcon } from "@/components/icons";
import { UISelect } from "@/components/ui-select";
import { LANG_KEY, LANGUAGES, THEME_KEY } from "@/lib/constants";

export function AppearanceSettings() {
  const { t, i18n } = useTranslation();
  const [theme, setTheme] = useState(() => localStorage.getItem(THEME_KEY) || "system");
  const [lang, setLang] = useState(() => localStorage.getItem(LANG_KEY) || "ky");

  const themeOptions = [
    { value: "system", label: t("common.themeSystem") },
    { value: "light", label: t("common.themeLight") },
    { value: "dark", label: t("common.themeDark") },
    { value: "light-hc", label: t("common.themeLightHc") },
    { value: "dark-hc", label: t("common.themeDarkHc") },
  ];

  function handleThemeChange(value: string) {
    setTheme(value);
    if (value === "system") {
      delete document.documentElement.dataset.theme;
      localStorage.removeItem(THEME_KEY);
    } else {
      document.documentElement.dataset.theme = value;
      localStorage.setItem(THEME_KEY, value);
    }
  }

  function handleLangChange(value: string) {
    setLang(value);
    i18n.changeLanguage(value);
    document.documentElement.lang = value;
    localStorage.setItem(LANG_KEY, value);
  }

  return (
    <div className="space-y-4">
      <h2 className="text-sm font-semibold text-fg">{t("profile.appearance")}</h2>

      <div className="space-y-3">
        <div className="space-y-1.5">
          <label className="text-sm text-fg-secondary">{t("common.theme")}</label>
          <div className="flex flex-wrap gap-2">
            {themeOptions.map((opt) => (
              <button
                key={opt.value}
                type="button"
                onClick={() => handleThemeChange(opt.value)}
                className={`inline-flex items-center gap-1.5 rounded-lg border px-3 py-1.5 text-sm transition-colors ${
                  theme === opt.value
                    ? "border-accent bg-accent/10 text-accent"
                    : "border-border-secondary bg-bg-secondary text-fg hover:border-border-primary"
                }`}
              >
                {theme === opt.value && <CheckIcon size={14} weight="bold" />}
                {opt.label}
              </button>
            ))}
          </div>
        </div>

        <div className="space-y-1.5">
          <label className="text-sm text-fg-secondary">{t("common.language")}</label>
          <UISelect value={lang} onValueChange={handleLangChange} options={LANGUAGES} />
        </div>
      </div>
    </div>
  );
}
