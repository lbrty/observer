import { type FormEvent, useState } from "react";

import { Field } from "@base-ui/react/field";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { FormDialog } from "@/components/form-dialog";
import { PageHeader } from "@/components/page-header";
import { RowActions } from "@/components/row-actions";
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
      render: (c) => <span className="text-fg-secondary">{c.description ?? "—"}</span>,
    },
    {
      key: "actions",
      header: "",
      render: (c) => (
        <RowActions onEdit={() => setEditTarget(c)} onDelete={() => setDeleteTarget(c)} />
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t("admin.reference.categories.title")}
        action={
          <Button onClick={() => setCreateOpen(true)}>
            {t("admin.reference.categories.add")}
          </Button>
        }
      />

      <DataTable
        columns={columns}
        data={data ?? []}
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
    <FormDialog
      open={open}
      onOpenChange={onOpenChange}
      title={title}
      loading={loading}
      onSubmit={handleSubmit}
    >
      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("admin.reference.categories.name")}
        </Field.Label>
        <Field.Control
          required
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent"
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
          className="block w-full rounded-lg border border-border-secondary bg-bg px-3 py-2 text-sm text-fg outline-none focus:border-accent"
        />
      </Field.Root>
    </FormDialog>
  );
}
