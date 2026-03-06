import { type SyntheticEvent, useState } from "react";

import { Field } from "@base-ui/react/field";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/button";
import { ConfirmDialog } from "@/components/confirm-dialog";
import { DataTable, type Column } from "@/components/data-table";
import { EmptyState } from "@/components/empty-state";
import { FormDialog } from "@/components/form-dialog";
import { GlobeIcon } from "@/components/icons";
import { PageHeader } from "@/components/page-header";
import { RowActions } from "@/components/row-actions";
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
      render: (c) => (
        <span className="inline-flex items-center justify-center rounded-full bg-bg-tertiary px-2.5 py-0.5 text-xs font-medium text-fg-secondary">
          {c.code}
        </span>
      ),
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
        title={t("admin.reference.countries.title")}
        action={
          <Button onClick={() => setCreateOpen(true)}>
            {t("admin.reference.countries.add")}
          </Button>
        }
      />

      <DataTable
        columns={columns}
        data={data ?? []}
        keyExtractor={(c) => c.id}
        onRowClick={(c) =>
          navigate({
            to: "/admin/reference/countries/$countryId",
            params: { countryId: c.id },
          })
        }
        isLoading={isLoading}
        emptyState={
          <EmptyState
            icon={GlobeIcon}
            title={t("admin.reference.countries.emptyTitle")}
            description={t("admin.reference.countries.emptyDescription")}
            action={
              <Button onClick={() => setCreateOpen(true)}>
                {t("admin.reference.countries.add")}
              </Button>
            }
          />
        }
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

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    await onSubmit(name, code);
    if (!initial) {
      setName("");
      setCode("");
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
          {t("admin.reference.countries.name")}
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
          {t("admin.reference.countries.code")}
        </Field.Label>
        <Field.Control
          required
          value={code}
          onChange={(e) => setCode(e.target.value)}
          className="block w-full rounded-lg border border-border-secondary bg-bg h-9 px-3 text-sm text-fg outline-none focus:border-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-bg"
        />
      </Field.Root>
    </FormDialog>
  );
}
