import { CheckIcon } from "@/components/icons";
import { UISelect } from "@/components/ui-select";
import { createFileRoute } from "@tanstack/react-router";
import { type SyntheticEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { Button } from "@/components/button";
import { FormField } from "@/components/form-field";
import { PageHeader } from "@/components/page-header";
import { LANG_KEY, LANGUAGES, THEME_KEY } from "@/lib/constants";
import { api, HTTPError } from "@/lib/api";
import { handleApiError } from "@/lib/form-error";
import { useAuth } from "@/stores/auth";
import { useToast } from "@/stores/toast";
import type { ChangePasswordInput, UpdateProfileInput, User } from "@/types/auth";

export const Route = createFileRoute("/_app/profile")({
  component: ProfilePage,
});

function ProfilePage() {
  const { t } = useTranslation();
  const { user, setUser } = useAuth();

  return (
    <div className="mx-auto w-full max-w-xl px-5 py-6">
      <PageHeader title={t("profile.title")} />
      <div className="space-y-6">
        <ProfileForm user={user} setUser={setUser} />
        <div className="h-px bg-border-secondary" />
        <AppearanceSettings />
        <div className="h-px bg-border-secondary" />
        <ChangePasswordForm />
      </div>
    </div>
  );
}

function ProfileForm({ user, setUser }: { user: User | null; setUser: (u: User) => void }) {
  const { t } = useTranslation();
  const toast = useToast();
  const [firstName, setFirstName] = useState(user?.first_name ?? "");
  const [lastName, setLastName] = useState(user?.last_name ?? "");
  const [phone, setPhone] = useState(user?.phone ?? "");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    setSaving(true);
    setError("");

    try {
      const data: UpdateProfileInput = {
        first_name: firstName,
        last_name: lastName,
        phone,
      };
      const updated = await api.patch("auth/me", { json: data }).json<User>();
      setUser(updated);
      toast.success(t("profile.saved"));
    } catch (err) {
      setError(await handleApiError(err, t));
    } finally {
      setSaving(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h2 className="text-sm font-semibold text-fg">{t("profile.personalInfo")}</h2>

      <ErrorBanner message={error} />

      <div className="grid grid-cols-2 gap-3">
        <FormField
          label={t("profile.firstName")}
          value={firstName}
          onChange={setFirstName}
        />
        <FormField
          label={t("profile.lastName")}
          value={lastName}
          onChange={setLastName}
        />
      </div>

      <FormField
        label={t("common.email")}
        value={user?.email ?? ""}
        onChange={() => {}}
        disabled
      />

      <FormField
        label={t("profile.phone")}
        value={phone}
        onChange={setPhone}
      />

      <Button type="submit" disabled={saving}>
        {saving ? t("profile.saving") : t("profile.save")}
      </Button>
    </form>
  );
}

function AppearanceSettings() {
  const { t, i18n } = useTranslation();
  const [theme, setTheme] = useState(
    () => localStorage.getItem(THEME_KEY) || "system",
  );
  const [lang, setLang] = useState(
    () => localStorage.getItem(LANG_KEY) || "ky",
  );

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
      <h2 className="text-sm font-semibold text-fg">
        {t("profile.appearance")}
      </h2>

      <div className="space-y-3">
        <div className="space-y-1.5">
          <label className="text-sm text-fg-secondary">
            {t("common.theme")}
          </label>
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
          <label className="text-sm text-fg-secondary">
            {t("common.language")}
          </label>
          <UISelect
            value={lang}
            onValueChange={handleLangChange}
            options={LANGUAGES}
          />
        </div>
      </div>
    </div>
  );
}

function ChangePasswordForm() {
  const { t } = useTranslation();
  const toast = useToast();
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    setError("");

    if (newPassword !== confirmPassword) {
      setError(t("auth.passwordsMismatch"));
      return;
    }

    if (newPassword.length < 8) {
      setError(t("auth.passwordTooShort"));
      return;
    }

    setSaving(true);

    try {
      const data: ChangePasswordInput = {
        current_password: currentPassword,
        new_password: newPassword,
      };
      await api.post("auth/change-password", { json: data });
      toast.success(t("profile.passwordChanged"));
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");
    } catch (err) {
      if (err instanceof HTTPError) {
        const body = await err.response.json().catch(() => null);
        if (body?.code === "errors.auth.invalidCredentials") {
          setError(t("profile.wrongPassword"));
        } else {
          const code = body?.code;
          const translated = code ? t(code, { defaultValue: "" }) : "";
          setError(translated || body?.error || err.message);
        }
      } else {
        setError(t("common.unexpectedError"));
      }
    } finally {
      setSaving(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h2 className="text-sm font-semibold text-fg">{t("profile.changePassword")}</h2>

      <ErrorBanner message={error} />

      <FormField
        label={t("profile.currentPassword")}
        type="password"
        required
        value={currentPassword}
        onChange={setCurrentPassword}
      />

      <FormField
        label={t("profile.newPassword")}
        type="password"
        required
        value={newPassword}
        onChange={setNewPassword}
      />

      <FormField
        label={t("auth.confirmPassword")}
        type="password"
        required
        value={confirmPassword}
        onChange={setConfirmPassword}
      />

      <Button type="submit" disabled={saving}>
        {saving ? t("profile.saving") : t("profile.changePassword")}
      </Button>
    </form>
  );
}
