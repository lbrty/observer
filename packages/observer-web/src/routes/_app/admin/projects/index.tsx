import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { DataTable, type Column } from "@/components/data-table";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { StatusBadge } from "@/components/status-badge";
import { UISelect } from "@/components/ui-select";
import { useProjects } from "@/hooks/use-projects";
import type { Project } from "@/types/project";

export const Route = createFileRoute("/_app/admin/projects/")({
  component: ProjectsPage,
});

function ProjectsPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const [page, setPage] = useState(1);
  const [status, setStatus] = useState("");

  const params = {
    page,
    per_page: 20,
    ...(status && { status }),
  };

  const { data, isLoading } = useProjects(params);

  const columns: Column<Project>[] = [
    {
      key: "name",
      header: t("admin.projects.name"),
      render: (p) => <span className="font-medium text-fg">{p.name}</span>,
    },
    {
      key: "status",
      header: t("admin.projects.status"),
      render: (p) => <StatusBadge label={p.status} />,
    },
    {
      key: "created",
      header: t("admin.projects.created"),
      render: (p) => (
        <span className="text-fg-tertiary">
          {new Date(p.created_at).toLocaleDateString()}
        </span>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t("admin.projects.title")}
        action={
          <Link
            to="/admin/projects/new"
            className="rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-accent-fg hover:opacity-90"
          >
            {t("admin.projects.newProject")}
          </Link>
        }
      />

      <div className="mb-4">
        <UISelect
          value={status}
          onValueChange={(v) => {
            setStatus(v);
            setPage(1);
          }}
          options={[
            { label: t("admin.projects.allStatuses"), value: "" },
            { label: "active", value: "active" },
            { label: "archived", value: "archived" },
            { label: "closed", value: "closed" },
          ]}
          placeholder={t("admin.projects.allStatuses")}
        />
      </div>

      <DataTable
        columns={columns}
        data={data?.projects ?? []}
        keyExtractor={(p) => p.id}
        onRowClick={(p) =>
          navigate({
            to: "/admin/projects/$projectId",
            params: { projectId: p.id },
          })
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
