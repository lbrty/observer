import { Field } from "@base-ui/react/field";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { PageHeader } from "@/components/page-header";
import { useCreateProject } from "@/hooks/use-projects";

export const Route = createFileRoute("/_app/admin/projects/new")({
  component: NewProjectPage,
});

function NewProjectPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const createProject = useCreateProject();

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    const project = await createProject.mutateAsync({
      name,
      description: description || undefined,
    });
    navigate({
      to: "/admin/projects/$projectId",
      params: { projectId: project.id },
    });
  }

  return (
    <div>
      <PageHeader title={t("admin.projects.createTitle")} />
      <form onSubmit={handleSubmit} className="max-w-lg space-y-4">
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.projects.name")}
          </Field.Label>
          <Field.Control
            required
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="block w-full rounded-md border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent"
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
            className="block w-full rounded-md border border-border-secondary bg-bg-secondary px-3 py-2 text-sm text-fg outline-none focus:border-accent"
          />
        </Field.Root>

        <button
          type="submit"
          disabled={createProject.isPending}
          className="cursor-pointer rounded-md bg-accent px-4 py-2 text-sm font-medium text-accent-fg hover:opacity-90 disabled:opacity-50"
        >
          {createProject.isPending
            ? t("admin.projects.saving")
            : t("admin.projects.save")}
        </button>
      </form>
    </div>
  );
}
