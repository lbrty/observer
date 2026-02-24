import { PencilSimple, Trash } from "@phosphor-icons/react";
import { Dialog } from "@base-ui/react/dialog";
import { Field } from "@base-ui/react/field";
import { createFileRoute } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { PageHeader } from "@/components/page-header";
import {
  useCategories,
  useCreateCategory,
  useDeleteCategory,
  useUpdateCategory,
} from "@/hooks/use-categories";
import type { Category } from "@/types/reference";

export const Route = createFileRoute("/_app/admin/reference/categories")({
  component: CategoriesPage,
});

function CategoriesPage() {
  const { t } = useTranslation();
  const { data, isLoading } = useCategories();
  const createCategory = useCreateCategory();
  const updateCategory = useUpdateCategory();
  const deleteCategory = useDeleteCategory();

  const [createOpen, setCreateOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<Category | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Category | null>(null);

  const columns: Column<Category>[] = [
    {
      key: "name",
      header: t("admin.reference.categories.name"),
      render: (c) => <span className="font-medium text-fg">{c.name}</span>,
    },
    {
      key: "description",
      header: t("admin.reference.categories.description"),
      render: (c) => (
        <span className="text-fg-secondary">{c.description ?? "—"}</span>
      ),
    },
    {
      key: "actions",
      header: "",
      render: (c) => (
        <div className="flex gap-2">
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              setEditTarget(c);
            }}
            className="cursor-pointer rounded p-1 text-fg-tertiary hover:bg-bg-tertiary hover:text-fg"
          >
            <PencilSimple size={16} />
          </button>
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              setDeleteTarget(c);
            }}
            className="cursor-pointer rounded p-1 text-fg-tertiary hover:bg-bg-tertiary hover:text-rose"
          >
            <Trash size={16} />
          </button>
        </div>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t("admin.reference.categories.title")}
        action={
          <button
            type="button"
            onClick={() => setCreateOpen(true)}
            className="cursor-pointer rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-accent-fg hover:opacity-90"
          >
            {t("admin.reference.categories.add")}
          </button>
        }
      />

      <DataTable
        columns={columns}
        data={data?.categories ?? []}
        keyExtractor={(c) => c.id}
        isLoading={isLoading}
      />

      <CategoryFormDialog
        open={createOpen}
        onOpenChange={setCreateOpen}
        title={t("admin.reference.categories.addTitle")}
        onSubmit={async (name, description) => {
          await createCategory.mutateAsync({
            name,
            description: description || undefined,
          });
          setCreateOpen(false);
        }}
        loading={createCategory.isPending}
      />

      {editTarget && (
        <CategoryFormDialog
          open={!!editTarget}
          onOpenChange={(open) => !open && setEditTarget(null)}
          title={t("admin.reference.categories.editTitle")}
          initial={editTarget}
          onSubmit={async (name, description) => {
            await updateCategory.mutateAsync({
              id: editTarget.id,
              data: { name, description: description || undefined },
            });
            setEditTarget(null);
          }}
          loading={updateCategory.isPending}
        />
      )}

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => !open && setDeleteTarget(null)}
        title={t("admin.common.delete")}
        description={t("admin.reference.categories.deleteConfirm")}
        onConfirm={async () => {
          if (deleteTarget) {
            await deleteCategory.mutateAsync(deleteTarget.id);
            setDeleteTarget(null);
          }
        }}
        loading={deleteCategory.isPending}
      />
    </div>
  );
}

function CategoryFormDialog({
  open,
  onOpenChange,
  title,
  initial,
  onSubmit,
  loading,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  initial?: Category;
  onSubmit: (name: string, description: string) => Promise<void>;
  loading: boolean;
}) {
  const { t } = useTranslation();
  const [name, setName] = useState(initial?.name ?? "");
  const [description, setDescription] = useState(initial?.description ?? "");

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    await onSubmit(name, description);
    if (!initial) {
      setName("");
      setDescription("");
    }
  }

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Backdrop className="fixed inset-0 bg-black/40" />
        <Dialog.Popup className="fixed top-1/2 left-1/2 w-full max-w-sm -translate-x-1/2 -translate-y-1/2 rounded-lg border border-border-secondary bg-bg-secondary p-6 shadow-elevated">
          <Dialog.Title className="text-lg font-semibold text-fg">
            {title}
          </Dialog.Title>
          <form onSubmit={handleSubmit} className="mt-4 space-y-3">
            <Field.Root>
              <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("admin.reference.categories.name")}
              </Field.Label>
              <Field.Control
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="block w-full rounded-md border border-border-secondary bg-bg px-3 py-2 text-sm text-fg outline-none focus:border-accent"
              />
            </Field.Root>
            <Field.Root>
              <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
                {t("admin.reference.categories.description")}
              </Field.Label>
              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={2}
                className="block w-full rounded-md border border-border-secondary bg-bg px-3 py-2 text-sm text-fg outline-none focus:border-accent"
              />
            </Field.Root>
            <div className="flex justify-end gap-3 pt-2">
              <Dialog.Close className="cursor-pointer rounded-md border border-border-secondary px-3 py-1.5 text-sm text-fg-secondary hover:bg-bg-tertiary">
                {t("admin.common.cancel")}
              </Dialog.Close>
              <button
                type="submit"
                disabled={loading}
                className="cursor-pointer rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-accent-fg hover:opacity-90 disabled:opacity-50"
              >
                {loading ? t("admin.common.saving") : t("admin.common.save")}
              </button>
            </div>
          </form>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
