import { Field } from "@base-ui/react/field";
import { createFileRoute, Link } from "@tanstack/react-router";
import { type FormEvent, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { PageHeader } from "@/components/page-header";
import { UISelect } from "@/components/ui-select";
import { useProject, useUpdateProject } from "@/hooks/use-projects";

export const Route = createFileRoute("/_app/admin/projects/$projectId/")({
  component: ProjectDetailPage,
});

function ProjectDetailPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const { data: project, isLoading } = useProject(projectId);
  const updateProject = useUpdateProject();

  const [form, setForm] = useState({
    name: "",
    description: "",
    status: "active" as string,
  });

  useEffect(() => {
    if (project) {
      setForm({
        name: project.name,
        description: project.description ?? "",
        status: project.status,
      });
    }
  }, [project]);

  if (isLoading) return null;
  if (!project) return null;

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    await updateProject.mutateAsync({
      id: projectId,
      data: {
        name: form.name,
        description: form.description || undefined,
        status: form.status as "active" | "archived" | "closed",
      },
    });
  }

  return (
    <div>
      <PageHeader
        title={t("admin.projects.editTitle")}
        action={
          <Link
            to="/admin/projects/$projectId/permissions"
            params={{ projectId }}
            className="rounded-lg border border-border-secondary px-4 py-2 text-sm text-fg-secondary shadow-card hover:bg-bg-tertiary"
          >
            {t("admin.projects.permissions")}
          </Link>
        }
      />

      <form onSubmit={handleSubmit} className="max-w-lg space-y-4">
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.projects.name")}
          </Field.Label>
          <Field.Control
            required
            value={form.name}
            onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
            className="block w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent"
          />
        </Field.Root>

        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.projects.description")}
          </Field.Label>
          <textarea
            value={form.description}
            onChange={(e) =>
              setForm((f) => ({ ...f, description: e.target.value }))
            }
            rows={3}
            className="block w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent"
          />
        </Field.Root>

        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.projects.status")}
          </Field.Label>
          <UISelect
            value={form.status}
            onValueChange={(v) => setForm((f) => ({ ...f, status: v }))}
            options={[
              { label: "active", value: "active" },
              { label: "archived", value: "archived" },
              { label: "closed", value: "closed" },
            ]}
            fullWidth
          />
        </Field.Root>

        <button
          type="submit"
          disabled={updateProject.isPending}
          className="cursor-pointer rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90 disabled:opacity-50"
        >
          {updateProject.isPending
            ? t("admin.projects.saving")
            : t("admin.projects.save")}
        </button>
      </form>
    </div>
  );
}
