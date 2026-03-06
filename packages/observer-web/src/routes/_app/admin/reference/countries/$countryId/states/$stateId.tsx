import { type FormEvent, useState } from "react";

import { ArrowLeftIcon } from "@/components/icons";
import { Field } from "@base-ui/react/field";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { FormDialog } from "@/components/form-dialog";
import { PageHeader } from "@/components/page-header";
import { RowActions } from "@/components/row-actions";
import { useCreatePlace, useDeletePlace, usePlaces, useUpdatePlace } from "@/hooks/use-places";
import type { Place } from "@/types/reference";

export const Route = createFileRoute("/_app/admin/reference/countries/$countryId/states/$stateId")({
  component: PlacesPage,
});

function PlacesPage() {
  const { t } = useTranslation();
  const { countryId, stateId } = Route.useParams();
  const { data, isLoading } = usePlaces(stateId);
  const createPlace = useCreatePlace();
  const updatePlace = useUpdatePlace();
  const deletePlace = useDeletePlace();

  const [createOpen, setCreateOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<Place | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Place | null>(null);

  const columns: Column<Place>[] = [
    {
      key: "name",
      header: t("admin.reference.places.name"),
      render: (p) => <span className="font-medium text-fg">{p.name}</span>,
    },
    {
      key: "lat",
      header: t("admin.reference.places.lat"),
      render: (p) => <span className="text-fg-secondary">{p.lat ?? "—"}</span>,
    },
    {
      key: "lon",
      header: t("admin.reference.places.lon"),
      render: (p) => <span className="text-fg-secondary">{p.lon ?? "—"}</span>,
    },
    {
      key: "actions",
      header: "",
      render: (p) => (
        <RowActions onEdit={() => setEditTarget(p)} onDelete={() => setDeleteTarget(p)} />
      ),
    },
  ];

  return (
    <div>
      <Link
        to="/admin/reference/countries/$countryId"
        params={{ countryId }}
        className="mb-3 inline-flex items-center gap-1 text-sm text-fg-tertiary hover:text-fg"
      >
        <ArrowLeftIcon size={14} />
        {t("admin.reference.states.title")}
      </Link>

      <PageHeader
        title={t("admin.reference.places.title")}
        action={
          <Button onClick={() => setCreateOpen(true)}>
            {t("admin.reference.places.add")}
          </Button>
        }
      />

      <DataTable
        columns={columns}
        data={data?.places ?? []}
        keyExtractor={(p) => p.id}
        isLoading={isLoading}
      />

      <PlaceFormDialog
        open={createOpen}
        onOpenChange={setCreateOpen}
        title={t("admin.reference.places.addTitle")}
        onSubmit={async (name, lat, lon) => {
          await createPlace.mutateAsync({
            stateId,
            data: { name, lat: lat || undefined, lon: lon || undefined },
          });
          setCreateOpen(false);
        }}
        loading={createPlace.isPending}
      />

      {editTarget && (
        <PlaceFormDialog
          open={!!editTarget}
          onOpenChange={(open) => !open && setEditTarget(null)}
          title={t("admin.reference.places.editTitle")}
          initial={editTarget}
          onSubmit={async (name, lat, lon) => {
            await updatePlace.mutateAsync({
              id: editTarget.id,
              data: { name, lat, lon },
            });
            setEditTarget(null);
          }}
          loading={updatePlace.isPending}
        />
      )}

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => !open && setDeleteTarget(null)}
        title={t("admin.common.delete")}
        description={t("admin.reference.places.deleteConfirm")}
        onConfirm={async () => {
          if (deleteTarget) {
            await deletePlace.mutateAsync(deleteTarget.id);
            setDeleteTarget(null);
          }
        }}
        loading={deletePlace.isPending}
      />
    </div>
  );
}

function PlaceFormDialog({
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
  initial?: Place;
  onSubmit: (name: string, lat?: number, lon?: number) => Promise<void>;
  loading: boolean;
}) {
  const { t } = useTranslation();
  const [name, setName] = useState(initial?.name ?? "");
  const [lat, setLat] = useState(initial?.lat?.toString() ?? "");
  const [lon, setLon] = useState(initial?.lon?.toString() ?? "");

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    await onSubmit(name, lat ? Number(lat) : undefined, lon ? Number(lon) : undefined);
    if (!initial) {
      setName("");
      setLat("");
      setLon("");
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
          {t("admin.reference.places.name")}
        </Field.Label>
        <Field.Control
          required
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent"
        />
      </Field.Root>
      <div className="grid grid-cols-2 gap-3">
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.reference.places.lat")}
          </Field.Label>
          <Field.Control
            type="number"
            step="any"
            value={lat}
            onChange={(e) => setLat(e.target.value)}
            className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent"
          />
        </Field.Root>
        <Field.Root>
          <Field.Label className="mb-1 block text-sm font-medium text-fg-secondary">
            {t("admin.reference.places.lon")}
          </Field.Label>
          <Field.Control
            type="number"
            step="any"
            value={lon}
            onChange={(e) => setLon(e.target.value)}
            className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent"
          />
        </Field.Root>
      </div>
    </FormDialog>
  );
}
