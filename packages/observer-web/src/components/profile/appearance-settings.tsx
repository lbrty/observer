import { useState } from "react";

import { useTranslation } from "react-i18next";

import { UISelect } from "@/components/ui/ui-select";
import { LANG_KEY, LANGUAGES, THEME_KEY } from "@/lib/constants";

interface ThemeSwatch {
  bg: string;
  card: string;
  accent: string;
  green: string;
}

const THEME_SWATCHES: Record<string, ThemeSwatch | null> = {
  system: null,
  light: { bg: "#f8f9fa", card: "#f0f1f5", accent: "#3e63dd", green: "#236e4a" },
  dark: { bg: "#19191c", card: "#222326", accent: "#5b8af8", green: "#3cb179" },
  "light-hc": { bg: "#f5f6fa", card: "#e0e3ea", accent: "#2454c7", green: "#1a6b41" },
  "dark-hc": { bg: "#111113", card: "#19191c", accent: "#6e9bff", green: "#54d48d" },
};

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
        <div className="space-y-2">
          <label className="text-sm text-fg-secondary">{t("common.theme")}</label>
          <div className="flex flex-wrap gap-2">
            {themeOptions.map((opt) => {
              const swatch = THEME_SWATCHES[opt.value];
              const isSelected = theme === opt.value;
              return (
                <button
                  key={opt.value}
                  type="button"
                  onClick={() => handleThemeChange(opt.value)}
                  className={`group flex flex-col items-center gap-1.5 rounded-xl border p-1.5 transition-all ${
                    isSelected
                      ? "border-accent bg-accent/8 ring-2 ring-accent/25"
                      : "border-border-secondary bg-bg-secondary hover:border-border"
                  }`}
                >
                  {swatch ? (
                    <div
                      className="flex h-20 w-28 flex-col overflow-hidden rounded-lg"
                      style={{ background: swatch.bg }}
                    >
                      <div
                        className="m-2 flex-1 rounded-md px-2 py-1.5"
                        style={{ background: swatch.card }}
                      >
                        <div
                          className="mb-1.5 h-1.5 w-3/4 rounded-full"
                          style={{ background: swatch.accent, opacity: 0.5 }}
                        />
                        <div
                          className="h-1.5 w-1/2 rounded-full"
                          style={{ background: swatch.accent, opacity: 0.25 }}
                        />
                      </div>
                      <div className="flex items-center gap-1.5 px-2.5 pb-2">
                        <div className="h-4 w-10 rounded" style={{ background: swatch.accent }} />
                        <div
                          className="h-2 w-2 rounded-full"
                          style={{ background: swatch.green }}
                        />
                      </div>
                    </div>
                  ) : (
                    /* System: split light/dark preview */
                    <div className="flex h-20 w-28 overflow-hidden rounded-lg border border-border-secondary">
                      <div className="flex flex-1 flex-col bg-[#f8f9fa]">
                        <div className="m-2 flex-1 rounded bg-[#f0f1f5]">
                          <div className="mx-1 mt-1 h-1 w-3/4 rounded-full bg-[#3e63dd]/50" />
                        </div>
                        <div className="mx-2 mb-2 h-3 w-6 rounded bg-[#3e63dd]" />
                      </div>
                      <div className="flex flex-1 flex-col bg-[#19191c]">
                        <div className="m-2 flex-1 rounded bg-[#222326]">
                          <div className="mx-1 mt-1 h-1 w-3/4 rounded-full bg-[#5b8af8]/50" />
                        </div>
                        <div className="mx-2 mb-2 h-3 w-6 rounded bg-[#5b8af8]" />
                      </div>
                    </div>
                  )}
                  <span
                    className={`text-xs font-medium leading-none ${
                      isSelected ? "text-accent" : "text-fg-secondary"
                    }`}
                  >
                    {opt.label}
                  </span>
                </button>
              );
            })}
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
