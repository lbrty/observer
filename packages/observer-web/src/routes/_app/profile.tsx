import { Field } from "@base-ui/react/field";
import { createFileRoute } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { PageHeader } from "@/components/page-header";
import { api, HTTPError } from "@/lib/api";
import { useAuth } from "@/stores/auth";
import type { ChangePasswordInput, UpdateProfileInput, User } from "@/types/auth";

export const Route = createFileRoute("/_app/profile")({
  component: ProfilePage,
});

function ProfilePage() {
  const { t } = useTranslation();
  const { user, setUser } = useAuth();

  return (
    <div className="page-bg-profile mx-auto w-full max-w-xl px-5 py-6">
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
  const [firstName, setFirstName] = useState(user?.first_name ?? "");
  const [lastName, setLastName] = useState(user?.last_name ?? "");
  const [phone, setPhone] = useState(user?.phone ?? "");
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setSaving(true);
    setMessage("");
    setError("");

    try {
      const data: UpdateProfileInput = {
        first_name: firstName,
        last_name: lastName,
        phone,
      };
      const updated = await api.patch("auth/me", { json: data }).json<User>();
      setUser(updated);
      setMessage(t("profile.saved"));
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

  const inputClass =
    "block h-9 w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent";

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h2 className="text-sm font-semibold text-fg">{t("profile.personalInfo")}</h2>

      {message && (
        <div className="rounded-lg bg-foam/10 px-3 py-2 text-sm text-foam">{message}</div>
      )}
      {error && <div className="rounded-lg bg-rose/10 px-3 py-2 text-sm text-rose">{error}</div>}

      <div className="grid grid-cols-2 gap-3">
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("profile.firstName")}
          </Field.Label>
          <Field.Control
            value={firstName}
            onChange={(e) => setFirstName(e.target.value)}
            className={inputClass}
          />
        </Field.Root>

        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("profile.lastName")}
          </Field.Label>
          <Field.Control
            value={lastName}
            onChange={(e) => setLastName(e.target.value)}
            className={inputClass}
          />
        </Field.Root>
      </div>

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("common.email")}
        </Field.Label>
        <Field.Control value={user?.email ?? ""} disabled className={`${inputClass} opacity-50`} />
      </Field.Root>

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("profile.phone")}
        </Field.Label>
        <Field.Control
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
          className={inputClass}
        />
      </Field.Root>

      <button
        type="submit"
        disabled={saving}
        className="cursor-pointer rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-50"
      >
        {saving ? t("profile.saving") : t("profile.save")}
      </button>
    </form>
  );
}

function ChangePasswordForm() {
  const { t } = useTranslation();
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setMessage("");
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
      setMessage(t("profile.passwordChanged"));
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

  const inputClass =
    "block h-9 w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent";

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h2 className="text-sm font-semibold text-fg">{t("profile.changePassword")}</h2>

      {message && (
        <div className="rounded-lg bg-foam/10 px-3 py-2 text-sm text-foam">{message}</div>
      )}
      {error && <div className="rounded-lg bg-rose/10 px-3 py-2 text-sm text-rose">{error}</div>}

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("profile.currentPassword")}
        </Field.Label>
        <Field.Control
          type="password"
          required
          value={currentPassword}
          onChange={(e) => setCurrentPassword(e.target.value)}
          autoComplete="current-password"
          className={inputClass}
        />
      </Field.Root>

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("profile.newPassword")}
        </Field.Label>
        <Field.Control
          type="password"
          required
          minLength={8}
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          autoComplete="new-password"
          className={inputClass}
        />
      </Field.Root>

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("auth.confirmPassword")}
        </Field.Label>
        <Field.Control
          type="password"
          required
          minLength={8}
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          autoComplete="new-password"
          className={inputClass}
        />
      </Field.Root>

      <button
        type="submit"
        disabled={saving}
        className="cursor-pointer rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-50"
      >
        {saving ? t("profile.saving") : t("profile.changePassword")}
      </button>
    </form>
  );
}
