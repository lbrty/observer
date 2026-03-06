import { type FormEvent, useState } from "react";

import { Field } from "@base-ui/react/field";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import { FormDialog } from "@/components/form-dialog";
import {
  ArrowsClockwiseIcon,
  PencilSimpleIcon,
  PlusIcon,
  TagIcon,
  TrashIcon,
} from "@/components/icons";
import { PageHeader } from "@/components/page-header";
import { useCreateTag, useDeleteTag, useUpdateTag, useTags } from "@/hooks/use-tags";
import { HTTPError } from "@/lib/api";
import { useToast } from "@/stores/toast";
import type { Tag } from "@/types/tag";

export const Route = createFileRoute("/_app/projects/$projectId/tags/")({
  component: TagsPage,
});

function colorFromName(name: string): string {
  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash);
  }
  const h = ((hash % 360) + 360) % 360;
  return `hsl(${h}, 55%, 55%)`;
}

function randomHex(): string {
  const hex = Math.floor(Math.random() * 0xffffff)
    .toString(16)
    .padStart(6, "0");
  return `#${hex}`;
}

function hslToHex(hsl: string): string {
  const match = hsl.match(/hsl\((\d+),\s*(\d+)%,\s*(\d+)%\)/);
  if (!match) return "#888888";
  const h = Number(match[1]) / 360;
  const s = Number(match[2]) / 100;
  const l = Number(match[3]) / 100;

  const hue2rgb = (p: number, q: number, t: number) => {
    if (t < 0) t += 1;
    if (t > 1) t -= 1;
    if (t < 1 / 6) return p + (q - p) * 6 * t;
    if (t < 1 / 2) return q;
    if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6;
    return p;
  };

  const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
  const p = 2 * l - q;
  const r = Math.round(hue2rgb(p, q, h + 1 / 3) * 255);
  const g = Math.round(hue2rgb(p, q, h) * 255);
  const b = Math.round(hue2rgb(p, q, h - 1 / 3) * 255);

  return `#${((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)}`;
}

function TagsPage() {
  const { t } = useTranslation();
  const { projectId } = Route.useParams();
  const toast = useToast();

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
    setColor(tag.color || hslToHex(colorFromName(tag.name)));
    setError("");
    setFormOpen(true);
  }

  async function handleSubmit(e: FormEvent) {
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
        const tagColor = color || hslToHex(colorFromName(name.trim()));
        await createTag.mutateAsync({ name: name.trim(), color: tagColor });
        toast.success(t("project.tags.saved"));
      }
      setFormOpen(false);
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
      toast.success(t("project.tags.deleted"));
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
          <span
            className="inline-flex size-7 shrink-0 items-center justify-center rounded-md"
            style={{ backgroundColor: tag.color || hslToHex(colorFromName(tag.name)) }}
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
        const hex = tag.color || hslToHex(colorFromName(tag.name));
        return (
          <div className="flex items-center gap-2">
            <span
              className="inline-block size-4 rounded-full border border-border-secondary"
              style={{ backgroundColor: hex }}
            />
            <span className="font-mono text-xs text-fg-tertiary">{hex}</span>
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
    {
      key: "actions",
      header: "",
      render: (tag) => (
        <div className="flex gap-1">
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
        </div>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t("project.tags.title")}
        action={
          <Button icon={<PlusIcon size={16} />} onClick={openCreate}>
            {t("project.tags.add")}
          </Button>
        }
      />

      <DataTable
        columns={columns}
        data={data?.tags ?? []}
        keyExtractor={(tag) => tag.id}
        isLoading={isLoading}
        emptyState={
          <EmptyState
            icon={TagIcon}
            title={t("project.tags.emptyTitle")}
            description={t("project.tags.emptyDescription")}
            action={
              <Button onClick={openCreate} icon={<PlusIcon size={16} />}>
                {t("project.tags.add")}
              </Button>
            }
          />
        }
      />

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
              value={color || hslToHex(colorFromName(name || "tag"))}
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
              {color || hslToHex(colorFromName(name || "tag"))}
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
    </div>
  );
}
