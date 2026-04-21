import { useState } from "react";

import { createLazyFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { CreateUserDialog } from "@/components/users/create-user-dialog";
import { DataTablePage } from "@/components/table/data-table-page";
import type { FilterDef } from "@/components/forms/filter-bar";
import { UsersIcon } from "@/components/ui/icons";
import { buildUsersColumns } from "@/components/users/users-columns";
import { useUsers } from "@/hooks/users/use-users";

export const Route = createLazyFileRoute("/_app/admin/users/")({
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

  const columns = buildUsersColumns({ t });

  const filters: FilterDef[] = [
    {
      type: "search",
      placeholder: t("admin.users.search"),
      value: search,
      onChange: (v) => {
        setSearch(v);
        setPage(1);
      },
    },
    {
      type: "select",
      value: role,
      onValueChange: (v) => {
        setRole(v);
        setPage(1);
      },
      options: roleOptions,
      placeholder: t("admin.users.allRoles"),
    },
    {
      type: "select",
      value: isActive,
      onValueChange: (v) => {
        setIsActive(v);
        setPage(1);
      },
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
      pagination={
        data
          ? { page: data.page, perPage: data.per_page, total: data.total, onChange: setPage }
          : undefined
      }
      emptyIcon={UsersIcon}
      emptyTitle={t("admin.users.emptyTitle")}
      createAction={<Button onClick={() => setCreateOpen(true)}>{t("admin.users.add")}</Button>}
    >
      <CreateUserDialog open={createOpen} onOpenChange={setCreateOpen} />
    </DataTablePage>
  );
}
