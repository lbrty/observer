import {
  ArrowSquareOutIcon,
  FolderSimpleIcon,
  PencilSimpleIcon,
  UsersIcon,
} from "@/components/icons";
import { Field } from "@base-ui/react/field";
import { Tabs } from "@base-ui/react/tabs";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { DataTable, type Column } from "@/components/data-table";
import { FormDialog } from "@/components/form-dialog";
import { PageHeader } from "@/components/page-header";
import { Pagination } from "@/components/pagination";
import { StatusBadge } from "@/components/status-badge";
import { useCreateProject, useProjects } from "@/hooks/use-projects";
import type { Project } from "@/types/project";

export const Route = createFileRoute("/_app/admin/projects/")({
  component: ProjectsPage,
});

const statusTabs = ["", "active", "archived", "closed"] as const;

function ProjectsPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const [page, setPage] = useState(1);
  const [status, setStatus] = useState<string>("");
  const [createOpen, setCreateOpen] = useState(false);

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
          <Link
            to="/projects/$projectId/people"
            params={{ projectId: p.id }}
            onClick={(e) => e.stopPropagation()}
            className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-accent"
            title={t("admin.projects.browse")}
          >
            <ArrowSquareOutIcon size={16} />
          </Link>
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
          <button
            type="button"
            onClick={() => setCreateOpen(true)}
            className="cursor-pointer rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90"
          >
            + {t("admin.projects.newProject")}
          </button>
        }
      />

      <Tabs.Root
        defaultValue=""
        value={status}
        onValueChange={(value) => {
          setStatus(value as string);
          setPage(1);
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

  async function handleSubmit(e: FormEvent) {
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
          className="block w-full rounded-lg border border-border-secondary bg-bg-secondary h-9 px-3 text-sm text-fg outline-none focus:border-accent"
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
          className="block w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent"
        />
      </Field.Root>
    </FormDialog>
  );
}
