import { Field } from "@base-ui/react/field";
import { createFileRoute } from "@tanstack/react-router";
import { type FormEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { ErrorBanner } from "@/components/alert-banner";
import { FormDialog } from "@/components/form-dialog";
import { FormField } from "@/components/form-field";
import { PageHeader } from "@/components/page-header";
import { UISelect } from "@/components/ui-select";
import { UISwitch } from "@/components/ui-switch";
import { useOffices } from "@/hooks/use-offices";
import { useUpdateUser, useUser } from "@/hooks/use-users";
import { api, HTTPError } from "@/lib/api";

export const Route = createFileRoute("/_app/admin/users/$userId")({
  component: UserDetailPage,
});

function UserDetailPage() {
  const { t } = useTranslation();
  const { userId } = Route.useParams();
  const { data: user, isLoading } = useUser(userId);
  const { data: officesData } = useOffices();
  const updateUser = useUpdateUser();

  const [form, setForm] = useState({
    first_name: "",
    last_name: "",
    email: "",
    phone: "",
    role: "",
    office_id: "",
    is_active: false,
    is_verified: false,
  });

  useEffect(() => {
    if (user) {
      setForm({
        first_name: user.first_name,
        last_name: user.last_name,
        email: user.email,
        phone: user.phone,
        role: user.role,
        office_id: user.office_id ?? "",
        is_active: user.is_active,
        is_verified: user.is_verified,
      });
    }
  }, [user]);

  if (isLoading) return null;
  if (!user) return null;

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    await updateUser.mutateAsync({
      id: userId,
      data: {
        first_name: form.first_name,
        last_name: form.last_name,
        email: form.email,
        phone: form.phone,
        role: form.role,
        office_id: form.office_id || null,
        is_active: form.is_active,
        is_verified: form.is_verified,
      },
    });
  }

  const offices = officesData ?? [];

  const roleOptions = [
    { label: t("admin.users.roleAdmin"), value: "admin" },
    { label: t("admin.users.roleStaff"), value: "staff" },
    { label: t("admin.users.roleConsultant"), value: "consultant" },
    { label: t("admin.users.roleGuest"), value: "guest" },
  ];

  const officeOptions = [
    { label: "—", value: "" },
    ...offices.map((o) => ({ label: o.name, value: o.id })),
  ];

  return (
    <div>
      <PageHeader title={t("admin.users.editTitle")} />
      <form onSubmit={handleSubmit} className="max-w-lg space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <FormField
            label={t("admin.users.firstName")}
            value={form.first_name}
            onChange={(v) => setForm((f) => ({ ...f, first_name: v }))}
          />
          <FormField
            label={t("admin.users.lastName")}
            value={form.last_name}
            onChange={(v) => setForm((f) => ({ ...f, last_name: v }))}
          />
        </div>

        <FormField
          label={t("admin.users.email")}
          type="email"
          value={form.email}
          onChange={(v) => setForm((f) => ({ ...f, email: v }))}
        />

        <FormField
          label={t("admin.users.phone")}
          type="tel"
          value={form.phone}
          onChange={(v) => setForm((f) => ({ ...f, phone: v }))}
        />

        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.users.role")}
          </Field.Label>
          <UISelect
            value={form.role}
            onValueChange={(v) => setForm((f) => ({ ...f, role: v }))}
            options={roleOptions}
            fullWidth
          />
        </Field.Root>

        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.users.office")}
          </Field.Label>
          <UISelect
            value={form.office_id}
            onValueChange={(v) => setForm((f) => ({ ...f, office_id: v }))}
            options={officeOptions}
            placeholder="—"
            fullWidth
          />
        </Field.Root>

        <div className="flex gap-6">
          <UISwitch
            checked={form.is_active}
            onCheckedChange={(v) => setForm((f) => ({ ...f, is_active: v }))}
            label={t("admin.users.active")}
          />
          <UISwitch
            checked={form.is_verified}
            onCheckedChange={(v) => setForm((f) => ({ ...f, is_verified: v }))}
            label={t("admin.users.verified")}
          />
        </div>

        <button
          type="submit"
          disabled={updateUser.isPending}
          className="cursor-pointer rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90 disabled:opacity-50"
        >
          {updateUser.isPending ? t("admin.users.saving") : t("admin.users.save")}
        </button>
      </form>

      <div className="mt-6 h-px bg-border-secondary" />

      <ResetPasswordSection userId={userId} />
    </div>
  );
}

function ResetPasswordSection({ userId }: { userId: string }) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const [password, setPassword] = useState("");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (password.length < 8) {
      setError(t("auth.passwordTooShort"));
      return;
    }

    setSaving(true);
    setError("");
    try {
      await api.post(`admin/users/${userId}/reset-password`, {
        json: { new_password: password },
      });
      setPassword("");
      setOpen(false);
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
    <div className="mt-4">
      <button
        type="button"
        onClick={() => setOpen(true)}
        className="cursor-pointer rounded-lg border border-border-secondary px-4 py-2 text-sm font-medium text-fg-secondary hover:bg-bg-tertiary"
      >
        {t("admin.users.resetPassword")}
      </button>

      <FormDialog
        open={open}
        onOpenChange={setOpen}
        title={t("admin.users.resetPasswordTitle")}
        loading={saving}
        onSubmit={handleSubmit}
      >
        <ErrorBanner message={error} />
        <FormField
          label={t("admin.users.newPassword")}
          type="password"
          required
          value={password}
          onChange={setPassword}
        />
      </FormDialog>
    </div>
  );
}
