import { MagnifyingGlassIcon } from "@/components/icons";
import { Field } from "@base-ui/react/field";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { DataTable, type Column } from "@/components/data-table";
import { FormDialog } from "@/components/form-dialog";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { StatusBadge, StatusDot } from "@/components/status-badge";
import { UISelect } from "@/components/ui-select";
import { UISwitch } from "@/components/ui-switch";
import { useOffices } from "@/hooks/use-offices";
import { useCreateUser, useUsers } from "@/hooks/use-users";
import type { AdminUser } from "@/types/admin";

export const Route = createFileRoute("/_app/admin/users/")({
  component: UsersPage,
});

function Initials({ first, last }: { first: string; last: string }) {
  const letters =
    `${first?.charAt(0) ?? ""}${last?.charAt(0) ?? ""}`.toUpperCase() || "?";
  return (
    <span className="inline-flex size-8 shrink-0 items-center justify-center rounded-full bg-accent/10 text-xs font-semibold text-accent">
      {letters}
    </span>
  );
}

function UsersPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [role, setRole] = useState("");
  const [isActive, setIsActive] = useState("");
  const [createOpen, setCreateOpen] = useState(false);

  const params = {
    page,
    per_page: 20,
    ...(search && { search }),
    ...(role && { role }),
    ...(isActive !== "" && { is_active: isActive === "true" }),
  };

  const { data, isLoading } = useUsers(params);

  const roleOptions = [
    { label: t("admin.users.allRoles"), value: "" },
    { label: "admin", value: "admin" },
    { label: "staff", value: "staff" },
    { label: "consultant", value: "consultant" },
    { label: "guest", value: "guest" },
  ];

  const statusOptions = [
    { label: t("admin.users.allStatuses"), value: "" },
    { label: t("admin.users.active"), value: "true" },
    { label: t("admin.users.inactive"), value: "false" },
  ];

  const columns: Column<AdminUser>[] = [
    {
      key: "name",
      header: t("admin.users.name"),
      render: (u) => (
        <div className="flex items-center gap-3">
          <Initials first={u.first_name} last={u.last_name} />
          <span className="font-medium text-fg">
            {u.first_name} {u.last_name}
          </span>
        </div>
      ),
    },
    {
      key: "email",
      header: t("admin.users.email"),
      render: (u) => <span className="text-fg-secondary">{u.email}</span>,
    },
    {
      key: "role",
      header: t("admin.users.role"),
      render: (u) => <StatusBadge label={u.role} />,
    },
    {
      key: "active",
      header: t("admin.users.active"),
      render: (u) => <StatusDot active={u.is_active} />,
    },
    {
      key: "verified",
      header: t("admin.users.verified"),
      render: (u) => <StatusDot active={u.is_verified} />,
    },
    {
      key: "created",
      header: t("admin.users.created"),
      render: (u) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {new Date(u.created_at).toLocaleDateString("en-CA")}
        </span>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t("admin.users.title")}
        action={
          <button
            type="button"
            onClick={() => setCreateOpen(true)}
            className="cursor-pointer rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90"
          >
            {t("admin.users.add")}
          </button>
        }
      />

      <div className="mb-4 flex gap-3">
        <div className="relative">
          <MagnifyingGlassIcon
            size={14}
            className="absolute top-1/2 left-3 -translate-y-1/2 text-fg-tertiary"
          />
          <input
            placeholder={t("admin.users.search")}
            value={search}
            onChange={(e) => {
              setSearch(e.target.value);
              setPage(1);
            }}
            className="rounded-lg border border-border-secondary bg-bg-secondary py-1.5 pr-3 pl-8 text-sm text-fg outline-none focus:border-accent"
          />
        </div>
        <UISelect
          value={role}
          onValueChange={(v) => {
            setRole(v);
            setPage(1);
          }}
          options={roleOptions}
          placeholder={t("admin.users.allRoles")}
        />
        <UISelect
          value={isActive}
          onValueChange={(v) => {
            setIsActive(v);
            setPage(1);
          }}
          options={statusOptions}
          placeholder={t("admin.users.allStatuses")}
        />
      </div>

      <DataTable
        columns={columns}
        data={data?.users ?? []}
        keyExtractor={(u) => u.id}
        onRowClick={(u) =>
          navigate({ to: "/admin/users/$userId", params: { userId: u.id } })
        }
        isLoading={isLoading}
      />

      {data && (
        <Pagination
          page={data.page}
          perPage={data.per_page}
          total={data.total}
          onChange={setPage}
        />
      )}

      <CreateUserDialog open={createOpen} onOpenChange={setCreateOpen} />
    </div>
  );
}

const inputClass =
  "block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent";

function CreateUserDialog({
  open,
  onOpenChange,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const { t } = useTranslation();
  const createUser = useCreateUser();
  const { data: officesData } = useOffices();

  const [form, setForm] = useState({
    first_name: "",
    last_name: "",
    email: "",
    phone: "",
    password: "",
    role: "consultant",
    office_id: "",
    is_active: true,
    is_verified: true,
  });

  const roleOptions = [
    { label: "admin", value: "admin" },
    { label: "staff", value: "staff" },
    { label: "consultant", value: "consultant" },
    { label: "guest", value: "guest" },
  ];

  const officeOptions = [
    { label: "—", value: "" },
    ...(officesData ?? []).map((o) => ({ label: o.name, value: o.id })),
  ];

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    await createUser.mutateAsync({
      first_name: form.first_name,
      last_name: form.last_name || undefined,
      email: form.email,
      phone: form.phone || undefined,
      password: form.password,
      role: form.role,
      office_id: form.office_id || null,
      is_active: form.is_active,
      is_verified: form.is_verified,
    });
    onOpenChange(false);
    setForm({
      first_name: "",
      last_name: "",
      email: "",
      phone: "",
      password: "",
      role: "consultant",
      office_id: "",
      is_active: true,
      is_verified: true,
    });
  }

  return (
    <FormDialog
      open={open}
      onOpenChange={onOpenChange}
      title={t("admin.users.addTitle")}
      loading={createUser.isPending}
      onSubmit={handleSubmit}
      maxWidth="md"
    >
      <div className="grid grid-cols-2 gap-3">
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.users.firstName")}
          </Field.Label>
          <Field.Control
            required
            value={form.first_name}
            onChange={(e) =>
              setForm((f) => ({ ...f, first_name: e.target.value }))
            }
            className={inputClass}
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
            className={inputClass}
          />
        </Field.Root>
      </div>

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("admin.users.email")}
        </Field.Label>
        <Field.Control
          type="email"
          required
          value={form.email}
          onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))}
          className={inputClass}
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
          className={inputClass}
        />
      </Field.Root>

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("admin.users.password")}
        </Field.Label>
        <Field.Control
          type="password"
          required
          minLength={8}
          value={form.password}
          onChange={(e) => setForm((f) => ({ ...f, password: e.target.value }))}
          className={inputClass}
        />
      </Field.Root>

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
    </FormDialog>
  );
}
