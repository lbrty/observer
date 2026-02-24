import { PencilSimple, Trash } from "@phosphor-icons/react";
import { Dialog } from "@base-ui/react/dialog";
import { Field } from "@base-ui/react/field";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { type FormEvent, useState } from "react";
import { useTranslation } from "react-i18next";

import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { PageHeader } from "@/components/page-header";
import {
  useCountries,
  useCreateCountry,
  useDeleteCountry,
  useUpdateCountry,
} from "@/hooks/use-countries";
import type { Country } from "@/types/reference";

export const Route = createFileRoute("/_app/admin/reference/countries/")({
  component: CountriesPage,
});

function CountriesPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { data, isLoading } = useCountries();
  const createCountry = useCreateCountry();
  const updateCountry = useUpdateCountry();
  const deleteCountry = useDeleteCountry();

  const [createOpen, setCreateOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<Country | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Country | null>(null);

  const columns: Column<Country>[] = [
    {
      key: "name",
      header: t("admin.reference.countries.name"),
      render: (c) => <span className="font-medium text-fg">{c.name}</span>,
    },
    {
      key: "code",
      header: t("admin.reference.countries.code"),
      render: (c) => <span className="text-fg-secondary">{c.code}</span>,
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
        title={t("admin.reference.countries.title")}
        action={
          <button
            type="button"
            onClick={() => setCreateOpen(true)}
            className="cursor-pointer rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-accent-fg hover:opacity-90"
          >
            {t("admin.reference.countries.add")}
          </button>
        }
      />

      <DataTable
        columns={columns}
        data={data?.countries ?? []}
        keyExtractor={(c) => c.id}
        onRowClick={(c) =>
          navigate({
            to: "/admin/reference/countries/$countryId",
            params: { countryId: c.id },
          })
        }
        isLoading={isLoading}
      />

      <CountryFormDialog
        open={createOpen}
        onOpenChange={setCreateOpen}
        title={t("admin.reference.countries.addTitle")}
        onSubmit={async (name, code) => {
          await createCountry.mutateAsync({ name, code });
          setCreateOpen(false);
        }}
        loading={createCountry.isPending}
      />

      {editTarget && (
        <CountryFormDialog
          open={!!editTarget}
          onOpenChange={(open) => !open && setEditTarget(null)}
          title={t("admin.reference.countries.editTitle")}
          initial={editTarget}
          onSubmit={async (name, code) => {
            await updateCountry.mutateAsync({
              id: editTarget.id,
              data: { name, code },
            });
            setEditTarget(null);
          }}
          loading={updateCountry.isPending}
        />
      )}

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => !open && setDeleteTarget(null)}
        title={t("admin.common.delete")}
        description={t("admin.reference.countries.deleteConfirm")}
        onConfirm={async () => {
          if (deleteTarget) {
            await deleteCountry.mutateAsync(deleteTarget.id);
            setDeleteTarget(null);
          }
        }}
        loading={deleteCountry.isPending}
      />
    </div>
  );
}

function CountryFormDialog({
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
  initial?: Country;
  onSubmit: (name: string, code: string) => Promise<void>;
  loading: boolean;
}) {
  const { t } = useTranslation();
  const [name, setName] = useState(initial?.name ?? "");
  const [code, setCode] = useState(initial?.code ?? "");

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    await onSubmit(name, code);
    if (!initial) {
      setName("");
      setCode("");
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
                {t("admin.reference.countries.name")}
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
                {t("admin.reference.countries.code")}
              </Field.Label>
              <Field.Control
                required
                value={code}
                onChange={(e) => setCode(e.target.value)}
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
