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
import { useCreateOffice, useDeleteOffice, useOffices, useUpdateOffice } from "@/hooks/use-offices";
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
        <span className="font-mono text-xs text-fg-secondary">{o.place_id ?? "—"}</span>
      ),
    },
    {
      key: "actions",
      header: "",
      render: (o) => (
        <RowActions onEdit={() => setEditTarget(o)} onDelete={() => setDeleteTarget(o)} />
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t("admin.reference.offices.title")}
        action={
          <Button onClick={() => setCreateOpen(true)}>
            {t("admin.reference.offices.add")}
          </Button>
        }
      />

      <DataTable
        columns={columns}
        data={data ?? []}
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
    <FormDialog
      open={open}
      onOpenChange={onOpenChange}
      title={title}
      loading={loading}
      onSubmit={handleSubmit}
    >
      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("admin.reference.offices.name")}
        </Field.Label>
        <Field.Control
          required
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
        />
      </Field.Root>
      <Field.Root>
        <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
          {t("admin.reference.offices.place")}
        </Field.Label>
        <Field.Control
          value={placeId}
          onChange={(e) => setPlaceId(e.target.value)}
          placeholder={t("admin.reference.offices.placeId")}
          className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
        />
      </Field.Root>
    </FormDialog>
  );
}
