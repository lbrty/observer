import { type FormEvent, useState } from "react";

import { Field } from "@base-ui/react/field";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { FormDialog } from "@/components/form-dialog";
import { TagIcon, TrashIcon } from "@/components/icons";
import { PageHeader } from "@/components/page-header";
import { useCreateTag, useDeleteTag, useTags } from "@/hooks/use-tags";
import { HTTPError } from "@/lib/api";
import type { Tag } from "@/types/tag";

export const Route = createFileRoute("/_app/projects/$projectId/tags/")({
  component: TagsPage,
});

function TagsPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();

  const { data, isLoading } = useTags(projectId);
  const createTag = useCreateTag(projectId);
  const deleteTag = useDeleteTag(projectId);

  const [createOpen, setCreateOpen] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState<Tag | null>(null);
  const [name, setName] = useState("");
  const [error, setError] = useState("");

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    setError("");
    try {
      await createTag.mutateAsync({ name: name.trim() });
      setName("");
      setCreateOpen(false);
    } catch (err) {
      if (err instanceof HTTPError) {
        const body = await err.response.json().catch(() => null);
        const code = body?.code;
        const translated = code ? t(code, { defaultValue: "" }) : "";
        setError(translated || body?.error || err.message);
      } else {
        setError(t("common.unexpectedError"));
      }
    }
  }

  async function handleDelete() {
    if (!deleteTarget) return;
    try {
      await deleteTag.mutateAsync(deleteTarget.id);
      setDeleteTarget(null);
    } catch (err) {
      if (err instanceof HTTPError) {
        const body = await err.response.json().catch(() => null);
        const code = body?.code;
        const translated = code ? t(code, { defaultValue: "" }) : "";
        setError(translated || body?.error || err.message);
      } else {
        setError(t("common.unexpectedError"));
      }
    }
  }

  const columns: Column<Tag>[] = [
    {
      key: "name",
      header: t("project.tags.name"),
      render: (tag) => (
        <div className="flex items-center gap-2.5">
          <span className="inline-flex size-7 shrink-0 items-center justify-center rounded-md bg-bg-tertiary text-fg-tertiary">
            <TagIcon size={14} />
          </span>
          <span className="font-medium text-fg">{tag.name}</span>
        </div>
      ),
    },
    {
      key: "created_at",
      header: t("admin.common.createdAt"),
      render: (tag) => (
        <span className="font-mono text-xs tabular-nums text-fg-tertiary">
          {new Date(tag.created_at).toLocaleDateString("en-CA")}
        </span>
      ),
    },
    {
      key: "actions",
      header: "",
      render: (tag) => (
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            setDeleteTarget(tag);
          }}
          className="cursor-pointer rounded-lg p-1.5 text-fg-tertiary hover:bg-bg-tertiary hover:text-rose"
        >
          <TrashIcon size={16} />
        </button>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t("project.tags.title")}
        action={
          <button
            type="button"
            onClick={() => {
              setError("");
              setName("");
              setCreateOpen(true);
            }}
            className="cursor-pointer rounded-lg bg-accent px-4 py-2 text-sm font-medium text-accent-fg shadow-card hover:opacity-90"
          >
            + {t("project.tags.add")}
          </button>
        }
      />

      <DataTable
        columns={columns}
        data={data?.tags ?? []}
        keyExtractor={(tag) => tag.id}
        isLoading={isLoading}
      />

      <FormDialog
        open={createOpen}
        onOpenChange={setCreateOpen}
        title={t("project.tags.add")}
        loading={createTag.isPending}
        onSubmit={handleCreate}
      >
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.tags.name")}
          </Field.Label>
          <Field.Control
            required
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="block h-9 w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent"
          />
        </Field.Root>
        {error && (
          <p className="text-sm text-rose">{error}</p>
        )}
      </FormDialog>

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => !open && setDeleteTarget(null)}
        title={t("admin.common.delete")}
        description={t("project.tags.deleteConfirm")}
        onConfirm={handleDelete}
        loading={deleteTag.isPending}
      />
    </div>
  );
}
