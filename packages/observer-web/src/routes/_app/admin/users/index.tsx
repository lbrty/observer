import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { DataTable, type Column } from "@/components/data-table";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { StatusBadge } from "@/components/status-badge";
import { useUsers } from "@/hooks/use-users";
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
  const [isActive, setIsActive] = useState<string>("");

  const params = {
    page,
    per_page: 20,
    ...(search && { search }),
    ...(role && { role }),
    ...(isActive !== "" && { is_active: isActive === "true" }),
  };

  const { data, isLoading } = useUsers(params);

  const columns: Column<AdminUser>[] = [
    {
      key: "name",
      header: t("admin.users.name"),
      render: (u) => (
        <span className="font-medium text-fg">
          {u.first_name} {u.last_name}
        </span>
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
      render: (u) => (
        <StatusBadge
          label={u.is_active ? t("admin.users.yes") : t("admin.users.no")}
          variant={u.is_active ? "foam" : "neutral"}
        />
      ),
    },
    {
      key: "verified",
      header: t("admin.users.verified"),
      render: (u) => (
        <StatusBadge
          label={u.is_verified ? t("admin.users.yes") : t("admin.users.no")}
          variant={u.is_verified ? "foam" : "neutral"}
        />
      ),
    },
    {
      key: "created",
      header: t("admin.users.created"),
      render: (u) => (
        <span className="text-fg-tertiary">
          {new Date(u.created_at).toLocaleDateString()}
        </span>
      ),
    },
  ];

  return (
    <div>
      <PageHeader title={t("admin.users.title")} />

      <div className="mb-4 flex gap-3">
        <input
          type="text"
          placeholder={t("admin.users.search")}
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            setPage(1);
          }}
          className="rounded-md border border-border-secondary bg-bg-secondary px-3 py-1.5 text-sm text-fg outline-none focus:border-accent"
        />
        <select
          value={role}
          onChange={(e) => {
            setRole(e.target.value);
            setPage(1);
          }}
          className="rounded-md border border-border-secondary bg-bg-secondary pl-3 pr-1 py-1.5 text-sm text-fg outline-none"
        >
          <option value="">{t("admin.users.allRoles")}</option>
          <option value="admin">admin</option>
          <option value="staff">staff</option>
          <option value="consultant">consultant</option>
          <option value="guest">guest</option>
        </select>
        <select
          value={isActive}
          onChange={(e) => {
            setIsActive(e.target.value);
            setPage(1);
          }}
          className="rounded-md border border-border-secondary bg-bg-secondary pl-3 pr-1 py-1.5 text-sm text-fg outline-none"
        >
          <option value="">{t("admin.users.allStatuses")}</option>
          <option value="true">{t("admin.users.active")}</option>
          <option value="false">{t("admin.users.no")}</option>
        </select>
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
    </div>
  );
}
