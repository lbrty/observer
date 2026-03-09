import { type SyntheticEvent, useState } from "react";

import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { Button } from "@/components/button";
import { FormField } from "@/components/form-field";
import { api, HTTPError } from "@/lib/api";
import { useToast } from "@/stores/toast";
import type { ChangePasswordInput } from "@/types/auth";

export function ChangePasswordForm() {
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
