import { createFileRoute } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { Button } from "@/components/button";
import { FormField } from "@/components/form-field";
import { PageHeader } from "@/components/page-header";
import { api, HTTPError } from "@/lib/api";
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

  async function handleSubmit(e: FormEvent) {
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
      if (err instanceof HTTPError) {
        const body = await err.response.json().catch(() => null);
        setError(body?.error ?? err.message);
      } else {
        setError(t("common.unexpectedError"));
      }
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

function ChangePasswordForm() {
  const { t } = useTranslation();
  const toast = useToast();
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(e: FormEvent) {
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
          setError(body?.error ?? err.message);
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
