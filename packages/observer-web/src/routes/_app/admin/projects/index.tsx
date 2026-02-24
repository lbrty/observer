import {
  FolderSimpleIcon,
  PencilSimpleIcon,
  UsersIcon,
} from "@/components/icons";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { DataTable, type Column } from "@/components/data-table";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { StatusBadge } from "@/components/status-badge";
import { useProjects } from "@/hooks/use-projects";
import type { Project } from "@/types/project";

export const Route = createFileRoute("/_app/admin/projects/")({
  component: ProjectsPage,
});

const statusTabs = ["", "active", "archived", "closed"] as const;

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

  const tabLabels: Record<string, string> = {
    "": t("admin.common.all"),
    active: "Active",
    archived: "Archived",
    closed: "Closed",
  };

  const columns: Column<Project>[] = [
    {
      key: "name",
      header: t("admin.projects.name"),
      render: (p) => (
        <div className="flex items-center gap-3">
          <span className="inline-flex size-8 shrink-0 items-center justify-center rounded-lg bg-bg-tertiary text-fg-tertiary">
            <FolderSimpleIcon size={16} />
          </span>
          <div className="min-w-0">
            <p className="truncate font-medium text-fg">{p.name}</p>
            {p.description && (
              <p className="truncate text-xs text-fg-tertiary">
                {p.description}
              </p>
            )}
          </div>
        </div>
      ),
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
          {new Date(p.created_at).toLocaleDateString("en-CA")}
        </span>
      ),
    },
    {
      key: "actions",
      header: "",
      render: (p) => (
        <div className="flex gap-2">
          <Link
            to="/admin/projects/$projectId/permissions"
            params={{ projectId: p.id }}
            onClick={(e) => e.stopPropagation()}
            className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-accent"
          >
            <UsersIcon size={16} />
          </Link>
          <Link
            to="/admin/projects/$projectId"
            params={{ projectId: p.id }}
            onClick={(e) => e.stopPropagation()}
            className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-accent"
          >
            <PencilSimpleIcon size={16} />
          </Link>
        </div>
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
            className="rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90"
          >
            + {t("admin.projects.newProject")}
          </Link>
        }
      />

      <div className="mb-4 flex gap-0 rounded-lg border border-border-secondary bg-bg-secondary p-0.5">
        {statusTabs.map((tab) => (
          <button
            key={tab}
            type="button"
            onClick={() => {
              setStatus(tab);
              setPage(1);
            }}
            className={`cursor-pointer rounded-md px-4 py-1.5 text-sm font-medium transition-colors ${
              status === tab
                ? "bg-bg shadow-card text-fg"
                : "text-fg-tertiary hover:text-fg"
            }`}
          >
            {tabLabels[tab]}
          </button>
        ))}
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
