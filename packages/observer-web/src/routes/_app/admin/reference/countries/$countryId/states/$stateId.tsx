import { ArrowLeft, PencilSimple, Trash } from "@phosphor-icons/react";
import { Dialog } from "@base-ui/react/dialog";
import { Field } from "@base-ui/react/field";
import { createFileRoute, Link } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { PageHeader } from "@/components/page-header";
import {
  useCreatePlace,
  useDeletePlace,
  usePlaces,
  useUpdatePlace,
} from "@/hooks/use-places";
import type { Place } from "@/types/reference";

export const Route = createFileRoute(
  "/_app/admin/reference/countries/$countryId/states/$stateId",
)({
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
      render: (p) => (
        <span className="text-fg-secondary">{p.lat ?? "\u2014"}</span>
      ),
    },
    {
      key: "lon",
      header: t("admin.reference.places.lon"),
      render: (p) => (
        <span className="text-fg-secondary">{p.lon ?? "\u2014"}</span>
      ),
    },
    {
      key: "actions",
      header: "",
      render: (p) => (
        <div className="flex gap-2">
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              setEditTarget(p);
            }}
            className="cursor-pointer rounded p-1 text-fg-tertiary hover:bg-bg-tertiary hover:text-fg"
          >
            <PencilSimple size={16} />
          </button>
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              setDeleteTarget(p);
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
      <Link
        to="/admin/reference/countries/$countryId"
        params={{ countryId }}
        className="mb-3 inline-flex items-center gap-1 text-sm text-fg-tertiary hover:text-fg"
      >
        <ArrowLeft size={14} />
        {t("admin.reference.states.title")}
      </Link>

      <PageHeader
        title={t("admin.reference.places.title")}
        action={
          <button
            type="button"
            onClick={() => setCreateOpen(true)}
            className="cursor-pointer rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-accent-fg hover:opacity-90"
          >
            {t("admin.reference.places.add")}
          </button>
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
    await onSubmit(
      name,
      lat ? Number(lat) : undefined,
      lon ? Number(lon) : undefined,
    );
    if (!initial) {
      setName("");
      setLat("");
      setLon("");
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
                {t("admin.reference.places.name")}
              </Field.Label>
              <Field.Control
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="block w-full rounded-md border border-border-secondary bg-bg px-3 py-2 text-sm text-fg outline-none focus:border-accent"
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
                  className="block w-full rounded-md border border-border-secondary bg-bg px-3 py-2 text-sm text-fg outline-none focus:border-accent"
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
                  className="block w-full rounded-md border border-border-secondary bg-bg px-3 py-2 text-sm text-fg outline-none focus:border-accent"
                />
              </Field.Root>
            </div>
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
