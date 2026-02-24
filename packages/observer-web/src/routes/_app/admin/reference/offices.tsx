import { Dialog } from "@base-ui/react/dialog";
import { Field } from "@base-ui/react/field";
import { createFileRoute } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { PageHeader } from "@/components/page-header";
import {
  useCreateOffice,
  useDeleteOffice,
  useOffices,
  useUpdateOffice,
} from "@/hooks/use-offices";
import type { Office } from "@/types/reference";

export const Route = createFileRoute("/_app/admin/reference/offices")({
  component: OfficesPage,
});

function OfficesPage() {
  const { t } = useTranslation();
  const { data, isLoading } = useOffices();
  const createOffice = useCreateOffice();
  const updateOffice = useUpdateOffice();
  const deleteOffice = useDeleteOffice();

  const [createOpen, setCreateOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<Office | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Office | null>(null);

  const columns: Column<Office>[] = [
    {
      key: "name",
      header: t("admin.reference.offices.name"),
      render: (o) => <span className="font-medium text-fg">{o.name}</span>,
    },
    {
      key: "place_id",
      header: t("admin.reference.offices.place"),
      render: (o) => (
        <span className="font-mono text-xs text-fg-secondary">
          {o.place_id ?? "—"}
        </span>
      ),
    },
    {
      key: "actions",
      header: t("admin.common.actions"),
      render: (o) => (
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => setEditTarget(o)}
            className="cursor-pointer text-xs text-accent hover:underline"
          >
            {t("admin.common.edit")}
          </button>
          <button
            type="button"
            onClick={() => setDeleteTarget(o)}
            className="cursor-pointer text-xs text-rose hover:underline"
          >
            {t("admin.common.delete")}
          </button>
        </div>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t("admin.reference.offices.title")}
        action={
          <button
            type="button"
            onClick={() => setCreateOpen(true)}
            className="cursor-pointer rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-accent-fg hover:opacity-90"
          >
            {t("admin.reference.offices.add")}
          </button>
        }
      />

      <DataTable
        columns={columns}
        data={data?.offices ?? []}
        keyExtractor={(o) => o.id}
        isLoading={isLoading}
      />

      <OfficeFormDialog
        open={createOpen}
        onOpenChange={setCreateOpen}
        title={t("admin.reference.offices.addTitle")}
        onSubmit={async (name, placeId) => {
          await createOffice.mutateAsync({
            name,
            place_id: placeId || undefined,
          });
          setCreateOpen(false);
        }}
        loading={createOffice.isPending}
      />

      {editTarget && (
        <OfficeFormDialog
          open={!!editTarget}
          onOpenChange={(open) => !open && setEditTarget(null)}
          title={t("admin.reference.offices.editTitle")}
          initial={editTarget}
          onSubmit={async (name, placeId) => {
            await updateOffice.mutateAsync({
              id: editTarget.id,
              data: { name, place_id: placeId || undefined },
            });
            setEditTarget(null);
          }}
          loading={updateOffice.isPending}
        />
      )}

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => !open && setDeleteTarget(null)}
        title={t("admin.common.delete")}
        description={t("admin.reference.offices.deleteConfirm")}
        onConfirm={async () => {
          if (deleteTarget) {
            await deleteOffice.mutateAsync(deleteTarget.id);
            setDeleteTarget(null);
          }
        }}
        loading={deleteOffice.isPending}
      />
    </div>
  );
}

function OfficeFormDialog({
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
  initial?: Office;
  onSubmit: (name: string, placeId: string) => Promise<void>;
  loading: boolean;
}) {
  const { t } = useTranslation();
  const [name, setName] = useState(initial?.name ?? "");
  const [placeId, setPlaceId] = useState(initial?.place_id ?? "");

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    await onSubmit(name, placeId);
    if (!initial) {
      setName("");
      setPlaceId("");
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
                {t("admin.reference.offices.name")}
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
                {t("admin.reference.offices.place")}
              </Field.Label>
              <Field.Control
                value={placeId}
                onChange={(e) => setPlaceId(e.target.value)}
                placeholder="Place ID"
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
