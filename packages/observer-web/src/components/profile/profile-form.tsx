import { type SyntheticEvent, useState } from "react";

import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/ui/alert-banner";
import { Button } from "@/components/ui/button";
import { FormField } from "@/components/forms/form-field";
import { api } from "@/lib/api";
import { handleApiError } from "@/lib/form-error";
import { useToast } from "@/stores/toast";
import type { UpdateProfileInput, User } from "@/types/auth";

export function ProfileForm({ user, setUser }: { user: User | null; setUser: (u: User) => void }) {
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
