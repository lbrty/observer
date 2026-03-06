import {
  ArrowSquareOutIcon,
  FolderSimpleIcon,
  PencilSimpleIcon,
  PlusIcon,
  UsersIcon,
} from "@/components/icons";
import { Field } from "@base-ui/react/field";
import { Tabs } from "@base-ui/react/tabs";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { type SyntheticEvent, useState } from "react";

import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import { FormDialog } from "@/components/form-dialog";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { StatusBadge } from "@/components/status-badge";
import { useCreateProject, useProjects } from "@/hooks/use-projects";
import type { Project } from "@/types/project";

export const Route = createFileRoute("/_app/admin/projects/")({
  component: ProjectsPage,
  validateSearch: (search: Record<string, unknown>): { status?: string; page?: number } => ({
    status: (search.status as string) || undefined,
    page: Number(search.page) || undefined,
  }),
});

const statusTabs = ["", "active", "archived", "closed"] as const;

function ProjectsPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { status = "", page = 1 } = Route.useSearch();

  const [createOpen, setCreateOpen] = useState(false);

  function setStatus(value: string) {
    navigate({ from: Route.fullPath, search: { status: value || undefined }, replace: true });
  }

  function setPage(value: number) {
    navigate({ from: Route.fullPath, search: (prev) => ({ ...prev, page: value > 1 ? value : undefined }), replace: true });
  }

  const params = {
    page,
    per_page: 20,
    ...(status && { status }),
  };

  const { data, isLoading } = useProjects(params);

  const tabLabels: Record<string, string> = {
    "": t("admin.common.all"),
    active: t("admin.projects.statusActive"),
    archived: t("admin.projects.statusArchived"),
    closed: t("admin.projects.statusClosed"),
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
            {p.description && <p className="truncate text-xs text-fg-tertiary">{p.description}</p>}
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
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {new Date(p.created_at).toLocaleDateString("en-CA")}
        </span>
      ),
    },
    {
      key: "actions",
      header: "",
      render: (p) => (
        <div className="flex gap-1">
          <Button variant="ghost" className="p-1.5" asChild>
            <Link
              to="/projects/$projectId/people"
              params={{ projectId: p.id }}
              onClick={(e) => e.stopPropagation()}
              title={t("admin.projects.browse")}
            >
              <ArrowSquareOutIcon size={16} />
            </Link>
          </Button>
          <Button variant="ghost" className="p-1.5" asChild>
            <Link
              to="/admin/projects/$projectId/permissions"
              params={{ projectId: p.id }}
              onClick={(e) => e.stopPropagation()}
            >
              <UsersIcon size={16} />
            </Link>
          </Button>
          <Button variant="ghost" className="p-1.5" asChild>
            <Link
              to="/admin/projects/$projectId"
              params={{ projectId: p.id }}
              onClick={(e) => e.stopPropagation()}
            >
              <PencilSimpleIcon size={16} />
            </Link>
          </Button>
        </div>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t("admin.projects.title")}
        action={
          <Button icon={<PlusIcon size={16} />} onClick={() => setCreateOpen(true)}>
            {t("admin.projects.newProject")}
          </Button>
        }
      />

      <Tabs.Root
        defaultValue=""
        value={status}
        onValueChange={(value) => {
          setStatus(value as string);
        }}
        className="mb-4"
      >
        <Tabs.List className="flex gap-0 rounded-lg border border-border-secondary bg-bg-secondary p-0.5">
          {statusTabs.map((tab) => (
            <Tabs.Tab
              key={tab}
              value={tab}
              className="cursor-pointer rounded-sm px-4 py-1.5 m-0.5 text-sm font-medium text-fg-tertiary transition-colors hover:text-fg data-active:bg-bg data-active:text-fg data-active:shadow-card"
            >
              {tabLabels[tab]}
            </Tabs.Tab>
          ))}
        </Tabs.List>
      </Tabs.Root>

      <DataTable
        columns={columns}
        data={data?.projects ?? []}
        keyExtractor={(p) => p.id}
        onRowClick={(p) =>
          navigate({
            to: "/projects/$projectId/people",
            params: { projectId: p.id },
          })
        }
        isLoading={isLoading}
        emptyState={
          <EmptyState
            icon={FolderSimpleIcon}
            title={t("admin.projects.emptyTitle")}
            description={t("admin.projects.emptyDescription")}
            action={
              <Button icon={<PlusIcon size={16} />} onClick={() => setCreateOpen(true)}>
                {t("admin.projects.newProject")}
              </Button>
            }
          />
        }
      />

      {data && (
        <Pagination
          page={data.page}
          perPage={data.per_page}
          total={data.total}
          onChange={setPage}
        />
      )}

      <CreateProjectDialog
        open={createOpen}
        onOpenChange={setCreateOpen}
        onCreated={(project) => {
          setCreateOpen(false);
          navigate({
            to: "/admin/projects/$projectId",
            params: { projectId: project.id },
          });
        }}
      />
    </div>
  );
}

function CreateProjectDialog({
  open,
  onOpenChange,
  onCreated,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onCreated: (project: Project) => void;
}) {
  const { t } = useTranslation();
  const createProject = useCreateProject();
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    const project = await createProject.mutateAsync({
      name,
      description: description || undefined,
    });
    setName("");
    setDescription("");
    onCreated(project);
  }

  return (
    <FormDialog
      open={open}
      onOpenChange={onOpenChange}
      title={t("admin.projects.createTitle")}
      loading={createProject.isPending}
      onSubmit={handleSubmit}
      maxWidth="md"
    >
      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("admin.projects.name")}
        </Field.Label>
        <Field.Control
          required
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="block w-full rounded-lg border border-border-secondary bg-bg-secondary h-9 px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
        />
      </Field.Root>

      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("admin.projects.description")}
        </Field.Label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          rows={3}
          className="block w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
        />
      </Field.Root>
    </FormDialog>
  );
}
