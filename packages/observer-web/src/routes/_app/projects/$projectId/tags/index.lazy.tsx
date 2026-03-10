import { type SyntheticEvent, useState } from "react";

import { Field } from "@base-ui/react/field";
import { createLazyFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { ConfirmDialog } from "@/components/confirm-dialog";
import { type Column } from "@/components/data-table";
import { DataTablePage } from "@/components/data-table-page";
import { FormDialog } from "@/components/form-dialog";
import {
  ArrowsClockwiseIcon,
  PencilSimpleIcon,
  PlusIcon,
  TagIcon,
  TrashIcon,
} from "@/components/icons";
import { useCreateTag, useDeleteTag, useUpdateTag, useTags } from "@/hooks/use-tags";
import { useProjectRole } from "@/hooks/use-project-role";
import { handleApiError } from "@/lib/form-error";
import { resolveTagColor } from "@/lib/tag-color";
import { useToast } from "@/stores/toast";
import type { Tag } from "@/types/tag";

export const Route = createLazyFileRoute("/_app/projects/$projectId/tags/")({
  component: TagsPage,
});

function randomHex(): string {
  const hex = Math.floor(Math.random() * 0xffffff)
    .toString(16)
    .padStart(6, "0");
  return `#${hex}`;
}

function TagsPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const toast = useToast();

  const { canWrite, canDelete } = useProjectRole(projectId);
  const { data, isLoading } = useTags(projectId);
  const createTag = useCreateTag(projectId);
  const updateTag = useUpdateTag(projectId);
  const deleteTag = useDeleteTag(projectId);

  const [formOpen, setFormOpen] = useState(false);
  const [editTag, setEditTag] = useState<Tag | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Tag | null>(null);
  const [name, setName] = useState("");
  const [color, setColor] = useState("");
  const [error, setError] = useState("");

  function openCreate() {
    setEditTag(null);
    setName("");
    setColor("");
    setError("");
    setFormOpen(true);
  }

  function openEdit(tag: Tag) {
    setEditTag(tag);
    setName(tag.name);
    setColor(resolveTagColor(tag.color, tag.name));
    setError("");
    setFormOpen(true);
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    setError("");
    try {
      if (editTag) {
        await updateTag.mutateAsync({
          id: editTag.id,
          data: { name: name.trim(), color },
        });
        toast.success(t("project.tags.saved"));
      } else {
        const tagColor = resolveTagColor(color, name.trim());
        await createTag.mutateAsync({ name: name.trim(), color: tagColor });
        toast.success(t("project.tags.saved"));
      }
      setFormOpen(false);
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  async function handleDelete() {
    if (!deleteTarget) return;
    try {
      await deleteTag.mutateAsync(deleteTarget.id);
      setDeleteTarget(null);
      toast.success(t("project.tags.deleted"));
    } catch (err) {
      setError(await handleApiError(err, t));
    }
  }

  const columns: Column<Tag>[] = [
    {
      key: "name",
      header: t("project.tags.name"),
      render: (tag) => (
        <div className="flex items-center gap-2.5">
          <span
            className="inline-flex size-7 shrink-0 items-center justify-center rounded-md"
            style={{ backgroundColor: resolveTagColor(tag.color, tag.name) }}
          >
            <TagIcon size={14} className="text-white" />
          </span>
          <span className="font-medium text-fg">{tag.name}</span>
        </div>
      ),
    },
    {
      key: "color",
      header: t("project.tags.color"),
      render: (tag) => {
        const hex = resolveTagColor(tag.color, tag.name);
        return (
          <div className="flex items-center gap-2">
            <span
              className="inline-block size-4 rounded-full border border-border-secondary"
              style={{ backgroundColor: hex }}
            />
            <span className="font-mono text-xs text-fg-tertiary">{hex}</span>
            {canWrite && (
              <button
                type="button"
                onClick={(e) => {
                  e.stopPropagation();
                  const newColor = randomHex();
                  updateTag.mutate(
                    { id: tag.id, data: { color: newColor } },
                    { onSuccess: () => toast.success(t("project.tags.saved")) },
                  );
                }}
                className="ml-1 inline-flex size-6 cursor-pointer items-center justify-center rounded-md text-fg-tertiary transition-colors hover:bg-bg-tertiary hover:text-fg"
              >
                <ArrowsClockwiseIcon size={14} />
              </button>
            )}
          </div>
        );
      },
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
    ...(canWrite || canDelete
      ? [
          {
            key: "actions",
            header: "",
            render: (tag: Tag) => (
              <div className="flex gap-1">
                {canWrite && (
                  <Button
                    variant="ghost"
                    className="p-1.5"
                    onClick={(e) => {
                      e.stopPropagation();
                      openEdit(tag);
                    }}
                  >
                    <PencilSimpleIcon size={16} />
                  </Button>
                )}
                {canDelete && (
                  <Button
                    variant="ghost"
                    className="p-1.5 hover:text-rose"
                    onClick={(e) => {
                      e.stopPropagation();
                      setDeleteTarget(tag);
                    }}
                  >
                    <TrashIcon size={16} />
                  </Button>
                )}
              </div>
            ),
          } satisfies Column<Tag>,
        ]
      : []),
  ];

  return (
    <DataTablePage
      title={t("project.tags.title")}
      columns={columns}
      data={data?.tags ?? []}
      keyExtractor={(tag) => tag.id}
      isLoading={isLoading}
      emptyIcon={TagIcon}
      emptyTitle={t("project.tags.emptyTitle")}
      emptyDescription={t("project.tags.emptyDescription")}
      emptyAction={
        canWrite ? (
          <Button onClick={openCreate} icon={<PlusIcon size={16} />}>
            {t("project.tags.add")}
          </Button>
        ) : undefined
      }
      createAction={
        canWrite ? (
          <Button icon={<PlusIcon size={16} />} onClick={openCreate}>
            {t("project.tags.add")}
          </Button>
        ) : undefined
      }
    >
      <FormDialog
        open={formOpen}
        onOpenChange={setFormOpen}
        title={editTag ? t("project.tags.edit") : t("project.tags.add")}
        loading={createTag.isPending || updateTag.isPending}
        onSubmit={handleSubmit}
      >
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.tags.name")}
          </Field.Label>
          <Field.Control
            required
            value={name}
            onChange={(e) => {
              setName(e.target.value);
              if (!editTag && !color) {
                // preview color will update automatically
              }
            }}
            className="block h-9 w-full rounded-lg border border-border-secondary bg-bg-secondary px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
          />
        </Field.Root>

        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("project.tags.color")}
          </Field.Label>
          <div className="flex items-center gap-3">
            <input
              type="color"
              value={resolveTagColor(color, name || "tag")}
              onChange={(e) => setColor(e.target.value)}
              className="size-9 cursor-pointer rounded-lg border border-border-secondary bg-bg-secondary p-0.5"
            />
            <button
              type="button"
              onClick={() => setColor(randomHex())}
              className="inline-flex size-9 cursor-pointer items-center justify-center rounded-lg border border-border-secondary bg-bg-secondary text-fg-tertiary transition-colors hover:text-fg"
            >
              <ArrowsClockwiseIcon size={16} />
            </button>
            <span className="font-mono text-sm text-fg-tertiary">
              {resolveTagColor(color, name || "tag")}
            </span>
          </div>
        </Field.Root>

        {error && <p className="text-sm text-rose">{error}</p>}
      </FormDialog>

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => !open && setDeleteTarget(null)}
        title={t("admin.common.delete")}
        description={t("project.tags.deleteConfirm")}
        onConfirm={handleDelete}
        loading={deleteTag.isPending}
      />
    </DataTablePage>
  );
}
