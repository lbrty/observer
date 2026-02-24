import { Field } from "@base-ui/react/field";
import { createFileRoute } from "@tanstack/react-router";
import { type FormEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { PageHeader } from "@/components/page-header";
import { useOffices } from "@/hooks/use-offices";
import { useUpdateUser, useUser } from "@/hooks/use-users";

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

  const offices = officesData?.offices ?? [];

  return (
    <div>
      <PageHeader title={t("admin.users.editTitle")} />
      <form onSubmit={handleSubmit} className="max-w-lg space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <Field.Root>
            <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
              {t("admin.users.firstName")}
            </Field.Label>
            <Field.Control
              value={form.first_name}
              onChange={(e) =>
                setForm((f) => ({ ...f, first_name: e.target.value }))
              }
              className="block w-full rounded-md border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent"
            />
          </Field.Root>
          <Field.Root>
            <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
              {t("admin.users.lastName")}
            </Field.Label>
            <Field.Control
              value={form.last_name}
              onChange={(e) =>
                setForm((f) => ({ ...f, last_name: e.target.value }))
              }
              className="block w-full rounded-md border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent"
            />
          </Field.Root>
        </div>

        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.users.email")}
          </Field.Label>
          <Field.Control
            type="email"
            value={form.email}
            onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))}
            className="block w-full rounded-md border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent"
          />
        </Field.Root>

        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.users.phone")}
          </Field.Label>
          <Field.Control
            type="tel"
            value={form.phone}
            onChange={(e) => setForm((f) => ({ ...f, phone: e.target.value }))}
            className="block w-full rounded-md border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent"
          />
        </Field.Root>

        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.users.role")}
          </Field.Label>
          <select
            value={form.role}
            onChange={(e) => setForm((f) => ({ ...f, role: e.target.value }))}
            className="block w-full rounded-md border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent"
          >
            <option value="admin">admin</option>
            <option value="staff">staff</option>
            <option value="consultant">consultant</option>
            <option value="guest">guest</option>
          </select>
        </Field.Root>

        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.users.office")}
          </Field.Label>
          <select
            value={form.office_id}
            onChange={(e) =>
              setForm((f) => ({ ...f, office_id: e.target.value }))
            }
            className="block w-full rounded-md border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent"
          >
            <option value="">—</option>
            {offices.map((o) => (
              <option key={o.id} value={o.id}>
                {o.name}
              </option>
            ))}
          </select>
        </Field.Root>

        <div className="flex gap-6">
          <label className="flex items-center gap-2 text-sm text-fg-secondary">
            <input
              type="checkbox"
              checked={form.is_active}
              onChange={(e) =>
                setForm((f) => ({ ...f, is_active: e.target.checked }))
              }
              className="accent-accent"
            />
            {t("admin.users.active")}
          </label>
          <label className="flex items-center gap-2 text-sm text-fg-secondary">
            <input
              type="checkbox"
              checked={form.is_verified}
              onChange={(e) =>
                setForm((f) => ({ ...f, is_verified: e.target.checked }))
              }
              className="accent-accent"
            />
            {t("admin.users.verified")}
          </label>
        </div>

        <button
          type="submit"
          disabled={updateUser.isPending}
          className="cursor-pointer rounded-md bg-accent px-4 py-2 text-sm font-medium text-accent-fg hover:opacity-90 disabled:opacity-50"
        >
          {updateUser.isPending
            ? t("admin.users.saving")
            : t("admin.users.save")}
        </button>
      </form>
    </div>
  );
}
