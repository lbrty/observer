import { Field } from "@base-ui/react/field";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { type SyntheticEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import type { Column } from "@/components/data-table";
import { DataTablePage } from "@/components/data-table-page";
import type { FilterDef } from "@/components/filter-bar";
import { FormDialog } from "@/components/form-dialog";
import { FormField } from "@/components/form-field";
import { UsersIcon } from "@/components/icons";
import { StatusBadge, StatusDot } from "@/components/status-badge";
import { UISelect } from "@/components/ui-select";
import { UISwitch } from "@/components/ui-switch";
import { UserInitials } from "@/components/user-initials";
import { useOffices } from "@/hooks/use-offices";
import { useCreateUser, useUsers } from "@/hooks/use-users";
import type { AdminUser } from "@/types/admin";

export const Route = createFileRoute("/_app/admin/users/")({
  component: UsersPage,
});

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
    { label: t("admin.users.roleAdmin"), value: "admin" },
    { label: t("admin.users.roleStaff"), value: "staff" },
    { label: t("admin.users.roleConsultant"), value: "consultant" },
    { label: t("admin.users.roleGuest"), value: "guest" },
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
          <UserInitials firstName={u.first_name} lastName={u.last_name} />
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

  const filters: FilterDef[] = [
    {
      type: "search",
      placeholder: t("admin.users.search"),
      value: search,
      onChange: (v) => { setSearch(v); setPage(1); },
    },
    {
      type: "select",
      value: role,
      onValueChange: (v) => { setRole(v); setPage(1); },
      options: roleOptions,
      placeholder: t("admin.users.allRoles"),
    },
    {
      type: "select",
      value: isActive,
      onValueChange: (v) => { setIsActive(v); setPage(1); },
      options: statusOptions,
      placeholder: t("admin.users.allStatuses"),
    },
  ];

  return (
    <DataTablePage
      title={t("admin.users.title")}
      columns={columns}
      data={data?.users ?? []}
      keyExtractor={(u) => u.id}
      onRowClick={(u) => navigate({ to: "/admin/users/$userId", params: { userId: u.id } })}
      isLoading={isLoading}
      filters={filters}
      pagination={data ? { page: data.page, perPage: data.per_page, total: data.total, onChange: setPage } : undefined}
      emptyIcon={UsersIcon}
      emptyTitle={t("admin.users.emptyTitle")}
      createAction={
        <Button onClick={() => setCreateOpen(true)}>
          {t("admin.users.add")}
        </Button>
      }
    >
      <CreateUserDialog open={createOpen} onOpenChange={setCreateOpen} />
    </DataTablePage>
  );
}

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
    { label: t("admin.users.roleAdmin"), value: "admin" },
    { label: t("admin.users.roleStaff"), value: "staff" },
    { label: t("admin.users.roleConsultant"), value: "consultant" },
    { label: t("admin.users.roleGuest"), value: "guest" },
  ];

  const officeOptions = [
    { label: "—", value: "" },
    ...(officesData ?? []).map((o) => ({ label: o.name, value: o.id })),
  ];

  async function handleSubmit(e: SyntheticEvent) {
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
        <FormField
          label={t("admin.users.firstName")}
          required
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
        required
        value={form.email}
        onChange={(v) => setForm((f) => ({ ...f, email: v }))}
      />

      <FormField
        label={t("admin.users.phone")}
        type="tel"
        value={form.phone}
        onChange={(v) => setForm((f) => ({ ...f, phone: v }))}
      />

      <FormField
        label={t("admin.users.password")}
        type="password"
        required
        value={form.password}
        onChange={(v) => setForm((f) => ({ ...f, password: v }))}
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
    </FormDialog>
  );
}
