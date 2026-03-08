import { CheckIcon } from "@/components/icons";
import { UISelect } from "@/components/ui-select";
import { createFileRoute } from "@tanstack/react-router";
import QRCode from "qrcode";
import { type SyntheticEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { Button } from "@/components/button";
import { FormField } from "@/components/form-field";
import { PageHeader } from "@/components/page-header";
import { LANG_KEY, LANGUAGES, THEME_KEY } from "@/lib/constants";
import { api, HTTPError } from "@/lib/api";
import { handleApiError } from "@/lib/form-error";
import { useDisableMFA, useEnableMFA, useMFASetup } from "@/hooks/use-mfa";
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
        <MFASettings />
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
        <FormField label={t("profile.firstName")} value={firstName} onChange={setFirstName} />
        <FormField label={t("profile.lastName")} value={lastName} onChange={setLastName} />
      </div>

      <FormField label={t("common.email")} value={user?.email ?? ""} onChange={() => {}} disabled />

      <FormField label={t("profile.phone")} value={phone} onChange={setPhone} />

      <Button type="submit" disabled={saving}>
        {saving ? t("profile.saving") : t("profile.save")}
      </Button>
    </form>
  );
}

function AppearanceSettings() {
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

function MFASettings() {
  const { t } = useTranslation();
  const { user, setUser } = useAuth();
  const toast = useToast();

  const [showSetup, setShowSetup] = useState(false);
  const [qrDataURL, setQrDataURL] = useState("");
  const [totpCode, setTotpCode] = useState("");
  const [disableCode, setDisableCode] = useState("");
  const [showDisable, setShowDisable] = useState(false);
  const [error, setError] = useState("");

  const setup = useMFASetup(showSetup);
  const enableMFA = useEnableMFA();
  const disableMFA = useDisableMFA();

  useEffect(() => {
    if (setup.data?.otpauth_url) {
      QRCode.toDataURL(setup.data.otpauth_url, { width: 200 }).then(setQrDataURL);
    }
  }, [setup.data?.otpauth_url]);

  async function handleEnable(e: SyntheticEvent) {
    e.preventDefault();
    setError("");
    if (!setup.data) return;
    try {
      await enableMFA.mutateAsync({ secret: setup.data.secret, totpCode });
      if (user) setUser({ ...user, mfa_enabled: true });
      toast.success(t("profile.mfaEnabled"));
      setShowSetup(false);
      setTotpCode("");
      setQrDataURL("");
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  async function handleDisable(e: SyntheticEvent) {
    e.preventDefault();
    setError("");
    try {
      await disableMFA.mutateAsync(disableCode);
      if (user) setUser({ ...user, mfa_enabled: false });
      toast.success(t("profile.mfaDisabled"));
      setShowDisable(false);
      setDisableCode("");
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  const isMFAEnabled = user?.mfa_enabled ?? false;

  return (
    <div className="space-y-4">
      <h2 className="text-sm font-semibold text-fg">{t("profile.twoFactor")}</h2>

      <ErrorBanner message={error} />

      {isMFAEnabled ? (
        <div className="rounded-lg border border-border-secondary p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-fg">{t("profile.mfaActive")}</p>
              <p className="text-xs text-fg-secondary mt-0.5">{t("profile.mfaActiveHint")}</p>
            </div>
            <span className="inline-flex items-center gap-1 rounded-full bg-green-100 px-2.5 py-0.5 text-xs font-medium text-green-800 dark:bg-green-900/30 dark:text-green-300">
              {t("common.enabled")}
            </span>
          </div>

          {!showDisable ? (
            <Button variant="danger" className="mt-3" onClick={() => setShowDisable(true)}>
              {t("profile.disableMFA")}
            </Button>
          ) : (
            <form onSubmit={handleDisable} className="mt-3 space-y-3">
              <p className="text-sm text-fg-secondary">{t("profile.disableMFAHint")}</p>
              <FormField
                label={t("auth.totpCode")}
                value={disableCode}
                onChange={setDisableCode}
                maxLength={6}
              />
              <div className="flex gap-2">
                <Button type="submit" variant="danger" disabled={disableMFA.isPending}>
                  {t("profile.confirmDisable")}
                </Button>
                <Button type="button" variant="secondary" onClick={() => setShowDisable(false)}>
                  {t("common.cancel")}
                </Button>
              </div>
            </form>
          )}
        </div>
      ) : (
        <div className="rounded-lg border border-border-secondary p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-fg">{t("profile.mfaInactive")}</p>
              <p className="text-xs text-fg-secondary mt-0.5">{t("profile.mfaInactiveHint")}</p>
            </div>
            <span className="inline-flex items-center gap-1 rounded-full bg-bg-secondary px-2.5 py-0.5 text-xs font-medium text-fg-tertiary">
              {t("common.disabled")}
            </span>
          </div>

          {!showSetup ? (
            <Button className="mt-3" onClick={() => setShowSetup(true)}>
              {t("profile.setupMFA")}
            </Button>
          ) : (
            <div className="mt-4 space-y-4">
              <p className="text-sm text-fg-secondary">{t("profile.scanQR")}</p>

              {setup.isLoading && <p className="text-sm text-fg-tertiary">{t("common.loading")}</p>}

              {qrDataURL && (
                <div className="flex flex-col items-start gap-3">
                  <img
                    src={qrDataURL}
                    alt={t("profile.totpQrCodeAlt")}
                    className="rounded-lg border border-border-secondary"
                    width={180}
                    height={180}
                  />
                  <details className="text-xs text-fg-tertiary">
                    <summary className="cursor-pointer">{t("profile.cantScanQR")}</summary>
                    <code className="mt-1 block break-all font-mono text-xs text-fg-secondary select-all">
                      {setup.data?.secret}
                    </code>
                  </details>
                </div>
              )}

              {setup.data && (
                <form onSubmit={handleEnable} className="space-y-3">
                  <FormField
                    label={t("profile.enterCode")}
                    value={totpCode}
                    onChange={setTotpCode}
                    maxLength={6}
                    required
                  />
                  <div className="flex gap-2">
                    <Button type="submit" disabled={enableMFA.isPending || totpCode.length !== 6}>
                      {t("profile.verifyAndEnable")}
                    </Button>
                    <Button
                      type="button"
                      variant="secondary"
                      onClick={() => {
                        setShowSetup(false);
                        setTotpCode("");
                        setQrDataURL("");
                      }}
                    >
                      {t("common.cancel")}
                    </Button>
                  </div>
                </form>
              )}
            </div>
          )}
        </div>
      )}
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
